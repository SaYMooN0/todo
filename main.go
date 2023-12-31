package main

import (
	"fmt"
	"my-todo-app/src_back/dbutils"
	mailUtils "my-todo-app/src_back/email"
	handlers "my-todo-app/src_back/http_handlers"

	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading .env file: ", err)
		return
	}
	err = initDBFromDotEnv()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = initMailUtilsFromDotEnv()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = initEncryptionKeysFromDotEnv()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbutils.CloseConnection()
	http.Handle("/src_front/", http.StripPrefix("/src_front/", http.FileServer(http.Dir("src_front"))))
	http.HandleFunc("/", handlers.IndexPage)
	http.HandleFunc("/registration", handlers.RegistrationPage)
	http.HandleFunc("/authorization", handlers.AuthorizationPage)
	http.HandleFunc("/passwordRecovery", handlers.PasswordRecoveryPage)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/signup", handlers.SignUp)
	http.HandleFunc("/confirmEmail", handlers.ConfirmEmail)
	http.HandleFunc("/renderTasks", handlers.RenderTasks)
	http.HandleFunc("/newTaskCreated", handlers.NewTaskCreated)
	http.HandleFunc("/newTaskForm", handlers.NewTaskForm)
	http.HandleFunc("/completeTask", handlers.CompleteTask)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("ListenAndServe: ", err)
		return
	}
}

func initDBFromDotEnv() error {
	dbConnStr := os.Getenv("DB_CONN_STR")
	if dbConnStr == "" {
		return fmt.Errorf("DB_CONN_STR is not set in .env file")
	}

	err := dbutils.InitDB(dbConnStr)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	return nil
}

func initMailUtilsFromDotEnv() error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	if host == "" || port == "" || user == "" || password == "" {
		return fmt.Errorf("SMTP configuration is not fully set in .env file")
	}

	mailUtils.InitMailUtils(host, port, user, password)
	return nil
}
func initEncryptionKeysFromDotEnv() error {
	idEncryption := os.Getenv("ID_ENCRYPTION")
	confirmationIdEncryption := os.Getenv("CONFIRMATION_ENCRYPTION")
	if idEncryption == "" || confirmationIdEncryption == "" {
		return fmt.Errorf("encryption keys are not fully set in .env file")
	}
	handlers.InitEncryptionKeys(idEncryption, confirmationIdEncryption)
	return nil
}
