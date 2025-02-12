/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"linear/write"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config <key> <value>",
	Short: "Sets a global configuration value",
	Long: `The only current accepted configuration value is the api key.

	This format keeps opens the possibility for settings display defaults in the
	config.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		validKeys := map[string]bool{"api_key": true, "ben": true}

		if !validKeys[key] {
			write.Std.Error("Invalid key")
			return
		}

		viper.Set(key, value)
		if err := viper.WriteConfig(); err != nil {
			write.Std.Error("Error saving config")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
