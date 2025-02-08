package handler

import (
	"fmt"
	"net/http"
	"os"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func init() {
	// Загружаем переменные окружения из .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Токен бота из переменных окружения
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		http.Error(w, "Bot token not set in environment", http.StatusInternalServerError)
		return
	}

	// Инициализация Telegram API
	bot, err := telegram.NewBotAPI(botToken)
	if err != nil {
		http.Error(w, "Failed to initialize bot", http.StatusInternalServerError)
		return
	}

	// Логирование активности бота
	bot.Debug = true

	// Ответ на запрос
	fmt.Fprintf(w, "Bot is running!")

	// Дальше добавляешь логику бота для скачивания видео...
}

func main() {
	// Ожидаем, что версель будет обращаться по этому пути
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":3000", nil)
}
