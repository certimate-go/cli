package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/certimate-go/cli/internal/models"
	"github.com/certimate-go/cli/internal/output"
)

var accessCmd = &cobra.Command{
	Use:   "access",
	Short: "Manage provider access credentials",
	Long:  `List and view provider access credentials (DNS providers, CDNs, etc.).`,
}

var accessListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all access credentials",
	RunE:  runAccessList,
}

var accessGetCmd = &cobra.Command{
	Use:   "get ACCESS_ID",
	Short: "Get access details",
	Args:  cobra.ExactArgs(1),
	RunE:  runAccessGet,
}

var accessCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new access credential",
	RunE:  runAccessCreate,
}

var accessEditCmd = &cobra.Command{
	Use:   "edit ACCESS_ID",
	Short: "Edit an access credential",
	Args:  cobra.ExactArgs(1),
	RunE:  runAccessEdit,
}

var accessDeleteCmd = &cobra.Command{
	Use:   "delete ACCESS_ID",
	Short: "Delete an access credential",
	Args:  cobra.ExactArgs(1),
	RunE:  runAccessDelete,
}

var (
	accessFilter   string
	accessLimit    int
	accessReveal   bool
	accessName     string
	accessProvider string
	accessConfig   string
)

func init() {
	rootCmd.AddCommand(accessCmd)
	accessCmd.AddCommand(accessListCmd)
	accessCmd.AddCommand(accessGetCmd)
	accessCmd.AddCommand(accessCreateCmd)
	accessCmd.AddCommand(accessEditCmd)
	accessCmd.AddCommand(accessDeleteCmd)

	accessListCmd.Flags().StringVar(&accessFilter, "filter", "", "PocketBase filter expression")
	accessListCmd.Flags().IntVar(&accessLimit, "limit", 100, "Maximum number of results")
	accessListCmd.Flags().BoolVar(&accessReveal, "reveal", false, "Show sensitive configuration values")

	accessGetCmd.Flags().BoolVar(&accessReveal, "reveal", false, "Show sensitive configuration values")

	accessCreateCmd.Flags().StringVar(&accessName, "name", "", "Access name")
	accessCreateCmd.Flags().StringVar(&accessProvider, "provider", "", "Provider name (e.g., cloudflare)")
	accessCreateCmd.Flags().StringVar(&accessConfig, "config", "", "Provider configuration (JSON string or @path/to/file.json)")
	accessCreateCmd.MarkFlagRequired("name")
	accessCreateCmd.MarkFlagRequired("provider")
	accessCreateCmd.MarkFlagRequired("config")

	accessEditCmd.Flags().StringVar(&accessName, "name", "", "Access name")
	accessEditCmd.Flags().StringVar(&accessProvider, "provider", "", "Provider name")
	accessEditCmd.Flags().StringVar(&accessConfig, "config", "", "Provider configuration (JSON string or @path/to/file.json)")
}

func runAccessList(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	resp, err := client.ListAccess(context.Background(), accessFilter, "", 1, accessLimit)
	if err != nil {
		return err
	}

	// Mask sensitive fields unless --reveal is set
	if !accessReveal {
		for i := range resp.Items {
			resp.Items[i] = *resp.Items[i].MaskSensitive()
		}
	}

	return output.Print(resp, outputFmt)
}

func runAccessGet(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	access, err := client.GetAccess(context.Background(), args[0])
	if err != nil {
		return err
	}

	// Mask sensitive fields unless --reveal is set
	if !accessReveal {
		access = access.MaskSensitive()
	}

	return output.Print(access, outputFmt)
}

// parseConfigInput handles config input from string or file (@path)
func parseConfigInput(input string) (models.AccessConfig, error) {
	var data []byte
	var err error

	if filePath, ok := strings.CutPrefix(input, "@"); ok {
		data, err = os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	} else {
		data = []byte(input)
	}

	var config models.AccessConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	return config, nil
}

func runAccessCreate(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	config, err := parseConfigInput(accessConfig)
	if err != nil {
		return err
	}

	access := &models.Access{
		Name:     accessName,
		Provider: accessProvider,
		Config:   config,
	}

	createdAccess, err := client.CreateAccess(context.Background(), access)
	if err != nil {
		return err
	}

	// Mask sensitive fields by default
	createdAccess = createdAccess.MaskSensitive()

	return output.Print(createdAccess, outputFmt)
}

func runAccessEdit(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	accessID := args[0]

	// Fetch existing access to merge updates
	existingAccess, err := client.GetAccess(context.Background(), accessID)
	if err != nil {
		return err
	}

	// Apply updates
	if accessName != "" {
		existingAccess.Name = accessName
	}
	if accessProvider != "" {
		existingAccess.Provider = accessProvider
	}
	if accessConfig != "" {
		config, err := parseConfigInput(accessConfig)
		if err != nil {
			return err
		}
		existingAccess.Config = config
	}

	updatedAccess, err := client.UpdateAccess(context.Background(), accessID, existingAccess)
	if err != nil {
		return err
	}

	// Mask sensitive fields by default
	updatedAccess = updatedAccess.MaskSensitive()

	return output.Print(updatedAccess, outputFmt)
}

func runAccessDelete(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	accessID := args[0]

	if err := client.DeleteAccess(context.Background(), accessID); err != nil {
		return err
	}

	result := map[string]any{
		"status":  "deleted",
		"id":      accessID,
		"message": "Access credential deleted successfully",
	}

	return output.Print(result, outputFmt)
}
