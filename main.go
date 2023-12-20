package main

import (
	"database/sql"
	"fmt"
	src "my-todo-app/src_back"
	"my-todo-app/src_back/dbutils"
	handlers "my-todo-app/src_back/http_handlers"

	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	tasks = make([]src.Task, 0)
	mutex = &sync.Mutex{}
	db    *sql.DB
)

func main() {
	var err error
	db, err = initDBFromDotEnv()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	http.Handle("/src_front/", http.StripPrefix("/src_front/", http.FileServer(http.Dir("src_front"))))
	http.HandleFunc("/", handlers.IndexPage)
	http.HandleFunc("/registration", handlers.RegistrationPage)
	http.HandleFunc("/add-todo", handlers.AddTodo)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("ListenAndServe: ", err)
		return
	}
}

func initDBFromDotEnv() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	dbConnStr := os.Getenv("DB_CONN_STR")
	if dbConnStr == "" {
		return nil, fmt.Errorf("DB_CONN_STR is not set in .env file")
	}

	db, err := dbutils.InitDB(dbConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return db, nil
}
