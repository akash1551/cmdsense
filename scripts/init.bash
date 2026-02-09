# explain-cmd Bash Integration

_explain_cmd_bash() {
    local cmd="$READLINE_LINE"
    
    if [[ -z "$cmd" ]]; then
        return
    fi

    echo ""
    explain-cmd "$cmd"
    
    # Reprint the prompt and the current line
    printf "\n"
    # Detect prompt if possible, or just force a redraw (Bash is trickier than Zsh)
    READLINE_POINT=$READLINE_POINT
}

# Bind to Ctrl+E
bind -x '"\C-e": _explain_cmd_bash'
