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

const idEncryptionKey = "JJllnnUU"

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
func CheckForAuthToken(writer http.ResponseWriter, request *http.Request) bool {

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