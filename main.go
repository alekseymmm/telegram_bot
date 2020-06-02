package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func getName(u *tgbotapi.User) string {
	res := ""
	if u == nil {
		return res
	}
	if u.FirstName != "" {
		res += u.FirstName + " "
	}
	if u.LastName != "" {
		res += u.LastName + " "
	}
	res += "(@" + u.UserName + ")"
	return res
}

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
	username := msg.CommandArguments()
	printedName := getName(msg.From)
	if username == "" {
		reply := printedName + ", no @username in your vote. Add someone.\n"
		replyMsg := tgbotapi.NewMessage(msg.Chat.ID, reply)
		bot.Send(replyMsg)
		return
	}
	key := myUserName + "_" + username[1:]

	if _, ok := votes[key]; ok {
		reply := fmt.Sprintln("Calm down,", printedName, " you have already voted this!")
		replyMsg := tgbotapi.NewMessage(msg.Chat.ID, reply)
		bot.Send(replyMsg)
		return
	}

	//reply := fmt.Sprintln(msg.From.FirstName, msg.From.LastName, "voted for someone!")
	reply := fmt.Sprintln(printedName, "voted for someone!")
	replyMsg := tgbotapi.NewMessage(msg.Chat.ID, reply)
	bot.Send(replyMsg)

	votes[key] = username[1:]
	log.Printf("Votes : %s", votes)
}

func pickVotedName(votes map[string]string) string {
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

	return votedUserName
}

func pullCmd(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, votes map[string]string) {
	myUserName := msg.From.UserName

	if len(votes) == 0 {
		reply := "No votes yet, you may suggest someone to fire.\n"
		replyMsg := tgbotapi.NewMessage(msg.Chat.ID, reply)
		bot.Send(replyMsg)
		return
	}

	arg := msg.CommandArguments()
	cnt := 0
	switch arg {
	case "all":
		cnt = len(votes)
	case "":
		cnt = 1
	default:
		val, err := strconv.Atoi(arg)
		if err != nil {
			replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "Wrong  /pull comand format")
			bot.Send(replyMsg)
			return
		}
		if val > len(votes) {
			cnt = len(votes)
		} else {
			cnt = val
		}
	}

	printedName := getName(msg.From)
	for i := 0; i < cnt; i++ {
		votedUserName := pickVotedName(votes)
		log.Printf("i=%d myUserName=%s, votedUserName=%s", i, myUserName, votedUserName)
		log.Printf("Votes : %s", votes)
		if myUserName == votedUserName {
			reply := fmt.Sprintln("Sorry,", printedName, "but you are fired!")
			if len(votes) > 0 {
				reply += "Does someone else feel lucky?\n"
			} else {
				reply += "No names left...\n"
			}
			deleteFiredUser(votes, myUserName)
			replyMsg := tgbotapi.NewMessage(msg.Chat.ID, reply)
			bot.Send(replyMsg)
			return
		}
	}
	replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "Lucky you "+printedName+". You are not fired.")
	bot.Send(replyMsg)
}

func countCmd(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, votes map[string]string) {
	log.Printf("Votes : %s", votes)

	reply := fmt.Sprintln("Names left: ", len(votes))
	replyMsg := tgbotapi.NewMessage(msg.Chat.ID, reply)
	bot.Send(replyMsg)
}

func helpCmd(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	reply := `This is hr_santa bot. You can propose the name to fire someone.

			 Run /vote @username . You can vote multiple times, but for differnet persons.

			 If you feel lucky try and pull some names. If you get your own name then you are fired.
			 Run /pull or /pull <number> to try <number> times or /pull all   if you are crazy.
			 
			 Run /count to see how many names are available to /pull.`
	replyMsg := tgbotapi.NewMessage(msg.Chat.ID, reply)
	bot.Send(replyMsg)
}

func main() {
	// подключаемся к боту с помощью токена
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	votes := make(map[string]string)
	chatVotesMap := make(map[int64]map[string]string)
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
		chatID := msg.Chat.ID
		if val, ok := chatVotesMap[chatID]; ok {
			votes = val
			log.Printf("Found votes map for chatID=%d\n", chatID)
		} else {
			votes = make(map[string]string)
			chatVotesMap[chatID] = votes
			log.Printf("Create votes map for chatID=%d\n", chatID)
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
		case "help":
			helpCmd(bot, msg)
		}

	}
}
