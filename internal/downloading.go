package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joshbarrass/ArchiverBot/pkg/uarchiver"
	"github.com/sirupsen/logrus"
)

func DownloadAuto(url, outdir string) error {
	// backup old working dir
	oldDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working dir: %w", err)
	}

	// move to outdir
	err = os.Chdir(outdir)
	if err != nil {
		return fmt.Errorf("failed to move to outdir: %w", err)
	}

	// call UArchiver
	err = uarchiver.DownloadAuto(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}

	// return to old dir
	err = os.Chdir(oldDir)
	if err != nil {
		return fmt.Errorf("failed to return to old working dir: %w", err)
	}

	return nil
}

func DownloadAutoOutput(url, outdir string, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	// backup old working dir
	oldDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working dir: %w", err)
	}

	// move to outdir
	err = os.Chdir(outdir)
	if err != nil {
		return fmt.Errorf("failed to move to outdir: %w", err)
	}

	// call UArchiver
	wait, stdout, stderr, err := uarchiver.DownloadAutoOutput(url)
	if err != nil {
		return fmt.Errorf("failed to start download: %w", err)
	}
	defer stdout.Close()
	defer stderr.Close()

	// start error outputter
	ctx, cancel := context.WithCancel(context.Background())
	go BotErrorLogger(ctx, stdout, stderr, bot, update)
	err = wait()
	cancel()
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}

	// return to old dir
	err = os.Chdir(oldDir)
	if err != nil {
		return fmt.Errorf("failed to return to old working dir: %w", err)
	}

	return nil
}

// BotLogger logs stdout and stderr at end of context
func BotErrorLogger(ctx context.Context, stdout, stderr uarchiver.StdPipe, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	allOut := []byte{}
	go func() { allOut, _ = ioutil.ReadAll(stdout) }()
	allError := []byte{}
	go func() { allError, _ = ioutil.ReadAll(stderr) }()
	<-ctx.Done()

	bufOut := bytes.NewReader(allOut)
	bufErr := bytes.NewReader(allError)

	lastMsgID := update.Message.MessageID

	for _, reader := range [2]io.Reader{bufOut, bufErr} {
		buf := make([]byte, 4096)
		n, err := reader.Read(buf)
		if err == nil {
			for n > 0 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(buf[:n]))
				msg.ReplyToMessageID = lastMsgID

				var sentMsg tgbotapi.Message
				sentMsg, err = bot.Send(msg)
				if err != nil {
					logrus.Errorf("Failed to send: %s", err)
					break
				}
				lastMsgID = sentMsg.MessageID

				n, err = reader.Read(buf)
				if err != nil {
					logrus.Errorf("Failed to read buffer: %s", err)
					break
				}
			}
		} else {
			logrus.Errorf("Failed to read buffer: %s", err)
		}
	}
}
