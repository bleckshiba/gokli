package main

import "github.com/google/uuid"

type Task struct {
	ID     uuid.UUID
	Status Status
	title  string
	Desc   string
}

func NewTask(title, desc string) Task {
	return Task{Status: PENDING, title: title, Desc: desc}
}

// TODO: to understand what is this
func (task *Task) Next() {
	if task.Status == DONE {
		task.Status = PENDING
	} else {
		task.Status++
	}
}

// list.Item interface
func (task Task) FilterValue() string {
	return task.title
}

func (task Task) Title() string {
	return task.title
}

func (task Task) Description() string {
	return task.Desc
}
