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
	CreatedAt   time.Time
}

func NewTaskWithoutDeadline(name, info string, importance int8, user int64) *Task {
	return &Task{
		Name:        name,
		Info:        info,
		IsCompleted: false,
		HasDeadline: false,
		Importance:  importance,
		User:        user,
		CreatedAt:   time.Now(),
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
		CreatedAt:   time.Now(),
	}
}
func (t *Task) DivShadow() string {
	if t.IsCompleted {
		return "completed"
	}
	if t.HasDeadline && time.Now().After(t.Deadline) {
		return "red"
	}
	return "default"
}
