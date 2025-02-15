package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	playerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("36")).Bold(true)
	clueStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Margin(1, 2)
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	successStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true)
	promptStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Margin(1, 2)
	headerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("57")).Bold(true)
	borderStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Foreground(lipgloss.Color("99"))
	players       = generatePlayerList()
	client        = &http.Client{Timeout: 30 * time.Second}
	currentPlayer string
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ChatResponse struct {
	Message Message `json:"message"`
}

type model struct {
	input         textinput.Model
	clues         []string
	attempts      int
	feedback      string
	gameOver      bool
	headerMessage string
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter your guess..."
	ti.Focus()

	currentPlayer = players[rand.Intn(len(players))]

	m := model{
		input:         ti,
		clues:         []string{},
		attempts:      1,
		headerMessage: headerStyle.Render("⚽ Guess the Football Player! (5 attempts) ⚽"),
	}

	m.getFirstClue()
	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.gameOver {
				return m, tea.Quit
			}
			return m.processGuess()
		case tea.KeyCtrlC:
			return m, tea.Quit
		}

	case clueMsg:
		m.clues = append(m.clues, msg.clue)
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
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

	clues := strings.Join(m.clues, "\n\n")
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
}

func (m *model) processGuess() (tea.Model, tea.Cmd) {
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

func (m *model) getFirstClue() {
	clue := m.getClueFromAI("Give the first funny and witty clue")
	m.clues = append(m.clues, clue)
}

func (m *model) getNextClue() {
	clue := m.getClueFromAI("Give the next clue, make it more specific but still funny")
	m.clues = append(m.clues, clue)
}

func (m *model) getClueFromAI(prompt string) string {
	messages := []Message{
		{
			Role: "system",
			Content: fmt.Sprintf(
				"You're a football expert comedian. The player is %s. "+
					"Give %d funny, witty clues using wordplay, pop culture refs, and funny analogies. "+
					"Never mention the name, team, or nationality directly. "+
					"Make it humorous and engaging under 40 characters! Format: emoji + clue",
				currentPlayer, len(m.clues)+1,
			),
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	requestBody := ChatRequest{
		Model:    "llama3.2",
		Messages: messages,
		Stream:   false,
	}

	jsonBody, _ := json.Marshal(requestBody)
	resp, err := client.Post("http://localhost:11434/api/chat", "application/json", strings.NewReader(string(jsonBody)))
	if err != nil {
		return "🚧 Oops, the clue machine broke! Try another guess..."
	}
	defer resp.Body.Close()

	var chatResp ChatResponse
	json.NewDecoder(resp.Body).Decode(&chatResp)
	return chatResp.Message.Content
}

type clueMsg struct{ clue string }

func generatePlayerList() []string {
	return []string{
		"Lionel Messi", "Cristiano Ronaldo", "Neymar Jr.", "Kylian Mbappé", "Robert Lewandowski",
		"Mohamed Salah", "Kevin De Bruyne", "Virgil van Dijk", "Erling Haaland", "Karim Benzema",
		"Zlatan Ibrahimović", "Gareth Bale", "Luka Modrić", "Sergio Ramos", "Manuel Neuer",
		"Harry Kane", "Toni Kroos", "Sadio Mané", "Paul Pogba", "Antoine Griezmann",
		"Eden Hazard", "Raheem Sterling", "Joshua Kimmich", "Jan Oblak", "Alisson Becker",
		"Ederson", "Riyad Mahrez", "Son Heung-min", "Romelu Lukaku", "Bruno Fernandes",
		"Frenkie de Jong", "Marquinhos", "Kalvin Phillips", "Phil Foden", "Jadon Sancho",
		"Marcus Rashford", "Trent Alexander-Arnold", "Andrew Robertson", "Thiago Alcântara",
		"Gerard Piqué", "Sergio Agüero", "Luis Suárez", "Ángel Di María", "N'Golo Kanté",
		"Gianluigi Buffon", "Andrea Pirlo", "Xabi Alonso", "Andrés Iniesta", "Xavi Hernández",
		"David Beckham", "Ronaldinho", "Thierry Henry", "Didier Drogba", "Steven Gerrard",
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
