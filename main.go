package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	// подключаемся к боту с помощью токена
	bot, err := tgbotapi.NewBotAPI("1174621006:AAF04nE-Cku5AhRIxPZeLUkEpGedVfwjYD4")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// инициализируем канал, куда будут прилетать обновления от API
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updates, err := bot.GetUpdatesChan(ucfg)
	// читаем обновления из канала
	for update := range updates {
		msg := update.Message
		if msg == nil { // ignore any non-Message updates
			continue
		}

		if !msg.IsCommand() { // ignore any non-command Messages
			continue
		}

		cmd := msg.Command()
		log.Printf("Get command: %s", cmd)
		switch msg.Command() {
		case "vote":
			delMsg := tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID)
			bot.DeleteMessage(delMsg)
			//reply := fmt.Sprintln(msg.From.FirstName, msg.From.LastName, "voted for someone!")
			reply := fmt.Sprintln("@"+msg.From.UserName, "voted for someone!")
			replyMsg := tgbotapi.NewMessage(msg.Chat.ID, reply)
			bot.Send(replyMsg)
			username := msg.CommandArguments()
			log.Printf("Add person: %s", username)
		}

		// // Пользователь, который написал боту
		// UserName := update.Message.From.UserName

		// // ID чата/диалога.
		// // Может быть идентификатором как чата с пользователем
		// // (тогда он равен UserID) так и публичного чата/канала
		// ChatID := update.Message.Chat.ID

		// // Текст сообщения
		// Text := update.Message.Text

		// log.Printf("[%s] %d %s", UserName, ChatID, Text)

		// // Ответим пользователю его же сообщением
		// //reply := Text
		// // Созадаем сообщение
		// delMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
		// bot.DeleteMessage(delMsg)
		//msg := tgbotapi.NewMessage(ChatID, reply)
		// и отправляем его
		//bot.Send(msg)

	}
}
