package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	MARGIN = 4
)

type Board struct {
	help    help.Model
	loaded  bool
	inFocus Status
	columns []Column
	quiting bool
}

func NewBoard() *Board {
	help := help.New()
	help.ShowAll = true
	return &Board{help: help, inFocus: PENDING}
}

func (board *Board) Init() tea.Cmd {
	return nil
}

func (board *Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		var cmds []tea.Cmd
		board.help.Width = msg.Width - MARGIN
		for i := range board.columns {
			var model tea.Model
			model, cmd = board.columns[i].Update(msg)
			board.columns[i] = model.(Column)
			cmds = append(cmds, cmd)
		}
		board.loaded = true
		return board, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Exit):
			board.quiting = true
			return board, tea.Quit
		case key.Matches(msg, keys.Left):
			board.columns[board.inFocus].Blur()
			board.inFocus = board.inFocus.getPrev()
			board.columns[board.inFocus].Focus()
		case key.Matches(msg, keys.Right):
			board.columns[board.inFocus].Blur()
			board.inFocus = board.inFocus.getNext()
			board.columns[board.inFocus].Focus()
		}
	case moveMsg:
		taskRepo.UpdateTask(msg.Task)
		return board, board.columns[board.inFocus.getNext()].Set(APPEND, msg.Task)
	case Form:
		task := msg.AddTask()
		taskRepo.InsertTask(task)
		return board, board.columns[board.inFocus].Set(msg.index, task)
	}
	model, cmd := board.columns[board.inFocus].Update(msg)
	if _, ok := model.(Column); ok {
		board.columns[board.inFocus] = model.(Column)
	} else {
		return model, cmd
	}
	return board, cmd
}

func (board *Board) View() string {
	if board.quiting {
		return ""
	}

	if !board.loaded {
		return "Initializing Data..."
	}

	view := lipgloss.JoinHorizontal(
		lipgloss.Left,
		board.columns[PENDING].View(),
		board.columns[WIP].View(),
		board.columns[DONE].View(),
	)

	return lipgloss.JoinVertical(lipgloss.Left, view, board.help.View(keys))
}
