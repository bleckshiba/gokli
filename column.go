package main

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	APPEND = -1
)

type Column struct {
	list   list.Model
	status Status
	focus  bool
	width  int
	height int
}

type moveMsg struct {
	Task
}

func (col *Column) Focus() {
	col.focus = true
}

func (col *Column) Blur() {
	col.focus = false
}

func (col *Column) HasFocus() bool {
	return col.focus
}

func newColumn(status Status) Column {
	var focus bool
	if status == PENDING {
		focus = true
	}

	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	return Column{focus: focus, status: status, list: defaultList}
}

func (col *Column) setWidth(width int) {
	col.width = width / MARGIN
}

var defaultStyle lipgloss.Style = lipgloss.NewStyle().Padding(1, 2)

func (col *Column) getStyle() lipgloss.Style {
	if col.HasFocus() {
		return defaultStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Height(col.height).
			Width(col.width)
	}
	return defaultStyle.
		Border(lipgloss.HiddenBorder()).
		Height(col.height).
		Width(col.width)
}

func (col *Column) Set(i int, t Task) tea.Cmd {
	log.Printf("column: %+v", col.list)
	if i != APPEND {
		return col.list.SetItem(i, t)
	}
	return col.list.InsertItem(APPEND, t)
}

func (col *Column) Delete(t Task) tea.Cmd {
	if len(col.list.VisibleItems()) > 0 {
		col.list.RemoveItem(col.list.Index())
	}

	taskRepo.DeleteTask(t)

	var cmd tea.Cmd
	col.list, cmd = col.list.Update(nil)
	return cmd
}

func (col *Column) MoveToNext() tea.Cmd {
	var task Task
	var ok bool

	if task, ok = col.list.SelectedItem().(Task); !ok {
		return nil
	}

	col.list.RemoveItem(col.list.Index())
	task.Status = col.status.getNext()

	var cmd tea.Cmd
	col.list, cmd = col.list.Update(nil)

	return tea.Sequence(cmd, func() tea.Msg { return moveMsg{task} })
}

func (col Column) Init() tea.Cmd {
	return nil
}

func (col Column) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		col.setWidth(msg.Width)
		col.list.SetSize(msg.Width/MARGIN, msg.Height/2)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Edit):
			if len(col.list.VisibleItems()) != 0 {
				task := col.list.SelectedItem().(Task)
				form := NewForm(task.Title(), task.Desc)
				form.index = col.list.Index()
				form.column = col
				return form.Update(nil)
			}
		case key.Matches(msg, keys.New):
			form := newDefaultForm()
			form.index = APPEND
			form.column = col
			return form.Update(nil)
		case key.Matches(msg, keys.Delete):
			if len(col.list.Items()) > 0 {
				return col, col.Delete(col.list.SelectedItem().(Task))
			}
		case key.Matches(msg, keys.Select):
			return col, col.MoveToNext()
		}
	}
	col.list, cmd = col.list.Update(msg)
	return col, cmd
}

func (col Column) View() string {
	return col.getStyle().Render(col.list.View())
}
