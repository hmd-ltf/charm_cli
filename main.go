package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
)

type TaskStatus int

const (
	Todo TaskStatus = iota
	Scheduled
	InProgress
	Done
)

type AppStatus int

const (
	Loading AppStatus = iota
	Loaded
	Quit
)

const divisor = 5

var (
	columnStyle  = lipgloss.NewStyle().Padding(1, 2)
	focusedStyle = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type Task struct {
	status      TaskStatus
	title       string
	description string
}

func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

func (t *Task) Next() {
	if t.status == Done {
		t.status = Todo
	} else {
		t.status++
	}
}

type Model struct {
	lists     []list.Model
	focused   TaskStatus
	err       error
	appStatus AppStatus
}

// Call on tea.window size msg
func (m *Model) initLists(width, height int) {
	v, h := focusedStyle.GetFrameSize()
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width-h, height-v)
	defaultList.SetShowHelp(false)

	m.lists = []list.Model{defaultList, defaultList, defaultList, defaultList}
	m.lists[Todo].Title = "Todo List"
	m.lists[Todo].SetItems([]list.Item{
		&Task{status: Todo, title: "Buy Milk", description: "Buy any milk you can find for the cat"},
		&Task{status: Todo, title: "Install Linux", description: "Cause fuck windows"},
	})
	m.lists[Scheduled].Title = "Scheduled"
	m.lists[Scheduled].SetItems([]list.Item{
		&Task{status: Scheduled, title: "Read book", description: "Read the habit book"},
	})
	m.lists[InProgress].Title = "InProgress"
	m.lists[InProgress].SetItems([]list.Item{
		&Task{status: InProgress, title: "Learn Charm Cli", description: "Learn Go Lang"},
	})
	m.lists[Done].Title = "Done"
	m.lists[Done].SetItems([]list.Item{
		&Task{status: Done, title: "Don't be a piece of trash", description: "Do something with your life"},
	})
}

func (m *Model) Init() tea.Cmd {
	return nil
}
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		columnStyle.Width(msg.Width / divisor)
		focusedStyle.Width(msg.Width / divisor)
		focusedStyle.Height(msg.Height - divisor)
		columnStyle.Height(msg.Height - divisor)
		m.initLists(msg.Width, msg.Height)
		if m.appStatus != Loaded {
			m.appStatus = Loaded
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.appStatus = Quit
		case "left", "h":
			m.Previous()
		case "right", "l":
			m.Next()
		case "enter":
			return m, m.MoveToNext
		}
	}

	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}
func (m *Model) View() string {
	if m.appStatus == Loading {
		return "Loading..."
	}

	if m.appStatus == Quit {
		return ""
	}

	todoView := m.lists[Todo].View()
	scheduledView := m.lists[Scheduled].View()
	inProgressView := m.lists[InProgress].View()
	doneView := m.lists[Done].View()

	switch m.focused {
	case Todo:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			focusedStyle.Render(todoView),
			columnStyle.Render(scheduledView),
			columnStyle.Render(inProgressView),
			columnStyle.Render(doneView),
		)
	case Scheduled:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			columnStyle.Render(todoView),
			focusedStyle.Render(scheduledView),
			columnStyle.Render(inProgressView),
			columnStyle.Render(doneView),
		)
	case InProgress:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			columnStyle.Render(todoView),
			columnStyle.Render(scheduledView),
			focusedStyle.Render(inProgressView),
			columnStyle.Render(doneView),
		)
	case Done:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			columnStyle.Render(todoView),
			columnStyle.Render(scheduledView),
			columnStyle.Render(inProgressView),
			focusedStyle.Render(doneView),
		)
	}

	return ""
}

func New() *Model {
	return &Model{
		focused:   Todo,
		appStatus: Loading,
	}
}

func (m *Model) Next() {
	if m.focused == Done {
		m.focused = Todo
	} else {
		m.focused++
	}
}

func (m *Model) Previous() {
	if m.focused == Todo {
		m.focused = Done
	} else {
		m.focused--
	}
}
func (m *Model) MoveToNext() tea.Msg {
	selectedItem := m.lists[m.focused].SelectedItem()
	selectedTask, ok := selectedItem.(*Task)
	if !ok {
		//log.Printf("some")
	}
	m.lists[selectedTask.status].RemoveItem(m.lists[m.focused].Index())
	selectedTask.Next()
	m.lists[selectedTask.status].InsertItem(len(m.lists[selectedTask.status].Items())-1, list.Item(selectedTask))
	return nil
}

func main() {
	p := tea.NewProgram(New())
	_, err := p.Run()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
