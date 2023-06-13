package main

import (
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	questions   []Question
	width       int
	height      int
	answerField textinput.Model
	qIndex      int
	styles      *Styles
}

type Question struct {
	question string
	answer   string
}

func newQuestion(question string) Question {
	return Question{question: question}
}

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

func DefaultStyles() *Styles {
	return &Styles{
		BorderColor: lipgloss.Color("#1ADFCD"),
		InputField:  lipgloss.NewStyle().BorderForeground(lipgloss.Color("#1ADFCD")).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80),
	}
}

func InitialModel(questions []Question) Model {
	answerField := textinput.New()
	answerField.Placeholder = "Enter your answer"
	answerField.Focus()
	styles := DefaultStyles()
	return Model{
		questions:   questions,
		answerField: answerField,
		qIndex:      0,
		styles:      styles,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	current := &m.questions[m.qIndex]
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			current.answer = m.answerField.Value()
			m.answerField.SetValue("")
			log.Printf("question: %s, answer: %s", current.question, current.answer)
			if m.qIndex+1 < len(m.questions) {
				m.qIndex++
			}
			return m, nil
		}
	}

	m.answerField, cmd = m.answerField.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.questions[m.qIndex].question,
			m.styles.InputField.Render(m.answerField.View())))
}

func main() {

	m := InitialModel([]Question{newQuestion("how are you?"), newQuestion("What is your favorite movie?"), newQuestion("What is your favorite pokemon?")})

	f, err := tea.LogToFile("debug.log", "prefix")
	if err != nil {
		log.Fatalf("err : %v", err)
	}

	defer f.Close()

	prg := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := prg.Run(); err != nil {
		log.Fatal(err)
	}
}
