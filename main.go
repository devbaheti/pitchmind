package main

import (
	"fmt"
	bubble "pitchmind/bubbletea"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(bubble.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
