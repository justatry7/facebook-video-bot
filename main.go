package main

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

// Структура для работы с ответом от API
type FBResponse struct {
	DownloadURL string `json:"download_url"`
}

func main() {
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

	// Логика обработки сообщений
	updates, err := bot.GetUpdatesChan(tgbotapi.NewUpdate(0))
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
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
					continue
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

	// Чтение ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Декодируем JSON ответ
	var response FBResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	// Возвращаем URL для скачивания
	if response.DownloadURL != "" {
		return response.DownloadURL, nil
	}
	return "", fmt.Errorf("unable to retrieve download URL")
}
