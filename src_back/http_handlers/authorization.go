package handlers

import (
	"net/http"
)

func AuthorizationPage(writer http.ResponseWriter, request *http.Request) {
	if CheckForAuthToken(writer, request) {
		http.Redirect(writer, request, "/index", http.StatusFound)
		return
	}
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
			SendErrorInErrorLine(writer, "Incorrect data")
		}
	}
}
