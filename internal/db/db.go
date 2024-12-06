package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // PostgreSQL драйвер
)

var db *sql.DB

// InitDB инициализирует подключение к PostgreSQL
func InitDB(conn string) error {
	var err error
	db, err = sql.Open("postgres", conn)
	if err != nil {
		log.Printf("Ошибка подключения к PostgreSQL: %v", err)
		return err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Ошибка пинга PostgreSQL: %v", err)
		return err
	}

	log.Println("Подключение к PostgreSQL установлено успешно")

	// Создание таблицы задач, если она еще не существует
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		description TEXT NOT NULL
	)`
	_, err = db.Exec(query)
	if err != nil {
		log.Printf("Ошибка создания таблицы: %v", err)
		return err
	}

	return nil
}

// AddTask добавляет новую задачу в таблицу
func AddTask(task string) error {
	query := `INSERT INTO tasks (description) VALUES ($1)`
	_, err := db.Exec(query, task)
	if err != nil {
		log.Printf("Ошибка добавления задачи: %v", err)
		return err
	}
	log.Printf("Задача '%s' успешно добавлена", task)
	return nil
}

// GetTasks возвращает список всех задач
func GetTasks() ([]string, error) {
	rows, err := db.Query(`SELECT description FROM tasks`)
	if err != nil {
		log.Printf("Ошибка получения задач: %v", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []string
	for rows.Next() {
		var task string
		if err := rows.Scan(&task); err != nil {
			log.Printf("Ошибка сканирования строки: %v", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// DeleteTask удаляет задачу по описанию
func DeleteTask(description string) error {
	query := `DELETE FROM tasks WHERE description = $1`
	_, err := db.Exec(query, description)
	if err != nil {
		log.Printf("Ошибка удаления задачи: %v", err)
		return err
	}

	log.Printf("Задача '%s' успешно удалена", description)
	return nil
}
