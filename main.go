package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var board *Board
var taskRepo TaskRepo

func main() {
  taskRepo.Init()

	form, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer form.Close()

	board = NewBoard()
	board.initList()
	program := tea.NewProgram(board)
	if _, err := program.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
