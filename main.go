package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	sub            chan ResponseMsg
	currentMessage string
}

type ResponseMsg struct {
	response string
	err      error
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		listenForActivity(m.sub),
		waitForActivity(m.sub),
	)
}

func (m Model) View() string {
	return "hello the current message is " + m.currentMessage
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case ResponseMsg:
		m.currentMessage = msg.response
		return m, waitForActivity(m.sub)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, cmd
}

// realtime
func listenForActivity(sub chan ResponseMsg) tea.Cmd {
	return func() tea.Msg {
		for {
			time.Sleep(time.Second * 10)
			sub <- ResponseMsg{response: fmt.Sprintf("we are loaded at %v", time.Now()), err: nil}
		}
	}
}

func waitForActivity(sub chan ResponseMsg) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func main() {
	model := Model{sub: make(chan ResponseMsg), currentMessage: "no msg"}

	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Printf("new error %v", err)
		os.Exit(1)
	}

}
