package tgbot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func mainMenu() tgbotapi.ReplyKeyboardMarkup {
	key := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Каталог"),
			tgbotapi.NewKeyboardButton("Инфо"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Заказать"),
			tgbotapi.NewKeyboardButton("Соц сети"),
		),
	)
	key.ResizeKeyboard = true
	key.Selective = true

	return key
}

func SetMenuCommands(bot *tgbotapi.BotAPI) error {
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Начать"},
		{Command: "catalog", Description: "Каталог"},
		{Command: "order", Description: "Заказать"},
		{Command: "info", Description: "Информация"},
	}
	config := tgbotapi.NewSetMyCommands(commands...)
	_, err := bot.Request(config)
	return err

}
