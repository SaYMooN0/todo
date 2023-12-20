package handlers

import (
	"database/sql"
	"my-todo-app/src_back/dbutils"
	"net/http"
	"strconv"
)

var db *sql.DB

func CheckForCookies(writer http.ResponseWriter, request *http.Request) (*int, error) {
	idCookie, err := request.Cookie("userId")
	if err != nil {
		return nil, err
	}

	userId, err := strconv.Atoi(idCookie.Value)
	if err != nil {
		return nil, err
	}

	authString, err := dbutils.GetAuthTokenByID(db, userId)
	if err != nil {
		return nil, err
	}

	authTokenCookie, err := request.Cookie("authToken")
	if err != nil || authString != authTokenCookie.Value {
		return nil, err
	}

	return &userId, nil
}
