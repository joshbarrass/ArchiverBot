package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joshbarrass/ArchiverBot/internal"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Configuration struct {
	BotToken  string `envconfig:"BOT_TOKEN" required:"true"`
	Admin     int    `envconfig:"ADMIN_ID" default:"0"`
	DebugMode bool   `envconfig:"DEBUG_MODE" default:"false"`
	OutDir    string `envconfig:"OUT_DIR" default:"./downloads"`
}

const DefaultDirMode os.FileMode = 0755

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

	config.OutDir, err = filepath.Abs(config.OutDir)
	if err != nil {
		logrus.Fatalf("Failed to make output dir absolute: %s", err)
	}
	if info, err := os.Stat(config.OutDir); os.IsNotExist(err) {
		err = os.Mkdir(config.OutDir, DefaultDirMode)
		if err != nil {
			logrus.Fatalf("Failed to make new output dir: %s", err)
		}
		logrus.Infof("Made new output dir '%s'", config.OutDir)
	} else if !info.IsDir() {
		logrus.Fatalf("OutDir '%s' already exists and is not a directory", config.OutDir)
	}

	// check perms are available in the output dir
	if !internal.TestReadPermission(config.OutDir) {
		logrus.Fatalf("No read permissions for output dir '%s'", config.OutDir)
	}
	if !internal.TestWritePermission(config.OutDir) {
		logrus.Fatalf("No write permissions for output dir '%s'", config.OutDir)
	}

	logrus.Infof("Output dir set to '%s'", config.OutDir)

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
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, internal.MessageIsNotAdmin)
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ParseMode = tgbotapi.ModeMarkdown
			bot.Send(msg)
			continue
		}
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.ReplyToMessageID = update.Message.MessageID
			switch update.Message.Command() {
			case "start":
				msg.Text = internal.MessageInitialStart
			case "download":
				arg := update.Message.CommandArguments()
				if err := internal.DownloadAutoOutput(arg, config.OutDir, bot, update); err != nil {
					msg.Text = fmt.Sprintf("Failed to download. Err: %s", err)
				} else {
					msg.Text = "Download successful!"
				}
			case "listdir":
				msg.Text = fmt.Sprintf("Content of '%s':", config.OutDir)
				dirs, err := ioutil.ReadDir(config.OutDir)
				if err != nil {
					msg.Text += fmt.Sprintf("Failed to list. Err: %s", err)
				} else {
					for _, f := range dirs {
						msg.Text += "\n" + f.Name()
					}
				}
			default:
				msg.Text = "Unrecognised command!"
			}

			bot.Send(msg)
			continue
		}
	}
}
