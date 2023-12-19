package main

import (
	src "my-todo-app/src_back"
	"net/http"
	"sync"
	"text/template"
)

var tasks = make([]src.Task, 0)
var mutex = &sync.Mutex{}
var tasksTemplate = template.Must(template.ParseFiles("templates/tasks.go.tmpl"))

func main() {
	http.Handle("/src_front/", http.StripPrefix("/src_front/", http.FileServer(http.Dir("src_front"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "src_front/index.html")
	})
	http.HandleFunc("/add-todo", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		newTask := src.Task{
			Name:     r.FormValue("name"),
			Info:     r.FormValue("dfds"),
			Complete: false,
		}
		mutex.Lock()
		tasks = append(tasks, newTask)
		defer mutex.Unlock()
		if err := tasksTemplate.Execute(w, tasks); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.ListenAndServe(":8080", nil)
}
