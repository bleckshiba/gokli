package main

import (
	"github.com/charmbracelet/bubbles/list"
)

var titles = []string{
	"TODO", "WIP", "DONE",
}

func (board *Board) initList() {
	board.columns = []Column{
		newColumn(PENDING),
		newColumn(WIP),
		newColumn(DONE),
	}

	for t := range titles {
		board.columns[t].list.Title = titles[t]
		board.columns[t].list.SetItems([]list.Item{})
	}

	tasks := taskRepo.FetchTasks()
	for t := range tasks {
		board.columns[tasks[t].Status].list.InsertItem(
			APPEND,
			tasks[t],
		)
	}
}
