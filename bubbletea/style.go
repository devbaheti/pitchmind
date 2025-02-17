package bubble

import "github.com/charmbracelet/lipgloss"

var (
	playerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("36")).Bold(true)
	clueStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Margin(1, 2)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true)
	promptStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Margin(1, 2)
	headerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("57")).Bold(true)
	borderStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Foreground(lipgloss.Color("99"))
	menuStyle    = lipgloss.NewStyle().Padding(1, 2).Margin(1, 2).Border(lipgloss.RoundedBorder())
	titleStyle   = lipgloss.NewStyle().
			Foreground(lipgloss.Color("228")).
			Background(lipgloss.Color("57")).
			Bold(true).
			Padding(0, 1).
			Margin(1, 0)
	creditStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
	currentPlayer  string
	activeTabStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("57")).
			Padding(0, 1).
			MarginRight(1)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Padding(0, 1).
				MarginRight(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			PaddingTop(1)
)
