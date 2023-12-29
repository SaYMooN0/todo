package forms

import "time"

type TaskForm struct {
	Name        string
	Info        string
	Importance  string
	HasDeadline string
	Deadline    string
	ErrorLine   string
}

func NewTaskForm(name, info, importance, hasDeadline, deadline, errorLine string) *TaskForm {
	return &TaskForm{
		Name:        name,
		Info:        info,
		Importance:  importance,
		HasDeadline: hasDeadline,
		Deadline:    deadline,
		ErrorLine:   errorLine,
	}
}
func DefaultTaskForm() *TaskForm {
	return &TaskForm{
		Name:        "",
		Info:        "",
		Importance:  "1",
		HasDeadline: "true",
		Deadline:    time.Now().Add(24 * time.Hour).Format("2006-01-02"),
		ErrorLine:   "",
	}
}
