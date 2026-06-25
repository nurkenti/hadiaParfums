package tgbot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Handler struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func (h *Handler) SendMessage(t string) {
	msg := tgbotapi.NewMessage(h.chatID, t)
	h.bot.Send(msg)
}
func NewHandler(bot *tgbotapi.BotAPI) *Handler {
	return &Handler{bot: bot}
}

func (h *Handler) HandlerStart() {
	msg := tgbotapi.NewMessage(h.chatID, "Салам! Это Hadia Parfums")
	msg.ReplyMarkup = mainMenu()
	h.bot.Send(msg)

}
func (h *Handler) HandlerInfo() {
	msg := tgbotapi.NewMessage(h.chatID, "Информация жайлы")
	msg.ReplyMarkup = mainMenu()
	h.bot.Send(msg)

}

func (h *Handler) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		h.chatID = update.Message.Chat.ID
		h.HandlerMessage(update.Message.Text)

	}

}

func (h *Handler) HandlerMessage(text string) {
	switch text {
	case "/start":
		h.HandlerStart()
	case "/info":
		h.HandlerInfo()
	default:
		h.SendMessage("Щщс дурстап жазсай")
	}
}
