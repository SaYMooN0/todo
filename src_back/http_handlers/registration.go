package handlers

import (
	"fmt"
	"my-todo-app/src_back/dbutils"
	"my-todo-app/src_back/structs"
	"net/http"
	"regexp"
	"time"
)

func RegistrationPage(writer http.ResponseWriter, request *http.Request) {
	if CheckForAuthToken(writer, request) {
		http.Redirect(writer, request, "/index", http.StatusFound)
		return
	}
	http.ServeFile(writer, request, "src_front/pages/registration/registration.html")
}
func SignUp(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		return
	}

	if err := request.ParseForm(); err != nil {
		http.Error(writer, "Form parsing error", http.StatusBadRequest)
		return
	}

	email, password, repeatedPassword, userName := request.FormValue("Email"), request.FormValue("Password"), request.FormValue("RepeatedPassword"), request.FormValue("Name")

	if message := validateSignUpForm(email, password, repeatedPassword, userName); message != "" {
		SendErrorInErrorLine(writer, message)
		return
	}

	registrationDate := time.Now()
	newUser := structs.User{Name: userName, Email: email, Password: password, RegistrationDate: registrationDate}

	userId, err := dbutils.AddNewUserToDb(newUser)
	if err != nil {
		fmt.Print("New user adding error: ", err)
		SendErrorInErrorLine(writer, "Server error")
		return
	}
	authToken, err := GenerateAuthToken(registrationDate, email, userId)
	if err != nil {
		SendErrorInErrorLine(writer, "Server error")
		return
	}
	idToken, err := GenerateIDToken(userId)
	if err != nil {
		SendErrorInErrorLine(writer, "Server error")
		return
	}
	SetAuthTokenCookie(writer, authToken)
	SetIDTokenCookie(writer, idToken)
	writer.Header().Set("HX-Redirect", "/index")
	writer.WriteHeader(http.StatusOK)
}

func validateSignUpForm(email, password, repeatedPassword, userName string) string {
	if email == "" || password == "" || repeatedPassword == "" || userName == "" {
		return "Fill all fields"
	}
	if !regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(email) {
		return "Invalid email format"
	}
	if _, err := dbutils.GetIDByEmail(email); err == nil {
		return "This email is already in use"
	}
	if len(password) < 8 || len(password) > 30 {
		return "Password must be between 8 and 30 characters long"
	}
	if password != repeatedPassword {
		return "Passwords don't match"
	}
	if len(userName) < 3 || len(userName) > 30 {
		return "Name must be between 3 and 30 characters long"
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9 _<>\[\]{}]+$`).MatchString(userName) {
		return "Name contains invalid characters"
	}
	return ""
}
