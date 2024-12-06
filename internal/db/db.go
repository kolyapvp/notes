package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB(conn string) error {
	var err error
	db, err = sql.Open("postgres", conn)
	if err != nil {
		return err
	}
	return db.Ping()
}

func AddTask(task string) error {
	_, err := db.Exec("INSERT INTO tasks (description) VALUES ($1)", task)
	return err
}

func GetTasks() ([]string, error) {
	rows, err := db.Query("SELECT description FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []string
	for rows.Next() {
		var task string
		if err := rows.Scan(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func DeleteTask(index int) error {
	_, err := db.Exec("DELETE FROM tasks WHERE id = $1", index+1)
	return err
}
