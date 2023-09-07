package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Exploring nested model and realtime

type NModel struct {
	sub   chan ResponseMsg
	title string
}

type ResponseMsg struct {
	response string
	err      error
}

func initNModel() (NModel, error) {
	var model NModel
	model.sub = make(chan ResponseMsg)
	model.title = ""

	return model, nil
}

func (m NModel) Init() tea.Cmd {
	return tea.Batch(
		listenForActivity(m.sub),
		waitForActivity(m.sub),
	)
}

func (m NModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m NModel) View() string {
	return fmt.Sprintf("time: %s", m.title)
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
