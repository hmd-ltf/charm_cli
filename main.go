package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
	"os"
)

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}
type Question struct {
	question  string
	answer    string
	inputType Input
}

type Model struct {
	questions []Question
	index     int
	width     int
	height    int
	styles    *Styles
	done      bool
}

func defaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = "56"
	s.InputField = lipgloss.
		NewStyle().
		BorderForeground(s.BorderColor).
		BorderStyle(lipgloss.NormalBorder()).
		Padding(1).
		Width(80)

	return s
}

func New(questions []Question) *Model {
	styles := defaultStyles()
	answerField := textinput.New()
	answerField.Placeholder = "Answer"
	answerField.Focus()
	return &Model{questions: questions, styles: styles, done: false}
}

func NewQuestion(question string) *Question {
	return &Question{question: question}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	current := &m.questions[m.index]
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			current.answer = current.inputType.Value()
			log.Printf("Question: %s, Answer: %s", current.question, current.answer)
			if m.index < len(m.questions) {
				m.index++
			} else {
				m.done = true
			}

			return m, current.inputType.Blur
		}
	}
	current.inputType, cmd = current.inputType.Update(msg)
	return m, cmd
}

func newShortQuestion(question string) *Question {
	shortQuestion := NewQuestion(question)
	shortQuestion.inputType = NewShortAnswerField()
	return shortQuestion
}

func newLongQuestion(question string) *Question {
	longQuestion := NewQuestion(question)
	longQuestion.inputType = NewLongAnswerField()
	return longQuestion
}

func (m *Model) View() string {
	current := m.questions[m.index]
	if m.width <= 0 || m.height <= 0 {
		return "Loading..."
	}

	if m.done {
		var output string
		for _, question := range m.questions {
			output += fmt.Sprintf("%s\n\n%s", question.question, question.answer)
		}
		return output
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			current.question,
			m.styles.InputField.Render(current.inputType.View()),
		),
	)
}

func (m *Model) next() {
	if m.index < len(m.questions) {
		m.index++
	} else {
		m.index = 0
	}

}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("err: %v", err)
		}
	}(f)

	model := New([]Question{
		*newShortQuestion("What is your name?"),
		*newShortQuestion("What is your age?"),
		*newLongQuestion("Where do you live?"),
	})

	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err = p.Run()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

}
