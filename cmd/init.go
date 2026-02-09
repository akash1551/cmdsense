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
		key, _ := cmd.Flags().GetString("key")

		switch shell {
		case "zsh":
			printZshScript(key)
		case "bash":
			printBashScript(key)
		default:
			fmt.Printf("Unsupported shell: %s. Supported shells: zsh, bash.\n", shell)
			os.Exit(1)
		}
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

func printBashScript(key string) {
	exePath, err := os.Executable()
	if err != nil {
		exePath = "cmdsense" // Fallback
	}

	// Map common keys to bash bind format if needed, but for now assumption is user provides correct bind key or we use default
	// Bash bind -x uses different syntax for keys than zsh, but for simplicity we'll assume standard sequences or default to Ctrl+g
	if key == "^g" {
		key = "\\C-g"
	}

	script := fmt.Sprintf(`
# cmdsense bash integration
_cmdsense_bash() {
    local cmd="$READLINE_LINE"

    if [[ -z "$cmd" ]]; then
        return
    fi

    echo ""
    # Call cmdsense
    "%s" "$cmd"

    # Reprint the prompt and the current line
    printf "\n"
    
    # Refresh the prompt
    if [ -n "$READLINE_LINE" ]; then
        READLINE_POINT="$READLINE_POINT"
    fi
}

# Bind to key
bind -x '"%s": _cmdsense_bash'
`, exePath, key)

	fmt.Println(script)
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().String("key", "^g", "Key binding for the widget (default: Ctrl+g)")
}
