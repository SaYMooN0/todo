package structs

import "time"

type Task struct {
	Id          int64
	Name        string
	Info        string
	IsCompleted bool
	HasDeadline bool
	Deadline    time.Time
	Importance  uint
	User        int64
}

func NewTask(id int64, name string, info string) Task {
	return Task{
		Id:          id,
		Name:        name,
		Info:        info,
		HasDeadline: false,
		Deadline:    time.Now(),
	}
}

func NewTaskWithDeadline(id int64, name string, info string, deadline time.Time) Task {
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
