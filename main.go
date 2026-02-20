package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"go_platform_template/internal/scaffold"
)

func init() {
	// Inject embedded scaffold FS into scaffold package
	scaffold.SetScaffoldFS(ScaffoldFS)
}

func main() {
	p := tea.NewProgram(scaffold.NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
