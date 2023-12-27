package handlers

import (
	"html/template"
	"net/http"
)

func SendErrorInErrorLine(writer http.ResponseWriter, errorText string) {
	tasksTemplate, err := template.ParseFiles("templates/error-line.go.tmpl")
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := tasksTemplate.Execute(writer, errorText); err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}
