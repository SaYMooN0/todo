package handlers

import (
	"fmt"
	"html/template"
	"my-todo-app/src_back/dbutils"
	forms "my-todo-app/src_back/form_structs"
	"my-todo-app/src_back/structs"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func IndexPage(writer http.ResponseWriter, request *http.Request) {
	if !CheckForAuthToken(request) {
		http.Redirect(writer, request, "/authorization", http.StatusFound)
		return
	}
	http.ServeFile(writer, request, "src_front/index.html")
}
func CompleteTask(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(writer, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := request.ParseForm(); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	taskID := request.FormValue("id")
	if taskID == "" {
		http.Error(writer, "ID not provided", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(writer, "Task %s completed", taskID)
}

func RenderTasks(writer http.ResponseWriter, request *http.Request) {
	id, err := GetUserIdFromCookie(request)
	if err != nil {
		SendErrorInErrorLine(writer, "Please, log out and log in one more time")
	}
	tasks, err := dbutils.GetTasksForUserID(id)
	if err != nil {
		SendErrorInErrorLine(writer, "An error has occurred, please try again later")
	}
	if len(tasks) < 1 {
		getNoTasksDiv(writer)

	} else {
		getTasksList(writer, tasks)
	}
}
func NewTaskForm(writer http.ResponseWriter, request *http.Request) {
	sendNewTaskForm(writer, *forms.DefaultTaskForm())
}
func NewTaskCreated(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := request.ParseForm(); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	taskForm := forms.NewTaskForm(request.FormValue("Name"),
		request.FormValue("Info"),
		request.FormValue("Importance"),
		request.FormValue("HasDeadline"),
		request.FormValue("Deadline"), "")
	id, err := GetUserIdFromCookie(request)
	if err != nil {
		taskForm.ErrorLine = "Please, log out and log in one more time"
		sendNewTaskForm(writer, *taskForm)
		return
	}
	task, errString := formNewTask(
		request.FormValue("Name"),
		request.FormValue("Info"),
		request.FormValue("Importance"),
		request.FormValue("HasDeadline"),
		request.FormValue("Deadline"), id)
	if errString != "" {
		taskForm.ErrorLine = errString
		sendNewTaskForm(writer, *taskForm)
		return
	}
	_, err = dbutils.AddTask(task)
	if err != nil {
		taskForm.ErrorLine = "Server error. Please try again later"
		sendNewTaskForm(writer, *taskForm)
		return
	}
	writer.Header().Set("HX-Redirect", "/index")
	writer.WriteHeader(http.StatusOK)
}
func formNewTask(name, info, importanceStr, hasDeadline, deadlineStr string, userId int64) (*structs.Task, string) {
	if strings.TrimSpace(name) == "" {
		return nil, "Task name cannot be empty"
	}

	var deadline time.Time
	var err error

	hasDeadlineBool := false
	if hasDeadline == "on" {
		hasDeadlineBool = true
	}

	if hasDeadlineBool {
		deadline, err = time.Parse("2006-01-02", deadlineStr)
		if err != nil {
			return nil, "Invalid deadline format"
		}

		if deadline.Before(time.Now()) {
			return nil, "Deadline must be in the future"
		}
	}

	importance, err := strconv.ParseInt(importanceStr, 10, 8)
	if err != nil {
		return nil, "Invalid value for 'Importance'"
	}

	var task *structs.Task
	if hasDeadlineBool {
		task = structs.NewTaskWithDeadline(name, info, deadline, int8(importance), userId)
	} else {
		task = structs.NewTaskWithoutDeadline(name, info, int8(importance), userId)
	}

	return task, ""
}
func getTasksList(writer http.ResponseWriter, tasks []structs.Task) {
	tasksTemplate, err := template.ParseFiles("templates/tasks.go.tmpl")
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := tasksTemplate.Execute(writer, tasks); err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}
func getNoTasksDiv(writer http.ResponseWriter) {
	tasksTemplate, err := template.ParseFiles("templates/no-tasks.go.tmpl")
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := tasksTemplate.Execute(writer, ""); err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}
func sendNewTaskForm(writer http.ResponseWriter, formStruct forms.TaskForm) {
	tasksTemplate, err := template.ParseFiles("templates/new-task-form.go.tmpl")
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := tasksTemplate.Execute(writer, formStruct); err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}
