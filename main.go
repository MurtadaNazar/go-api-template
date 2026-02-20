package main

import (
	"fmt"
	"os"

	"go_platform_template/internal/scaffold"

	tea "github.com/charmbracelet/bubbletea"
)

func init() {
	scaffold.SetScaffoldFS(ScaffoldFS)
}

func main() {
	p := tea.NewProgram(scaffold.NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
