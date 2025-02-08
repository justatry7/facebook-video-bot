package handler

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"github.com/joho/godotenv"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"encoding/json"
)

// Структура для работы с ответом от API
type FBResponse struct {
	DownloadURL string `json:"download_url"`
}

// Экспортируемая функция для Vercel
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	// Загружаем переменные окружения из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Получаем токен бота из переменных окружения
	apiToken := os.Getenv("BOT_TOKEN")
	if apiToken == "" {
		log.Fatal("BOT_TOKEN is not set in the .env file")
	}

	// Создаем бота
	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	// Обработка запроса
	if r.Method == "POST" {
		var update tgbotapi.Update
		// Читаем тело запроса
		err := json.NewDecoder(r.Body).Decode(&update)
		if err != nil {
			http.Error(w, "Failed to parse request", http.StatusBadRequest)
			return
		}

		if update.Message != nil {
			if update.Message.Text == "/start" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Send me a Facebook video link, and I'll download it for you!")
				bot.Send(msg)
			} else if isFacebookLink(update.Message.Text) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Downloading video, please wait...")
				bot.Send(msg)

				// Получаем ссылку на видео с помощью API
				downloadURL, err := getFacebookVideoDownloadURL(update.Message.Text)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "An error occurred while downloading the video.")
					bot.Send(msg)
					return
				}

				// Отправляем видео пользователю
				videoMsg := tgbotapi.NewVideo(update.Message.Chat.ID, downloadURL)
				bot.Send(videoMsg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "This is not a Facebook video link!")
				bot.Send(msg)
			}
		}
	}

	// Ответ на запрос
	fmt.Fprintf(w, "Bot is running!")
}

// Проверка, является ли ссылка ссылкой на видео с Facebook
func isFacebookLink(url string) bool {
	return contains(url, "facebook.com")
}

// Проверка наличия подстроки в строке
func contains(str, substr string) bool {
	return len(str) > 0 && len(substr) > 0 && (str == substr || contains(str[1:], substr))
}

// Получение URL для скачивания видео из Facebook
func getFacebookVideoDownloadURL(url string) (string, error) {
	// API для скачивания видео с Facebook
	apiURL := fmt.Sprintf("https://api.fdown.net/api/download?url=%s", url)
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Декодируем JSON ответ
	var response FBResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	// Возвращаем URL для скачивания
	if response.DownloadURL != "" {
		return response.DownloadURL, nil
	}
	return "", fmt.Errorf("unable to retrieve download URL")
}
