package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	channel int64
}

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
	}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.ForwardFromMessageID != 0 {
			b.handleForwardedMessage(update.Message)
		}
	}
}

func (b *Bot) handleForwardedMessage(message *tgbotapi.Message) {
	if message.Photo != nil && len(message.Photo) > 0 {
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
	
	file, err := b.api.GetFile(tgbotapi.FileConfig{FileID: photo.FileID})
	if err != nil {
		log.Printf("Error getting photo file: %v", err)
		return
	}

	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", b.api.Token, file.FilePath)
	
	tempFile, err := b.downloadFile(fileURL)
	if err != nil {
		log.Printf("Error downloading photo: %v", err)
		return
	}
	defer os.Remove(tempFile)

	photoMsg := tgbotapi.NewPhoto(b.channel, tgbotapi.FilePath(tempFile))

	_, err = b.api.Send(photoMsg)
	if err != nil {
		log.Printf("Error sending photo: %v", err)
	}
}

func (b *Bot) handleVideo(message *tgbotapi.Message) {
	video := message.Video
	
	file, err := b.api.GetFile(tgbotapi.FileConfig{FileID: video.FileID})
	if err != nil {
		log.Printf("Error getting video file: %v", err)
		return
	}

	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", b.api.Token, file.FilePath)
	
	tempFile, err := b.downloadFile(fileURL)
	if err != nil {
		log.Printf("Error downloading video: %v", err)
		return
	}
	defer os.Remove(tempFile)

	videoMsg := tgbotapi.NewVideo(b.channel, tgbotapi.FilePath(tempFile))

	_, err = b.api.Send(videoMsg)
	if err != nil {
		log.Printf("Error sending video: %v", err)
	}
}

func (b *Bot) handleAnimation(message *tgbotapi.Message) {
	animation := message.Animation
	
	file, err := b.api.GetFile(tgbotapi.FileConfig{FileID: animation.FileID})
	if err != nil {
		log.Printf("Error getting animation file: %v", err)
		return
	}

	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", b.api.Token, file.FilePath)
	
	tempFile, err := b.downloadFile(fileURL)
	if err != nil {
		log.Printf("Error downloading animation: %v", err)
		return
	}
	defer os.Remove(tempFile)

	animationMsg := tgbotapi.NewAnimation(b.channel, tgbotapi.FilePath(tempFile))

	_, err = b.api.Send(animationMsg)
	if err != nil {
		log.Printf("Error sending animation: %v", err)
	}
}

func (b *Bot) handleDocument(message *tgbotapi.Message) {
	doc := message.Document
	
	if !b.isVideoDocument(doc.MimeType) {
		return
	}

	file, err := b.api.GetFile(tgbotapi.FileConfig{FileID: doc.FileID})
	if err != nil {
		log.Printf("Error getting document file: %v", err)
		return
	}

	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", b.api.Token, file.FilePath)
	
	tempFile, err := b.downloadFile(fileURL)
	if err != nil {
		log.Printf("Error downloading document: %v", err)
		return
	}
	defer os.Remove(tempFile)

	videoMsg := tgbotapi.NewVideo(b.channel, tgbotapi.FilePath(tempFile))

	_, err = b.api.Send(videoMsg)
	if err != nil {
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

func (b *Bot) downloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tempFile, err := os.CreateTemp("", "telegram_media_*")
	if err != nil {
		return "", err
	}

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return "", err
	}

	tempFile.Close()
	return tempFile.Name(), nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	channel := os.Getenv("TELEGRAM_CHANNEL")
	if channel == "" {
		log.Fatal("TELEGRAM_CHANNEL environment variable is required")
	}

	bot, err := NewBot(token, channel)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Bot started for channel: %s", channel)
	bot.Start()
}
