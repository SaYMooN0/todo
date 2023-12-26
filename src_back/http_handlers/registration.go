package handlers

import (
	"fmt"
	"html/template"
	"math/rand"
	"my-todo-app/src_back/dbutils"
	mailUtils "my-todo-app/src_back/email"
	"my-todo-app/src_back/structs"
	"net/http"
	"regexp"
)

func RegistrationPage(writer http.ResponseWriter, request *http.Request) {
	if CheckForAuthToken(request) {
		http.Redirect(writer, request, "/index", http.StatusFound)
		return
	}

	if CheckForConfirmationId(request) {
		DeleteConfirmationIDCookie(writer)
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

	formData := structs.NewRegistrationForm(request.FormValue("Name"), request.FormValue("Email"), request.FormValue("Password"), request.FormValue("RepeatedPassword"))
	_, err := dbutils.GetIDByEmail(formData.Email)
	if err == nil {
		formData.ErrorLine = "The email has already been registered"
		sendRegistrationFormBack(writer, *formData)
		return
	}
	if message := validateSignUpForm(formData.Email, formData.Password, formData.RepeatedPassword, formData.Name); message != "" {
		formData.ErrorLine = message
		sendRegistrationFormBack(writer, *formData)
		return
	}
	emailConfirmation(writer, request, *formData)
}

func emailConfirmation(writer http.ResponseWriter, request *http.Request, formData structs.RegistrationForm) {

	user := structs.User{Name: formData.Name, Email: formData.Email, Password: formData.Password}
	confirmationCode := generateConfirmationCode()
	confirmationId, err := dbutils.UpsertEmailConfirmation(confirmationCode, user)
	if err != nil {
		fmt.Print("New user adding error: ", err)
		formData.ErrorLine = "Server error"
		sendRegistrationFormBack(writer, formData)
		return
	}
	confirmationToken, err := GenerateConfirmationIdToken(confirmationId)
	if err != nil {
		formData.ErrorLine = "Server error"
		sendRegistrationFormBack(writer, formData)
		return
	}
	err = mailUtils.SendConfirmationCode(formData.Email, confirmationCode)
	if err != nil {
		formData.ErrorLine = "Server error"
		sendRegistrationFormBack(writer, formData)
		return
	}
	SetConfirmationIdCookie(writer, confirmationToken)
	goToEmailConfirmation(writer)

}
func ConfirmEmail(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		return
	}
	if err := request.ParseForm(); err != nil {
		http.Error(writer, "Form parsing error", http.StatusBadRequest)
		return
	}
	confirmationId, err := GetConfirmationIdFromCookie(request)
	if err != nil {
		SendErrorInErrorLine(writer, "Server error")
		return
	}
	code, err := dbutils.GetConfirmationCodeByID(confirmationId)
	if err != nil {
		SendErrorInErrorLine(writer, "Server error")
		return
	}
	receivedCode := request.FormValue("Code")
	if receivedCode == "" {
		SendErrorInErrorLine(writer, "No confirmation code recieved")
		return
	}
	if string(code) != receivedCode {
		SendErrorInErrorLine(writer, "Wrong code")
		return
	}

	user, err := dbutils.GetUserFromConfirmationTable(confirmationId)
	if err != nil {
		SendErrorInErrorLine(writer, "Server error. Please try again later")
		return
	}
	userId, err := dbutils.AddNewUserToDb(user)
	if err != nil {
		fmt.Print("New user adding error: ", err)
		SendErrorInErrorLine(writer, "Server error. Please try again later")
		return
	}
	authToken, err := GenerateAuthToken(user.RegistrationDate, user.Email, userId)
	if err != nil {
		SendErrorInErrorLine(writer, "Server error. Please try again later")
		return
	}
	idToken, err := GenerateIDToken(userId)
	if err != nil {
		SendErrorInErrorLine(writer, "Server error. Please try again later")
		return
	}
	DeleteConfirmationIDCookie(writer)
	SetIDTokenCookie(writer, idToken)
	SetAuthTokenCookie(writer, authToken)
	dbutils.DeleteUserFromConfirmationTable(user.Id)
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

func sendRegistrationFormBack(writer http.ResponseWriter, formStruct structs.RegistrationForm) {
	tasksTemplate, err := template.ParseFiles("templates/registration-form.go.tmpl")
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := tasksTemplate.Execute(writer, formStruct); err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}
func goToEmailConfirmation(writer http.ResponseWriter) {
	tasksTemplate, err := template.ParseFiles("templates/email-confirmation.go.tmpl")
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := tasksTemplate.Execute(writer, ""); err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}
func generateConfirmationCode() string {
	const digits = "0123456789"
	const length = 12
	code := make([]byte, length)

	for i := range code {
		code[i] = digits[rand.Intn(len(digits))]
	}
	return string(code)
}
