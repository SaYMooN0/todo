package handlers

import (
	"html/template"
	"my-todo-app/src_back/dbutils"
	"my-todo-app/src_back/structs"
	"net/http"
)

func IndexPage(writer http.ResponseWriter, request *http.Request) {
	if !CheckForAuthToken(request) {
		http.Redirect(writer, request, "/authorization", http.StatusFound)
		return
	}
	http.ServeFile(writer, request, "src_front/index.html")
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

func AddTodo(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := request.ParseForm(); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// newTask := src.Task{
	// 	Name:     request.FormValue("name"),
	// 	Info:     request.FormValue("info"),
	// 	Complete: false,
	// }
	// tasksTemplate := template.Must(template.ParseFiles("templates/tasks.go.tmpl"))
	// if err := tasksTemplate.Execute(writer, tasks); err != nil {
	// 	http.Error(writer, err.Error(), http.StatusInternalServerError)
	// }
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
