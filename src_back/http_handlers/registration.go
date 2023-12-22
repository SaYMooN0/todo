package handlers

import (
	"fmt"
	src "my-todo-app/src_back"
	"my-todo-app/src_back/dbutils"
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
	if request.Method == "POST" {
		err := request.ParseForm()
		if err != nil {
			http.Error(writer, "Form parsing error", http.StatusBadRequest)
			return
		}

		email := request.FormValue("Email")
		password := request.FormValue("Password")
		repeatedPassword := request.FormValue("RepeatedPassword")
		userName := request.FormValue("Name")
		if email == "" || password == "" || repeatedPassword == "" || userName == "" {
			SendErrorInErrorLine(writer, "Fill all fields")
			return
		}
		emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
		if !emailRegex.MatchString(email) {
			SendErrorInErrorLine(writer, "Invalid email format")
			return
		}

		_, err = dbutils.GetIDByEmail(email)
		if err == nil {
			SendErrorInErrorLine(writer, "This email is already in use")
			return
		}

		if len(password) < 8 || len(password) > 30 {
			SendErrorInErrorLine(writer, "Password must be between 8 and 30 characters long")
			return
		}

		if password != repeatedPassword {
			SendErrorInErrorLine(writer, "Passwords don't match")
			return
		}

		if len(userName) < 3 || len(userName) > 30 {
			SendErrorInErrorLine(writer, "Name must be between 3 and 30 characters long")
			return
		}
		nameRegex := regexp.MustCompile(`^[a-zA-Z0-9 _<>\[\]{}]+$`)
		if !nameRegex.MatchString(userName) {
			SendErrorInErrorLine(writer, "Name contains invalid characters")
			return
		}
		var registrationDate = time.Now()
		newUser := src.User{
			Name:             userName,
			Email:            email,
			Password:         password,
			RegistrationDate: registrationDate,
		}
		userId, err := dbutils.AddNewUserToDb(newUser)
		if err != nil {
			fmt.Print("New user adding error: ", err)
			SendErrorInErrorLine(writer, "Server error")
			return
		}
		//yyyymmddhhmmss format
		formattedDate := registrationDate.Format("20060102150405")
		SetAuthToken(writer, userId, email, formattedDate)
		writer.Header().Set("HX-Redirect", "/index")
		writer.WriteHeader(http.StatusOK)
	}
}
