package main

import (
	"firstProject/internal/db"
	"firstProject/internal/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	// Загрузка переменных окружения
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла:", err)
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	dbURL := os.Getenv("POSTGRES_URL")

	// Подключение к PostgreSQL
	err = db.InitDB(dbURL)
	if err != nil {
		log.Fatal("Ошибка подключения к PostgreSQL:", err)
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
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Команды:\n/add <задача> - добавить задачу\n/list - список задач\n/delete <описание> - удалить задачу")
				bot.Send(msg)
			}
		}
	}
}
