package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// Model: This is the main model
type Model struct {
	nestedModel    tea.Model
	currentMessage string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	return m.nestedModel.View()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case ResponseMsg:
		fmt.Printf("correct we waited")

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, cmd
}

func main() {
	nestedModel, err := initNModel()
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
	model := Model{nestedModel: nestedModel}

	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Printf("new error %v", err)
		os.Exit(1)
	}

}
