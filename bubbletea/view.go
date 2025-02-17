package bubble

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {

	topBar := m.renderTopBar()
	help := helpStyle.Render("Ctrl+C to quit • Enter to submit")

	var mainContent string

	switch m.currentPage {
	case PageIntro:
		mainContent = m.RenderIntro()
	case PageHowToPlay:
		mainContent = m.RenderHowToPlay()
	case PageLicense:
		mainContent = m.RenderLicense()
	case PagePlayers:
		mainContent = m.RenderPlayersList()
	case PageGame:
		if m.gameOver {
			return borderStyle.Render(
				lipgloss.JoinVertical(
					lipgloss.Center,
					headerStyle.Render("Game Over!"),
					clueStyle.Render(fmt.Sprintf("The player was: %s", playerStyle.Render(currentPlayer))),
					promptStyle.Render("Press Enter to exit"),
				),
			)
		}

		clues := strings.Join(m.Clues, "\n\n")
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			m.headerMessage,
			clueStyle.Render(clues),
			borderStyle.Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					promptStyle.Render(fmt.Sprintf("Attempt %d/5", m.attempts)),
					m.input.View(),
					errorStyle.Render(m.feedback),
				),
			),
		)

		return lipgloss.NewStyle().Padding(1, 2).Render(content)
	default:
		return ""
	}

	fullUI := lipgloss.JoinVertical(lipgloss.Left,
		topBar,
		lipgloss.NewStyle().Padding(1).Render(mainContent),
		help,
	)

	if m.currentPage == PageIntro {
		return lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			fullUI,
		)
	}

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Render(fullUI)

}
