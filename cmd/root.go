package cmd

import (
	"cmdsense/pkg/ai"
	"cmdsense/pkg/config"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cmdsense [command]",
	Short: "AI-powered shell command explainer",
	Long: `cmdsense is a CLI tool that uses AI to explain the command currently in your terminal buffer.
It helps you understand complex flags, warnings, and suggests fixes.`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		// Load Config
		if err := config.LoadConfig(); err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		// Join args in case the user didn't quote the command
		commandToExplain := strings.Join(args, " ")

		fmt.Printf("🤔 analyzing: %s...\n", commandToExplain)

		// Initialize Client
		client, err := ai.NewClient(config.AppConfig)
		if err != nil {
			fmt.Printf("Error initializing AI client: %v\n", err)
			os.Exit(1)
		}

		// Call AI
		explanation, err := client.Explain(cmd.Context(), commandToExplain)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}

		// Render Markdown
		r, _ := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(100),
		)
		out, err := r.Render(explanation)
		if err != nil {
			fmt.Printf("Error rendering output: %v\n", err)
			fmt.Println(explanation) // Fallback
			return
		}

		fmt.Print(out)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Initialize config here if needed, or in main
}
