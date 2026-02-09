# explain-cmd Zsh Integration

explain-command-widget() {
    # Get the current buffer
    local cmd="$BUFFER"
    
    # If buffer is empty, do nothing
    if [[ -z "$cmd" ]]; then
        return
    fi

    # Print a newline to avoid overwriting the prompt
    echo ""
    
    # Call the explain-cmd binary
    # We use -- to ensure flags in the command aren't parsed by explain-cmd
    explain-cmd "$cmd"
    
    # Redraw the prompt
    zle reset-prompt
}

# Register the widget
zle -N explain-command-widget

# Bind to Ctrl+E
bindkey '^e' explain-command-widget
