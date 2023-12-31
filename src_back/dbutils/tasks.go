package dbutils

import (
	"database/sql"
	"fmt"
	structs "my-todo-app/src_back/structs"

	_ "github.com/lib/pq"
)

func GetTasksForUserID(userID int64) ([]structs.Task, error) {
	tasks := []structs.Task{}
	rows, err := db.Query("SELECT id, name, info, is_completed, has_deadline, deadline, importance, user_id FROM tasks WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t structs.Task
		if err := rows.Scan(&t.Id, &t.Name, &t.Info, &t.IsCompleted, &t.HasDeadline, &t.Deadline, &t.Importance, &t.User); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
func AddTask(t *structs.Task) (int64, error) {
	var taskID int64
	err := db.QueryRow(`
		INSERT INTO tasks(name, info, is_completed, has_deadline, deadline, importance, user_id, created_at) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		t.Name, t.Info, t.IsCompleted, t.HasDeadline, t.Deadline, t.Importance, t.User, t.CreatedAt,
	).Scan(&taskID)
	if err != nil {
		return 0, err
	}
	return taskID, nil
}
func CompleteTask(taskID int64) error {
	_, err := db.Exec("UPDATE tasks SET is_completed = true WHERE id = $1", taskID)
	if err != nil {
		return fmt.Errorf("CompleteTask: %v", err)
	}
	return nil
}
func GetTaskByID(taskID int64) (*structs.Task, error) {
	var t structs.Task
	query := `SELECT id, name, info, is_completed, has_deadline, deadline, importance, user_id, created_at FROM tasks WHERE id = $1`
	row := db.QueryRow(query, taskID)
	err := row.Scan(&t.Id, &t.Name, &t.Info, &t.IsCompleted, &t.HasDeadline, &t.Deadline, &t.Importance, &t.User, &t.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no task with ID %d", taskID)
		}
		return nil, err
	}
	return &t, nil
}
