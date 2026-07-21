package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/certimate-go/cli/internal/output"
)

var certificateCmd = &cobra.Command{
	Use:   "certificate",
	Short: "Manage SSL certificates",
	Long:  `List, view, and download SSL certificates.`,
}

var certificateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all certificates",
	RunE:  runCertificateList,
}

var certificateGetCmd = &cobra.Command{
	Use:   "get CERTIFICATE_ID",
	Short: "Get certificate details",
	Args:  cobra.ExactArgs(1),
	RunE:  runCertificateGet,
}

var certificateDownloadCmd = &cobra.Command{
	Use:   "download CERTIFICATE_ID",
	Short: "Download certificate",
	Args:  cobra.ExactArgs(1),
	RunE:  runCertificateDownload,
}

var (
	certFilter string
	certLimit  int
	certFormat string
	certDest   string
)

func init() {
	rootCmd.AddCommand(certificateCmd)
	certificateCmd.AddCommand(certificateListCmd)
	certificateCmd.AddCommand(certificateGetCmd)
	certificateCmd.AddCommand(certificateDownloadCmd)

	certificateListCmd.Flags().StringVar(&certFilter, "filter", "", "PocketBase filter expression")
	certificateListCmd.Flags().IntVar(&certLimit, "limit", 100, "Maximum number of results")

	certificateDownloadCmd.Flags().StringVar(&certFormat, "format", "PEM", "Certificate format (PEM, PFX, JKS)")
	// Intentionally not "--output": that is the global output-format flag
	// (json|table). --dest is the download destination instead.
	certificateDownloadCmd.Flags().StringVar(&certDest, "dest", "", "Destination: a directory (writes <domain>-<id>.zip inside), a .zip filename, or '-' for stdout (default: current directory)")
}

func runCertificateList(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	resp, err := client.ListCertificates(context.Background(), certFilter, "", 1, certLimit)
	if err != nil {
		return err
	}

	return output.Print(resp, outputFmt)
}

func runCertificateGet(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	cert, err := client.GetCertificate(context.Background(), args[0])
	if err != nil {
		return err
	}

	return output.Print(cert, outputFmt)
}

func runCertificateDownload(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	data, err := client.DownloadCertificate(ctx, args[0], certFormat)
	if err != nil {
		return err
	}

	// "-" streams the raw archive bytes to stdout.
	if certDest == "-" {
		_, err = os.Stdout.Write(data)
		return err
	}

	// An empty --dest means the current directory.
	dest := certDest
	if dest == "" {
		dest = "."
	}

	// A directory (existing dir, or a trailing path separator) means "place a
	// <domain>-<certificateID>.zip inside it"; any other value is the file
	// path itself, gaining a ".zip" suffix when missing.
	destIsDir := strings.HasSuffix(dest, "/") || strings.HasSuffix(dest, string(os.PathSeparator))
	if !destIsDir {
		if info, statErr := os.Stat(dest); statErr == nil && info.IsDir() {
			destIsDir = true
		}
	}

	// Only the directory case needs the domain for the generated filename.
	var domain string
	if destIsDir {
		if cert, getErr := client.GetCertificate(ctx, args[0]); getErr == nil {
			domain = canonicalDomain(cert.SubjectAltNames)
		}
	}

	outPath := resolveDownloadDestPath(dest, destIsDir, domain, args[0])

	if destIsDir {
		if err := os.MkdirAll(dest, 0o755); err != nil {
			return fmt.Errorf("create output directory: %w", err)
		}
	}

	if err := os.WriteFile(outPath, data, 0600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	result := map[string]interface{}{
		"status":     "downloaded",
		"format":     certFormat,
		"dest":       outPath,
		"size_bytes": len(data),
	}
	return output.Print(result, outputFmt)
}

// canonicalDomain returns the certificate's primary SAN with wildcard
// asterisks replaced by underscores, matching the server's archive naming.
func canonicalDomain(sans []string) string {
	if len(sans) == 0 {
		return ""
	}
	return strings.ReplaceAll(sans[0], "*", "_")
}

// resolveDownloadDestPath decides where to write the downloaded ZIP archive.
//
// When dest is a directory, the file is written inside it as
// "<domain>-<certificateID>.zip" (or "<certificateID>.zip" when the domain is
// unknown). Otherwise dest is treated as a file path: names ending in ".zip"
// (case-insensitive) are used verbatim, and any other name gets ".zip" appended.
func resolveDownloadDestPath(dest string, destIsDir bool, domain, certificateID string) string {
	if destIsDir {
		name := certificateID + ".zip"
		if domain != "" {
			name = fmt.Sprintf("%s-%s.zip", domain, certificateID)
		}
		return filepath.Join(dest, name)
	}
	if strings.HasSuffix(strings.ToLower(dest), ".zip") {
		return dest
	}
	return dest + ".zip"
}
