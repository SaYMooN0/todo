package handlers

import (
	"net/http"
)

func RegistrationPage(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "src_front/pages/registration/registration.html")
}
