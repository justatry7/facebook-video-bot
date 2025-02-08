import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"  // Импортируем godotenv
	"github.com/telegram-bot-api/telegram-bot-api"
)

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/axios/axios" // подключение для скачивания видео
)

func main() {
	// Загружаем переменные окружения из .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Получаем BOT_TOKEN из переменных окружения
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN is not set in the .env file")
	}

	// Инициализация бота
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Bot authorized with username", bot.Self.UserName)

	// Создаем обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	// Обрабатываем команды и сообщения
	for update := range updates {
		if update.Message == nil { // Игнорируем не сообщения
			continue
		}

		// Команда /start
		if update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Send me a Facebook video link, and I'll download it for you!")
			bot.Send(msg)
		}

		// Проверка на ссылку
		if strings.Contains(update.Message.Text, "facebook.com") {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Downloading video, please wait...")
			bot.Send(msg)

			// Логика скачивания видео
			downloadVideo(update.Message.Text)

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Here is your video!")
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "This is not a Facebook video link!")
			bot.Send(msg)
		}
	}
}

// Функция скачивания видео
func downloadVideo(url string) {
	// Код для скачивания видео с Facebook через axios
	_, err := axios.Get("https://api.fdown.net/api/download?url=" + url)
	if err != nil {
		log.Fatal("Error downloading video: ", err)
	}
}
