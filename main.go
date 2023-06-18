package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	questions []Question
	width     int
	height    int
	qIndex    int
	styles    *Styles
	done      bool
}

type Question struct {
	question string
	answer   string
	input    Input
}

func newQuestion(question string) Question {
	return Question{question: question}
}

func newShortQuestion(question string) Question {
	q := newQuestion(question)
	field := NewShortAnswerField()
	q.input = field
	return q
}

func newLongQuestion(question string) Question {
	q := newQuestion(question)
	field := NewLongAnswerField()
	q.input = field
	return q
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
	styles := DefaultStyles()
	return Model{
		questions: questions,
		qIndex:    0,
		styles:    styles,
		done:      false,
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
			current.answer = current.input.Value()
			if m.qIndex+1 < len(m.questions) {
				m.qIndex++
			} else {
				m.done = true
			}
			return m, current.input.Blur
		}
	}

	current.input, cmd = current.input.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	current := m.questions[m.qIndex]

	if m.done {
		var output string
		for _, q := range m.questions {
			output += fmt.Sprintf("%s: %s\n", q.question, q.answer)
		}
		return output
	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.questions[m.qIndex].question,
			m.styles.InputField.Render(current.input.View())))
}

func main() {

	m := InitialModel([]Question{newLongQuestion("how are you?"), newShortQuestion("What is your favorite movie?"), newShortQuestion("What is your favorite pokemon?")})

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
