package handlers

import (
	"my-todo-app/src_back/dbutils"
	"net/http"
)

func AuthorizationPage(writer http.ResponseWriter, request *http.Request) {
	if CheckForAuthToken(request) {
		http.Redirect(writer, request, "/index", http.StatusFound)
		return
	}
	http.ServeFile(writer, request, "src_front/pages/authorization/authorization.html")
}
func Login(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		return
	}

	if err := request.ParseForm(); err != nil {
		http.Error(writer, "Form parsing error", http.StatusBadRequest)
		return
	}
	email, password := request.FormValue("Email"), request.FormValue("Password")
	if email == "" || password == "" {
		SendErrorInErrorLine(writer, "Fill all fields")
		return
	}
	user, err := dbutils.AuthenticateUser(email, password)
	if err != nil {
		SendErrorInErrorLine(writer, err.Error())
		return
	}
	authToken, err := GenerateAuthTokenFromUser(user)
	if err != nil {
		SendErrorInErrorLine(writer, "Server error")
		return
	}
	idToken, err := GenerateIDToken(user.Id)
	if err != nil {
		SendErrorInErrorLine(writer, "Server error")
		return
	}
	SetAuthTokenCookie(writer, authToken)
	SetIDTokenCookie(writer, idToken)
	writer.Header().Set("HX-Redirect", "/index")
	writer.WriteHeader(http.StatusOK)
}
