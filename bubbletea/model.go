package bubble

import (
	"fmt"
	"math/rand"
	"pitchmind/constants"
	"pitchmind/llm"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Page string

const (
	PageIntro     Page = "intro"
	PageHowToPlay Page = "howtoplay"
	PageLicense   Page = "license"
	PagePlayers   Page = "players"
	PageGame      Page = "game"
)

type Model struct {
	input         textinput.Model
	currentPage   Page
	showMainMenu  bool
	menuCursor    int
	menuChoices   []string
	playersList   []string
	Clues         []string
	attempts      int
	feedback      string
	gameOver      bool
	headerMessage string
	width         int
	height        int
}

func InitialModel() Model {

	ti := textinput.New()
	ti.Placeholder = "Enter your guess..."
	ti.Focus()

	currentPlayer = constants.GeneratePlayerList()[rand.Intn(len(constants.GeneratePlayerList()))]

	m := Model{
		input:         ti,
		Clues:         []string{},
		attempts:      1,
		headerMessage: headerStyle.Render("⚽ Guess the Football Player! (5 attempts) ⚽"),
		currentPage:   PageIntro,
		showMainMenu:  false,
		menuCursor:    0,
		menuChoices:   []string{"Start Game", "How to Play", "Player List", "License", "Quit"},
		playersList:   constants.GeneratePlayerList(),
	}

	m.getFirstClue()
	return m
}

func (m Model) RenderIntro() string {
	introText := titleStyle.Render("⚽ Guess the Football Legend! ⚽") + "\n\n" +
		"  Crafted by " + playerStyle.Render("Dev Baheti") +
		" - the Messi of code\n" +
		"  (but with better hair and worse footwork)\n\n" +
		"  \"Because guessing Cristiano's age should be a sport too!\"\n\n" +
		menuStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				"Press any key to continue...",
				creditStyle.Render("Ctrl+C to quit"),
			),
		)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		introText,
	)
}

func (m Model) RenderHowToPlay() string {
	content := []string{
		titleStyle.Render("📖 How to Play"),
		"",
		"1. You get 5 attempts to guess the player",
		"2. Each wrong guess reveals a new clue",
		"3. Clues get more obvious (and more ridiculous)",
		"4. Partial names work! (e.g. 'Messi' for Lionel Messi)",
		"5. If you lose... well, better luck next time!",
		"",
		menuStyle.Render("Press ESC to return"),
	}
	return lipgloss.JoinVertical(lipgloss.Center, content...)
}

func (m Model) RenderLicense() string {
	return lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Render("📜 MIT License"),
		"",
		"Copyright (c) 2025 Dev Baheti",
		"",
		"Permission is hereby granted, free of charge, to any person obtaining a copy",
		"of this software and associated documentation files (the \"Software\"), to deal",
		"in the Software without restriction, including without limitation the rights",
		"to use, copy, modify, merge, publish, distribute, sublicense, and/or sell",
		"copies of the Software, and to permit persons to whom the Software is",
		"furnished to do so, subject to the following conditions:",
		"",
		"The above copyright notice and this permission notice shall be included in all",
		"copies or substantial portions of the Software.",
		"",
		"THE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR",
		"IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,",
		"FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE",
		"AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER",
		"LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,",
		"OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE",
		"SOFTWARE.",
		"",
		menuStyle.Render("Press ESC to return"),
	)
}

func (m Model) RenderPlayersList() string {
	columns := make([]string, 0, 5)
	colSize := len(m.playersList)/5 + 1

	playerItemStyle := lipgloss.NewStyle().
		MarginBottom(1).
		MarginRight(2).
		PaddingLeft(1).
		Foreground(lipgloss.Color("255"))

	for i := 0; i < 5; i++ {
		start := i * colSize
		end := start + colSize
		if end > len(m.playersList) {
			end = len(m.playersList)
		}

		columnItems := make([]string, 0)
		for _, player := range m.playersList[start:end] {
			columnItems = append(columnItems, playerItemStyle.Render(player))
		}

		columns = append(columns,
			lipgloss.JoinVertical(lipgloss.Left, columnItems...))
	}

	return lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Render("🌟 Star Players List"),
		"\n",
		lipgloss.JoinHorizontal(lipgloss.Top, columns...),
		"\n\n"+menuStyle.Render("Press ESC to return"),
	)
}

func (m Model) renderTopBar() string {
	var tabs []string
	pages := []struct {
		key  string
		name string
	}{
		{"SPACE", "Intro"},
		{"1", "How To Play"},
		{"2", "License"},
		{"3", "Players"},
		{"ESC", "Back"},
	}

	for _, p := range pages {
		tab := fmt.Sprintf("%s (%s)", p.name, p.key)
		if m.currentPage == PageIntro && p.name == "Intro" ||
			m.currentPage == PageHowToPlay && p.name == "How To Play" ||
			m.currentPage == PageLicense && p.name == "License" ||
			m.currentPage == PagePlayers && p.name == "Players" {
			tabs = append(tabs, activeTabStyle.Render(tab))
		} else {
			tabs = append(tabs, inactiveTabStyle.Render(tab))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func (m *Model) processGuess() (tea.Model, tea.Cmd) {
	guess := strings.TrimSpace(m.input.Value())
	m.input.Reset()

	if strings.Contains(strings.ToLower(currentPlayer), strings.ToLower(guess)) {
		m.gameOver = true
		m.headerMessage = successStyle.Render(fmt.Sprintf("🎉 Correct! You guessed it in %d attempts!", m.attempts))
		return m, nil
	}

	m.feedback = errorStyle.Render("❌ Incorrect. Next clue...")
	m.attempts++
	if m.attempts > 5 {
		m.gameOver = true
		return m, nil
	}

	m.getNextClue()
	return m, nil
}

func (m *Model) getFirstClue() {
	clue := m.getClueFromAI("Give the first funny and witty clue")
	m.Clues = append(m.Clues, clue)
}

func (m *Model) getNextClue() {
	clue := m.getClueFromAI("Give the next clue, make it more specific but still funny")
	m.Clues = append(m.Clues, clue)
}

func (m *Model) getClueFromAI(prompt string) string {
	return llm.New().GenerateClueFromLLM(m.Clues, currentPlayer, prompt)
}

type clueMsg struct{ clue string }
