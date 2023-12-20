package handlers

import (
	"net/http"
)

func IndexPage(writer http.ResponseWriter, request *http.Request) {
	_, err := CheckForCookies(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/registration", http.StatusFound)
		return
	}
	http.ServeFile(writer, request, "src_front/index.html")
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
