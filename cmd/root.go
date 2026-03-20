package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	outputFmt string
	profile   string

	// Version info set by main
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "certimate",
	Short: "CLI for Certimate SSL certificate management",
	Long: `Certimate CLI manages SSL certificate workflows, certificates,
and access credentials through the command line.

Configure with: certimate config set --server URL --token TOKEN
Then use: certimate workflow list, certimate certificate list, etc.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.config/certimate-cli/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&outputFmt, "output", "o", "json", "output format (json|table)")
	rootCmd.PersistentFlags().StringVar(&profile, "profile", "default", "configuration profile")
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug output")

	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		viper.AddConfigPath(home + "/.config/certimate-cli")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("CERTIMATE_CLI")
	viper.AutomaticEnv()

	// Read config file (ignore error if not exists)
	viper.ReadInConfig()
}
