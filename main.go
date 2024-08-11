package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"net/http"
	"os"
	"time"
)

const defaultUrl = "https://charm.sh"

type errMsg struct{ err error }
type statusMsg int

type model struct {
	status int
	err    error
}

func (e errMsg) Error() string {
	return e.err.Error()
}

func checkServer(url string) tea.Cmd {
	return func() tea.Msg {
		c := &http.Client{Timeout: 10 * time.Second}
		resp, err := c.Get(url)

		if err != nil {
			return errMsg{err: err}
		}

		defer resp.Body.Close()

		return statusMsg(resp.StatusCode)
	}
}

func (m model) Init() tea.Cmd {
	return checkServer(defaultUrl)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusMsg:
		m.status = int(msg)
		return m, tea.Quit
	case errMsg:
		m.err = msg.err
		return m, tea.Quit
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nError: %v\n\n", m.err)
	}

	s := fmt.Sprintf("Checking status of %s... \n\n", defaultUrl)

	if m.status > 0 {
		s += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}

	return s
}

func main() {
	_, err := tea.NewProgram(model{}).Run()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

}
