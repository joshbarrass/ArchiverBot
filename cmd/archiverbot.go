package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joshbarrass/ArchiverBot/internal"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Configuration struct {
	BotToken  string `envconfig:"BOT_TOKEN" required:"true"`
	Admin     int    `envconfig:"ADMIN_ID" default:"0"`
	DebugMode bool   `envconfig:"DEBUG_MODE" default:"false"`
}

// based on example from https://github.com/go-telegram-bot-api/telegram-bot-api
func main() {
	var config Configuration
	err := envconfig.Process("AB", &config)
	if err != nil {
		logrus.Fatalf("Failed to process config: %s", err)
	}

	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		logrus.Panic(err)
	}

	bot.Debug = config.DebugMode

	logrus.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if config.Admin == 0 {
		logrus.Warn("No bot admin set! Bot will not function until an admin is set.")
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if config.Admin == 0 {
			replyText := fmt.Sprintf(internal.MessageNoAdmin, update.Message.From.ID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, replyText)
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ParseMode = tgbotapi.ModeMarkdown
			bot.Send(msg)
			continue
		} else if update.Message.From.ID != config.Admin {
			continue
		}
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.ReplyToMessageID = update.Message.MessageID
			switch update.Message.Command() {
			case "start":
				msg.Text = internal.MessageInitialStart
			default:
				msg.Text = "Unrecognised command!"
			}

			bot.Send(msg)
			continue
		}
	}
}
