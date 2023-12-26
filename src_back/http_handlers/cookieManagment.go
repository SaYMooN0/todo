package handlers

import (
	"fmt"
	src "my-todo-app/src_back"
	"my-todo-app/src_back/dbutils"
	"my-todo-app/src_back/structs"
	"net/http"
	"strconv"
	"time"
)

var idEncryptionKey string = ""
var confirmationIdEncryptionKey string = ""

func InitEncryptionKeys(idEncryption, confirmationIdEncryption string) {
	idEncryptionKey = idEncryption
	confirmationIdEncryptionKey = confirmationIdEncryption
}
func GetConfirmationIdFromCookie(request *http.Request) (int64, error) {
	cookie, err := request.Cookie("confirmationId")
	if err != nil {
		return 0, err
	}

	decryptedId, err := src.Decrypt(cookie.Value, confirmationIdEncryptionKey)
	if err != nil {
		return 0, err
	}

	confirmationId, err := strconv.ParseInt(decryptedId, 10, 64)
	if err != nil {
		return 0, err
	}

	return confirmationId, nil
}

func GenerateConfirmationIdToken(id int64) (string, error) {
	encryptedId, err := src.Encrypt(fmt.Sprint(id), confirmationIdEncryptionKey)
	if err != nil {
		return "", err
	}
	return encryptedId, nil
}
func SetConfirmationIdCookie(writer http.ResponseWriter, idToken string) {
	http.SetCookie(writer, &http.Cookie{
		Name:  "confirmationId",
		Value: idToken,
		Path:  "/",
	})
}
func GetUserIdFromCookie(request *http.Request) (int64, error) {
	idCookie, err := request.Cookie("userId")
	if err != nil {
		return 0, err
	}
	decryptedId, err := src.Decrypt(idCookie.Value, idEncryptionKey)
	if err != nil {
		return 0, err
	}

	userId, err := strconv.ParseInt(decryptedId, 10, 64)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func GenerateIDToken(id int64) (string, error) {
	encryptedId, err := src.Encrypt(fmt.Sprint(id), idEncryptionKey)
	if err != nil {
		return "", err
	}
	return encryptedId, nil
}
func SetIDTokenCookie(writer http.ResponseWriter, idToken string) {
	http.SetCookie(writer, &http.Cookie{
		Name:  "userId",
		Value: idToken,
		Path:  "/",
	})
}
func GenerateAuthTokenFromUser(user structs.User) (string, error) {
	return GenerateAuthToken(user.RegistrationDate, user.Email, user.Id)
}

func GenerateAuthToken(registrationDate time.Time, email string, id int64) (string, error) {
	formattedDate := registrationDate.Format("20060102150405")
	authToken := fmt.Sprintf("%s|%s", formattedDate, email)
	encryptedToken, err := src.Encrypt(authToken, fmt.Sprint(id, idEncryptionKey))
	if err != nil {
		return "", err
	}
	return encryptedToken, nil
}

func SetAuthTokenCookie(writer http.ResponseWriter, authToken string) {
	http.SetCookie(writer, &http.Cookie{
		Name:  "authToken",
		Value: authToken,
		Path:  "/",
	})
}

func CheckForAuthToken(request *http.Request) bool {

	idCookie, err := request.Cookie("userId")
	if err != nil {
		return false
	}
	decryptedId, err := src.Decrypt(idCookie.Value, idEncryptionKey)
	if err != nil {
		return false
	}

	userId, err := strconv.Atoi(decryptedId)
	if err != nil {
		return false
	}
	authString, err := dbutils.GetAuthTokenByID(userId)
	if err != nil {
		return false
	}
	authTokenCookie, err := request.Cookie("authToken")
	if err != nil {
		return false
	}
	decryptedToken, err := src.Decrypt(authTokenCookie.Value, fmt.Sprint(userId, idEncryptionKey))
	if err != nil {
		return false
	}

	return authString == decryptedToken
}
func CheckForConfirmationId(request *http.Request) bool {
	_, err := request.Cookie("confirmationId")
	return err == nil
}
func DeleteConfirmationIDCookie(writer http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:    "confirmationId",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	}
	http.SetCookie(writer, cookie)
}
