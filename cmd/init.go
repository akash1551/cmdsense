package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [shell]",
	Short: "Generate shell integration script",
	Long:  `Generate shell integration script for your shell. Currently supports zsh.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shell := args[0]
		if shell != "zsh" {
			fmt.Printf("Unsupported shell: %s. Only 'zsh' is currently supported.\n", shell)
			os.Exit(1)
		}

		key, _ := cmd.Flags().GetString("key")
		printZshScript(key)
	},
}

func printZshScript(key string) {
	exePath, err := os.Executable()
	if err != nil {
		exePath = "cmdsense" // Fallback
	}

	script := fmt.Sprintf(`
# cmdsense shell integration
_cmdsense_widget() {
    # Get the current buffer
    local current_cmd="$BUFFER"

    if [[ -z "$current_cmd" ]]; then
        echo "No command to explain."
        return
    fi

    echo "" # New line
    # Call cmdsense with the buffer content
    # We use the absolute path to ensure it works even if not in PATH
    "%s" "$current_cmd"

    # Redraw prompt
    zle reset-prompt
}

# Register the widget
zle -N _cmdsense_widget

# Bind the key
bindkey '%s' _cmdsense_widget
`, exePath, key)

	fmt.Println(script)
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().String("key", "^g", "Key binding for the widget (default: Ctrl+g)")
}
