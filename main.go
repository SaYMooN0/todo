package main

import (
	"fmt"
	"my-todo-app/src_back/dbutils"
	handlers "my-todo-app/src_back/http_handlers"

	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	var err error
	err = initDBFromDotEnv()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbutils.CloseConnection()
	http.Handle("/src_front/", http.StripPrefix("/src_front/", http.FileServer(http.Dir("src_front"))))
	http.HandleFunc("/", handlers.IndexPage)
	http.HandleFunc("/registration", handlers.RegistrationPage)
	http.HandleFunc("/authorization", handlers.AuthorizationPage)
	http.HandleFunc("/password_recovery", handlers.PasswordRecoveryPage)
	http.HandleFunc("/add-todo", handlers.AddTodo)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/signup", handlers.SignUp)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("ListenAndServe: ", err)
		return
	}
}

func initDBFromDotEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	dbConnStr := os.Getenv("DB_CONN_STR")
	if dbConnStr == "" {
		return fmt.Errorf("DB_CONN_STR is not set in .env file")
	}

	err = dbutils.InitDB(dbConnStr)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	return nil
}
