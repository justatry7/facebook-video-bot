package handler

import (
  "fmt"
  "github.com/go-telegram-bot-api/telegram-bot-api"
  "os"
  "net/http"
  "github.com/joho/godotenv"
)

func init() {
  err := godotenv.Load()
  if err != nil {
    fmt.Println("Error loading .env file")
  }
}

func Handler(w http.ResponseWriter, r *http.Request) {
  botToken := os.Getenv("BOT_TOKEN")
  if botToken == "" {
    fmt.Fprintf(w, "BOT_TOKEN is not set in environment variables!")
    return
  }

  bot, err := tgbotapi.NewBotAPI(botToken)
  if err != nil {
    fmt.Fprintf(w, "Failed to initialize bot: %v", err)
    return
  }

  fmt.Fprintf(w, "Bot successfully initialized!")
}
