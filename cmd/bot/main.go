package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	tgbot "github.com/nurkenti/hadiaParfums/internal/tgbot"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	tgbot.SetMenuCommands(bot)

	bot.Debug = false // Для max ditales

	log.Printf("Authorized on account %s", bot.Self.UserName)

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("Ошибка с токеном")
	}
	handler := tgbot.NewHandler(bot)
	handler.Start()

	/*u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
		for update := range updates {
			if update.Message != nil { // If we got a message
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

				bot.Send()
			}
		}*/

}
