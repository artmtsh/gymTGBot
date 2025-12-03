package main

import (
	"github.com/artmtsh/gymTGBot/gym"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strings"
)

func LoadToken(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("cannot read token file: %v", err)
	}
	return strings.TrimSpace(string(data))
}

func main() {
	token := LoadToken("token.txt")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	app := gym.NewApp()
	updateConfig := tgbotapi.NewUpdate(0)
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, app.HandleText(update.Message.From.ID, update.Message.Text))
		msg.ReplyToMessageID = update.Message.MessageID
		if _, err = bot.Send(msg); err != nil {
			panic(err)
		}
	}
}
