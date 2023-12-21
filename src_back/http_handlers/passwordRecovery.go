package handlers

import (
	"net/http"
)

func PasswordRecoveryPage(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "src_front/pages/password_recovery/password_recovery.html")
}
