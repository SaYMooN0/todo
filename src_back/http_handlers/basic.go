package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"my-todo-app/src_back/dbutils"
	"net/http"
	"strconv"
)

var db *sql.DB

func CheckForAuthToken(writer http.ResponseWriter, request *http.Request) bool {
	idCookie, err := request.Cookie("userId")
	if err != nil {
		return false
	}

	userId, err := strconv.Atoi(idCookie.Value)
	if err != nil {
		return false
	}

	authString, err := dbutils.GetAuthTokenByID(userId)
	if err != nil {
		return false
	}

	authTokenCookie, err := request.Cookie("authToken")
	if err != nil || authString != authTokenCookie.Value {
		return false
	}

	return true
}
func SetAuthToken(writer http.ResponseWriter, id int64, email string, registrationDate string) error {

	authToken := fmt.Sprintf("%s | %s", registrationDate, email)
	http.SetCookie(writer, &http.Cookie{
		Name:  "userId",
		Value: fmt.Sprint(id),
		Path:  "/",
	})
	http.SetCookie(writer, &http.Cookie{
		Name:  "authToken",
		Value: authToken,
		Path:  "/",
	})

	return nil
}
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
