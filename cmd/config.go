package cmd

import (
	"cmdsense/pkg/config"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.LoadConfig(); err != nil {
			fmt.Printf("Error loading config: %v\n", err)
		}

		provider, _ := cmd.Flags().GetString("provider")
		openaiKey, _ := cmd.Flags().GetString("openai-key")
		ollamaURL, _ := cmd.Flags().GetString("ollama-url")
		model, _ := cmd.Flags().GetString("model")

		changed := false

		if provider != "" {
			viper.Set("provider", provider)
			fmt.Printf("Provider set to: %s\n", provider)
			changed = true
		}
		if openaiKey != "" {
			viper.Set("openai_key", openaiKey)
			fmt.Println("OpenAI API Key updated.")
			changed = true
		}
		if ollamaURL != "" {
			viper.Set("ollama_base_url", ollamaURL)
			fmt.Printf("Ollama URL set to: %s\n", ollamaURL)
			changed = true
		}
		if model != "" {
			viper.Set("model", model)
			fmt.Printf("Model set to: %s\n", model)
			changed = true
		}

		if changed {
			if err := config.SaveConfig(); err != nil {
				fmt.Printf("Error saving config: %v\n", err)
			} else {
				fmt.Println("Configuration saved successfully.")
			}
		} else {
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().String("provider", "", "Set AI provider (openai/ollama)")
	configCmd.Flags().String("openai-key", "", "Set OpenAI API Key")
	configCmd.Flags().String("ollama-url", "", "Set Ollama Base URL")
	configCmd.Flags().String("model", "", "Set Model Name (e.g. gpt-4o, llama3)")
}
