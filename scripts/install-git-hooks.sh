#!/bin/sh

# Git hooks installation script

# Set the source and destination directories
HOOK_DIR="scripts/git-hooks"
GIT_HOOK_DIR=".git/hooks"

# Ensure we're in the project root
if [ ! -d "$HOOK_DIR" ]; then
    echo "Error: Could not find hooks directory ($HOOK_DIR)"
    exit 1
fi

# Create git hooks directory if it doesn't exist
mkdir -p "$GIT_HOOK_DIR"

# Install all hooks
for hook in "$HOOK_DIR"/*; do
    if [ -f "$hook" ]; then
        # Get the hook name
        hook_name=$(basename "$hook")
        # Create symlink
        ln -sf "../../$HOOK_DIR/$hook_name" "$GIT_HOOK_DIR/$hook_name"
        # Make it executable
        chmod +x "$GIT_HOOK_DIR/$hook_name"
        echo "Installed: $hook_name"
    fi
done

echo "Git hooks installation completed! 🎉"
echo "The following hooks are now active:"
ls -l "$GIT_HOOK_DIR"
