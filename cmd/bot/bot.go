package main

import (
	"firstProject/internal/db"
	"firstProject/internal/handlers"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Загрузка переменных окружения
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла:", err)
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	dbConnection := os.Getenv("DATABASE_URL")

	// Подключение к PostgreSQL с попытками повторного подключения
	for i := 0; i < 5; i++ {
		err = db.InitDB(dbConnection)
		if err == nil {
			log.Println("Успешное подключение к БД")
			break
		}
		log.Printf("Попытка %d: не удалось подключиться к БД: %v", i+1, err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatal("Ошибка подключения к БД после нескольких попыток:", err)
	}

	// Создание Telegram-бота
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal("Ошибка создания бота:", err)
	}
	bot.Debug = true
	log.Printf("Бот авторизован под именем %s", bot.Self.UserName)

	// Настройка обработчиков команд
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Command() {
			case "add":
				go handlers.AddTaskHandler(bot, update)
			case "list":
				go handlers.ListTasksHandler(bot, update)
			case "delete":
				go handlers.DeleteTaskHandler(bot, update)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Команды:\n/add <задача> - добавить задачу\n/list - список задач\n/delete <номер> - удалить задачу")
				bot.Send(msg)
			}
		}
	}
}
