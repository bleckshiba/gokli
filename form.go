package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

type Form struct {
	help   help.Model
	title  textinput.Model
	desc   textarea.Model
	column Column
	index  int
}

func newDefaultForm() *Form {
	return NewForm("Task name", "")
}

func NewForm(title, desc string) *Form {
	form := Form{
		help:  help.New(),
		title: textinput.New(),
		desc:  textarea.New(),
	}
	form.title.Placeholder = title
	form.desc.Placeholder = desc
	form.title.Focus()
	return &form
}

func (f Form) AddTask() Task {
	return Task{
		ID:     uuid.New(),
		title:  f.title.Value(),
		Desc:   f.desc.Value(),
		Status: f.column.status,
	}
}

func (f Form) Init() tea.Cmd {
	return nil
}

func (f Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case Column:
		f.column = msg
		f.column.list.Index()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Exit):
			return f, tea.Quit

		case key.Matches(msg, keys.Back):
			return board.Update(nil)
		case key.Matches(msg, keys.Select):
			if f.title.Focused() {
				f.title.Blur()
				f.desc.Focus()
				return f, textarea.Blink
			}
			// Return the completed form as a message.
			return board.Update(f)
		}
	}
	if f.title.Focused() {
		f.title, cmd = f.title.Update(msg)
		return f, cmd
	}
	f.desc, cmd = f.desc.Update(msg)
	return f, cmd
}

func (f Form) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"Create a new task",
		f.title.View(),
		f.desc.View(),
		f.help.View(keys),
	)
}
