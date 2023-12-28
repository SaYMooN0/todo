package structs

import "time"

type Task struct {
	Id          int64
	Name        string
	Info        string
	IsCompleted bool
	HasDeadline bool
	Deadline    time.Time
	Importance  int8
	User        int64
}

func NewTaskWithoutDeadline(name, info string, importance int8, user int64) *Task {
	return &Task{
		Name:        name,
		Info:        info,
		IsCompleted: false,
		HasDeadline: false,
		Importance:  importance,
		User:        user,
	}
}
func NewTaskWithDeadline(name, info string, deadline time.Time, importance int8, user int64) *Task {
	return &Task{
		Name:        name,
		Info:        info,
		IsCompleted: false,
		HasDeadline: true,
		Deadline:    deadline,
		Importance:  importance,
		User:        user,
	}
}

func (t *Task) SetImportance(importance int8) {
	if importance < 1 {
		t.Importance = 1
	} else if importance > 10 {
		t.Importance = 10
	} else {
		t.Importance = importance
	}
}
