package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // PostgreSQL драйвер
)

var db *sql.DB

// InitDB инициализирует подключение к базе данных
func InitDB(conn string) error {
	var err error
	db, err = sql.Open("postgres", conn)
	if err != nil {
		log.Printf("Ошибка при открытии подключения: %v", err)
		return err
	}

	// Проверка соединения с базой
	if err := db.Ping(); err != nil {
		log.Printf("Ошибка при проверке подключения: %v", err)
		return err
	}

	log.Println("Подключение к базе данных установлено успешно")
	return nil
}

// AddTask добавляет новую задачу в таблицу tasks
func AddTask(task string) error {
	_, err := db.Exec("INSERT INTO tasks (description) VALUES ($1)", task)
	if err != nil {
		log.Printf("Ошибка при добавлении задачи: %v", err)
		return err
	}
	log.Printf("Задача '%s' успешно добавлена", task)
	return nil
}

// GetTasks возвращает список всех задач из таблицы tasks
func GetTasks() ([]string, error) {
	rows, err := db.Query("SELECT description FROM tasks")
	if err != nil {
		log.Printf("Ошибка при получении задач: %v", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []string
	for rows.Next() {
		var task string
		if err := rows.Scan(&task); err != nil {
			log.Printf("Ошибка при чтении строки: %v", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Ошибка при обработке результатов: %v", err)
		return nil, err
	}

	log.Printf("Получено задач: %d", len(tasks))
	return tasks, nil
}

// DeleteTask удаляет задачу по ID
func DeleteTask(id int) error {
	// Удаляем задачу с указанным ID
	_, err := db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		log.Printf("Ошибка при удалении задачи с ID %d: %v", id, err)
		return err
	}

	log.Printf("Задача с ID %d успешно удалена", id)
	return nil
}
