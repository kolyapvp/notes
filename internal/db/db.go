package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var tasksCollection *mongo.Collection

// InitDB инициализирует подключение к MongoDB
func InitDB(conn string) error {
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(conn))
	if err != nil {
		log.Printf("Ошибка создания клиента MongoDB: %v", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Printf("Ошибка подключения к MongoDB: %v", err)
		return err
	}

	// Проверяем соединение
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Ошибка пинга MongoDB: %v", err)
		return err
	}

	// Устанавливаем коллекцию
	tasksCollection = client.Database("tasksdb").Collection("tasks")
	log.Println("Подключение к MongoDB установлено успешно")
	return nil
}

// AddTask добавляет новую задачу в коллекцию
func AddTask(task string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := tasksCollection.InsertOne(ctx, bson.M{"description": task})
	if err != nil {
		log.Printf("Ошибка добавления задачи: %v", err)
		return err
	}
	log.Printf("Задача '%s' успешно добавлена", task)
	return nil
}

// GetTasks возвращает список всех задач
func GetTasks() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := tasksCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Ошибка получения задач: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []string
	for cursor.Next(ctx) {
		var task struct {
			Description string `bson:"description"`
		}
		if err := cursor.Decode(&task); err != nil {
			log.Printf("Ошибка декодирования задачи: %v", err)
			return nil, err
		}
		tasks = append(tasks, task.Description)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Ошибка курсора: %v", err)
		return nil, err
	}

	log.Printf("Получено задач: %d", len(tasks))
	return tasks, nil
}

// DeleteTask удаляет задачу по описанию
func DeleteTask(description string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := tasksCollection.DeleteOne(ctx, bson.M{"description": description})
	if err != nil {
		log.Printf("Ошибка удаления задачи: %v", err)
		return err
	}

	log.Printf("Задача '%s' успешно удалена", description)
	return nil
}
