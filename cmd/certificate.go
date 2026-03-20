package cmd

import (
	"context"
	"fmt"
	"os"

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
	certOutput string
)

func init() {
	rootCmd.AddCommand(certificateCmd)
	certificateCmd.AddCommand(certificateListCmd)
	certificateCmd.AddCommand(certificateGetCmd)
	certificateCmd.AddCommand(certificateDownloadCmd)

	certificateListCmd.Flags().StringVar(&certFilter, "filter", "", "PocketBase filter expression")
	certificateListCmd.Flags().IntVar(&certLimit, "limit", 100, "Maximum number of results")

	certificateDownloadCmd.Flags().StringVar(&certFormat, "format", "PEM", "Certificate format (PEM, PFX, JKS)")
	certificateDownloadCmd.Flags().StringVar(&certOutput, "output", "", "Output file path (default: stdout)")
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
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	data, err := client.DownloadCertificate(context.Background(), args[0], certFormat)
	if err != nil {
		return err
	}

	if certOutput != "" {
		if err := os.WriteFile(certOutput, data, 0600); err != nil {
			return fmt.Errorf("write file: %w", err)
		}
		result := map[string]interface{}{
			"status":     "downloaded",
			"format":     certFormat,
			"output":     certOutput,
			"size_bytes": len(data),
		}
		return output.Print(result, outputFmt)
	}

	// Output to stdout
	fmt.Print(string(data))
	return nil
}
