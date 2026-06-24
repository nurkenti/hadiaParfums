package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Handler struct {
	bot *tgbotapi.BotAPI
}

func NewHandler(token string) (*Handler, error) {
	tgbot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Handler{bot: tgbot}, nil
}

func (h *Handler) HandlerStart(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Салам! Это Hadia Parfums")
	h.bot.Send(msg)
}

func (h *Handler) Start() {
	u := tgbotapi.NewUpdate(0)
	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.Text == "/start" {
			h.HandlerStart(update.Message.Chat.ID)
		}
	}

}
