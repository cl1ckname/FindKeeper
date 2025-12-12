package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"find-keeper/internal/bot"

	"github.com/joho/godotenv"
)

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

	bot, err := bot.NewBot(token, channel)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Bot started for channel: %s", channel)
	go bot.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit
	bot.Stop()
	time.Sleep(time.Millisecond * 10)
}
