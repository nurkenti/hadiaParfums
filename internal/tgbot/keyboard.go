package tgbot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func menuLang() tgbotapi.ReplyKeyboardMarkup {
	key := tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Казахский"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Руский"),
		),
	)
	/*key.ResizeKeyboard = true
	key.Selective = true*/

	return key
}
func mainMenu() tgbotapi.ReplyKeyboardMarkup {
	key := tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Каталог"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Заказать"),
		),
	)
	/*key.ResizeKeyboard = true
	key.Selective = true*/

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

func CommandForAdmin() tgbotapi.ReplyKeyboardMarkup {
	commands := tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Товар"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Заказы"),
		),
	)
	return commands
}

func CommandForProdAdmin() tgbotapi.ReplyKeyboardMarkup {
	commands := tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавить"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Изменить"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Список"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Удалить"),
		),
	)
	return commands
}
