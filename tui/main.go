package main

import (
	"fmt"
	"os"

	"github.com/awwwkshay/alpha-aide/agent/config"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg := config.Load()

	m, err := initialModel(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
