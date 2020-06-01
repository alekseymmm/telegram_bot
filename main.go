package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func deleteFiredUser(votes map[string]string, username string) {
	for k := range votes {
		if votes[k] == username {
			delete(votes, k)
		}
	}
}

func voteCmd(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, votes map[string]string) {
	delMsg := tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID)
	bot.DeleteMessage(delMsg)

	myUserName := msg.From.UserName
	if _, ok := votes[myUserName]; ok {
		reply := fmt.Sprintln("Calm down, @"+myUserName, " you have already voted!")
		replyMsg := tgbotapi.NewMessage(msg.Chat.ID, reply)
		bot.Send(replyMsg)
		return
	}
	//reply := fmt.Sprintln(msg.From.FirstName, msg.From.LastName, "voted for someone!")
	reply := fmt.Sprintln("@"+msg.From.UserName, "voted for someone!")
	replyMsg := tgbotapi.NewMessage(msg.Chat.ID, reply)
	bot.Send(replyMsg)
	username := msg.CommandArguments()
	votes[msg.From.UserName] = username[1:]
	log.Printf("Votes : %s", votes)
}

func pullCmd(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, votes map[string]string) {
	myUserName := msg.From.UserName

	keys := make([]string, len(votes))
	i := 0
	for k := range votes {
		keys[i] = k
		i++
	}
	ind := rand.Intn(len(keys))
	pick := keys[ind]

	votedUserName := votes[pick]
	delete(votes, pick)
	log.Printf("myUserName=%s, votedUserName=%s", myUserName, votedUserName)
	log.Printf("Votes : %s", votes)
	if myUserName == votedUserName {
		reply := fmt.Sprintln("Sorry, @"+myUserName, "but you are fired!")
		if len(votes) > 0 {
			reply += "Does someone else feel lucky?\n"
		} else {
			reply += "No names left...\n"
		}
		deleteFiredUser(votes, myUserName)
		replyMsg := tgbotapi.NewMessage(msg.Chat.ID, reply)
		bot.Send(replyMsg)
	}
}

func countCmd(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, votes map[string]string) {
	log.Printf("Votes : %s", votes)

	reply := fmt.Sprintln("Names left: ", len(votes))
	replyMsg := tgbotapi.NewMessage(msg.Chat.ID, reply)
	bot.Send(replyMsg)
}

func main() {
	// подключаемся к боту с помощью токена
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	votes := make(map[string]string)

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
			voteCmd(bot, msg, votes)
		case "pull":
			pullCmd(bot, msg, votes)
		case "count":
			countCmd(bot, msg, votes)
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
