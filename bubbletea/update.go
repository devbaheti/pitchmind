package bubble

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyRunes:
			switch strings.ToLower(string(msg.Runes)) {
			case "1":
				m.currentPage = PageHowToPlay
				return m, nil
			case "2":
				m.currentPage = PageLicense
				return m, nil
			case "3":
				m.currentPage = PagePlayers
				return m, nil
			}

		case tea.KeySpace:
			if m.currentPage == PageIntro {
				m.currentPage = PageGame
				return m, nil
			}

		case tea.KeyEsc:
			if m.currentPage != PageGame {
				m.currentPage = PageIntro
			}
			return m, nil

		}
		switch m.currentPage {
		case PageIntro:
			if msg.Type == tea.KeyEnter || msg.Type == tea.KeySpace {
				m.currentPage = PageGame
			}
		case PageHowToPlay, PageLicense, PagePlayers:
			if msg.Type == tea.KeyEsc {
				m.currentPage = PageIntro
			}
		case PageGame:
			switch msg.Type {
			case tea.KeyEnter:
				if m.gameOver {
					return m, tea.Quit
				}
				return m.processGuess()
			case tea.KeyCtrlC:
				return m, tea.Quit
			}
		}
	case clueMsg:
		m.Clues = append(m.Clues, msg.clue)
	}
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}
