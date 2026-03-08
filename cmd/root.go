package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Version = "1.0.0"

var (
	jsonOutput bool
	websiteID  string
)

var rootCmd = &cobra.Command{
	Use:     "datafast",
	Short:   "DataFast CLI — analytics from your terminal",
	Long:    "A command-line interface for the DataFast analytics API.\nManage your website analytics, goals, and payments.",
	Version: Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().StringVar(&websiteID, "website-id", "", "Website ID (overrides config)")
}

func initConfig() {
	if val := os.Getenv("DATAFAST_WEBSITE_ID"); val != "" && websiteID == "" {
		websiteID = val
	}
	viper.AutomaticEnv()
}
