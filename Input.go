package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Input interface {
	Value() string
	Blur() tea.Msg
	Update(tea.Msg) (Input, tea.Cmd)
	View() string
}

type ShortAnswerField struct {
	textInput textinput.Model
}

func (s ShortAnswerField) Value() string {
	return s.textInput.Value()
}
func (s ShortAnswerField) Blur() tea.Msg {
	return s.textInput.Blur
}
func (s ShortAnswerField) Update(msg tea.Msg) (Input, tea.Cmd) {
	return s.textInput.Update(msg)
}
func (s ShortAnswerField) View() string {
	return s.textInput.View()
}
func NewShortAnswerField() Input {
	ti := textinput.New()
	ti.Placeholder = "Your Answer"
	ti.Focus()
	return ShortAnswerField{textInput: ti}
}

type LongAnswerField struct {
	textInput textarea.Model
}

func (l LongAnswerField) Value() string {
	return l.textInput.Value()
}
func (l LongAnswerField) Blur() tea.Msg {
	return l.textInput.Blur
}
func (l LongAnswerField) Update(msg tea.Msg) (Input, tea.Cmd) {
	return l.textInput.Update(msg)
}
func (l LongAnswerField) View() string {
	return l.textInput.View()
}
func NewLongAnswerField() Input {
	ta := textarea.New()
	ta.Placeholder = "Your Answer"
	ta.Focus()
	return LongAnswerField{textInput: ta}
}
