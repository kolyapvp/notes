package handlers

import (
	"firstProject/internal/db"
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AddTaskHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	task := strings.TrimSpace(update.Message.CommandArguments())
	if task == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите задачу после команды /add.")
		bot.Send(msg)
		return
	}

	err := db.AddTask(task)
	if err != nil {
		log.Println("Ошибка добавления задачи:", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось добавить задачу.")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Задача добавлена!")
	bot.Send(msg)
}

func ListTasksHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	tasks, err := db.GetTasks()
	if err != nil {
		log.Println("Ошибка получения задач:", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось получить список задач.")
		bot.Send(msg)
		return
	}

	if len(tasks) == 0 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Список задач пуст.")
		bot.Send(msg)
		return
	}

	var response string
	for i, task := range tasks {
		response += fmt.Sprintf("%d. %s\n", i+1, task)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}

func DeleteTaskHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	arg := strings.TrimSpace(update.Message.CommandArguments())
	index, err := strconv.Atoi(arg)
	if err != nil || index <= 0 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите номер задачи для удаления.")
		bot.Send(msg)
		return
	}

	err = db.DeleteTask(index - 1)
	if err != nil {
		log.Println("Ошибка удаления задачи:", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось удалить задачу.")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Задача удалена!")
	bot.Send(msg)
}
