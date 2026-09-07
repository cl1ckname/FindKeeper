package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	channel int64
	close   closeChan
}

type closeChan = chan struct{}

func NewBot(token, channel string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	channelID, err := strconv.ParseInt(channel, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid channel ID: %v", err)
	}

	return &Bot{
		api:     bot,
		channel: channelID,
		close:   make(closeChan),
	}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)
	for {
		select {
		case update := <-updates:
			if update.Message != nil && update.Message.ForwardFromMessageID != 0 {
				b.handleForwardedMessage(update.Message)
			}
		case <-b.close:
			return
		}
	}
}

func (b *Bot) Stop() {
	b.close <- struct{}{}
}

func (b *Bot) handleForwardedMessage(message *tgbotapi.Message) {
	if len(message.Photo) > 0 {
		b.handlePhoto(message)
	} else if message.Video != nil {
		b.handleVideo(message)
	} else if message.Animation != nil {
		b.handleAnimation(message)
	} else if message.Document != nil {
		b.handleDocument(message)
	}
}

func (b *Bot) handlePhoto(message *tgbotapi.Message) {
	photo := message.Photo[len(message.Photo)-1]
	msg := tgbotapi.NewPhoto(b.channel, tgbotapi.FileID(photo.FileID))
	msg.Caption = message.Caption
	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Error sending photo: %v", err)
	}
}

func (b *Bot) handleVideo(message *tgbotapi.Message) {
	msg := tgbotapi.NewVideo(b.channel, tgbotapi.FileID(message.Video.FileID))
	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Error sending video: %v", err)
	}
}

func (b *Bot) handleAnimation(message *tgbotapi.Message) {
	msg := tgbotapi.NewAnimation(b.channel, tgbotapi.FileID(message.Animation.FileID))
	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Error sending animation: %v", err)
	}
}

func (b *Bot) handleDocument(message *tgbotapi.Message) {
	doc := message.Document
	if !b.isVideoDocument(doc.MimeType) {
		return
	}
	msg := tgbotapi.NewVideo(b.channel, tgbotapi.FileID(doc.FileID))
	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Error sending document as video: %v", err)
	}
}

func (b *Bot) isVideoDocument(mimeType string) bool {
	videoTypes := []string{"video/", "application/x-mpegURL", "application/vnd.apple.mpegurl"}
	for _, vType := range videoTypes {
		if strings.HasPrefix(mimeType, vType) {
			return true
		}
	}
	return false
}
