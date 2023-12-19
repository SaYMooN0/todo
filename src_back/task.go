package src

import (
	"time"
)

type Task struct {
	Id          uint
	Name        string
	Info        string
	Complete    bool
	HasDeadline bool
	Deadline    time.Time
	Importance  uint
}

func NewTask(id uint, name string, info string) Task {
	return Task{
		Id:          id,
		Name:        name,
		Info:        info,
		HasDeadline: false,
		Deadline:    time.Now(),
	}
}

func NewTaskWithDeadline(id uint, name string, info string, deadline time.Time) Task {
	return Task{
		Id:          id,
		Name:        name,
		Info:        info,
		HasDeadline: true,
		Deadline:    deadline,
	}
}
func (t *Task) SetImportance(importance uint) {
	if importance < 1 {
		t.Importance = 1
	} else if importance > 10 {
		t.Importance = 10
	} else {
		t.Importance = importance
	}
}
