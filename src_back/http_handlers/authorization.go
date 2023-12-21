package handlers

import (
	"html/template"
	"net/http"
)

func AuthorizationPage(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "src_front/pages/authorization/authorization.html")
}
func Login(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		err := request.ParseForm()
		if err != nil {
			http.Error(writer, "Form parsing error", http.StatusBadRequest)
			return
		}

		// email := request.FormValue("Email")
		// password := request.FormValue("Password")

		if false {
			http.Redirect(writer, request, "/index", http.StatusFound)
		} else {
			tasksTemplate := template.Must(template.ParseFiles("templates/error-line.go.tmpl"))
			if err := tasksTemplate.Execute(writer, "Incorrect data"); err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			} else {
				writer.WriteHeader(http.StatusBadRequest)
			}
		}
	}
}
