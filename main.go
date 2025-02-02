package main
//я люблю Свету
import (
	"encoding/json"
	"fmt"
	"log"
	//"math/rand"
	"net/http"
	//"os"
	//"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Структура для хранения шуток
type Joke struct {
	Setup   string `json:"setup"`
	Punchline string `json:"punchline"`
}
//go get github.com/go-telegram-bot-api/telegram-bot-api/v5
// Функция для получения случайной шутки из API
func getJoke() (string, error) {
	resp, err := http.Get("https://official-joke-api.appspot.com/jokes/random")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var joke Joke
	if err := json.NewDecoder(resp.Body).Decode(&joke); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n%s", joke.Setup, joke.Punchline), nil
}

func main() {
	// Читаем токен бота из переменных окружения
	token := "token"
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil { // Игнорируем некорректные обновления
			continue
		}

		switch update.Message.Text {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Отправь /joke, чтобы получить шутку!")
			bot.Send(msg)
		case "/joke":
			joke, err := getJoke()
			if err != nil {
				joke = "Не удалось получить шутку. Попробуйте позже."
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, joke)
			bot.Send(msg)
		case "/help":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Доступные команды:\n/start - Начать\n/joke - Получить шутку\n/help - Справка")
			bot.Send(msg)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда. Используйте /help для списка команд.")
			bot.Send(msg)
		}
	}
}
