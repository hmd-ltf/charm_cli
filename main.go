package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

type Status int

const (
	todo Status = iota
	Scheduled
	InProgress
	Done
)

type Task struct {
	status      Status
	title       string
	description string
}

func (t *Task) FilterValue() string {
	return t.title
}
func (t *Task) Title() string {
	return t.title
}
func (t *Task) Description() string {
	return t.description
}

type Model struct {
	list list.Model
	err  error
}

// Call on tea.window size msg
func (m *Model) initList(width, height int) {
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	m.list.Title = "Todo List"
	m.list.SetItems([]list.Item{
		&Task{status: todo, title: "Buy Milk", description: "Buy any milk you can find for the cat"},
		&Task{status: todo, title: "Install Linux", description: "Cause fuck windows"},
		&Task{status: todo, title: "Learn Charm Cli", description: "Learn Go Lang"},
	})
}

func (m *Model) Init() tea.Cmd {
	return nil
}
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.initList(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
func (m *Model) View() string {
	return m.list.View()
}

func New() *Model {
	return &Model{}
}

func main() {
	p := tea.NewProgram(New())
	_, err := p.Run()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
