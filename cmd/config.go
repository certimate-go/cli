package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/certimate-go/cli/internal/api"
	"github.com/certimate-go/cli/internal/config"
	"github.com/certimate-go/cli/internal/output"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
	Long:  `Configure server connection and authentication for certimate-cli.`,
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values",
	Long: `Set the server URL and authentication token.

Example:
  certimate config set --server http://127.0.0.1:8090 --token YOUR_TOKEN`,
	RunE: runConfigSet,
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "View current configuration",
	RunE:  runConfigGet,
}

var configCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current profile name",
	RunE:  runConfigCurrent,
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured profiles",
	RunE:  runConfigList,
}

var (
	configServer   string
	configToken    string
	configValidate bool
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configCurrentCmd)
	configCmd.AddCommand(configListCmd)

	configSetCmd.Flags().StringVar(&configServer, "server", "", "Certimate server URL")
	configSetCmd.Flags().StringVar(&configToken, "token", "", "Authentication token")
	configSetCmd.Flags().BoolVar(&configValidate, "validate", true, "Validate connection after setting")

	configSetCmd.MarkFlagRequired("server")
	configSetCmd.MarkFlagRequired("token")
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	profileName := profile
	if profileName == "" {
		profileName = "default"
	}

	// Save configuration
	if err := config.Save(profileName, configServer, configToken); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	result := map[string]interface{}{
		"profile": profileName,
		"server":  configServer,
		"status":  "saved",
	}

	// Validate connection if requested
	if configValidate {
		cfg := &config.Config{
			Server: configServer,
			Token:  configToken,
		}
		client := api.NewClient(cfg)

		if err := client.ValidateConnection(context.Background()); err != nil {
			result["validation"] = "failed"
			result["validation_error"] = err.Error()
		} else {
			result["validation"] = "passed"
		}
	}

	return output.Print(result, outputFmt)
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	profiles, current, err := config.GetAllProfiles()
	if err != nil {
		return err
	}

	result := map[string]interface{}{
		"current_profile": current,
		"profiles":        profiles,
	}

	return output.Print(result, outputFmt)
}

func runConfigCurrent(cmd *cobra.Command, args []string) error {
	current := config.GetCurrentProfile()
	return output.Print(map[string]string{"current_profile": current}, outputFmt)
}

func runConfigList(cmd *cobra.Command, args []string) error {
	profiles, current, err := config.GetAllProfiles()
	if err != nil {
		return err
	}

	result := map[string]interface{}{
		"current_profile": current,
		"profiles":        profiles,
	}

	return output.Print(result, outputFmt)
}
