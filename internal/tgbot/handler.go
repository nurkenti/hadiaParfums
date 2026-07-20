package tgbot

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgtype"
	service "github.com/nurkenti/hadiaParfums/internal/core/service/admin"
)

type Handler struct {
	bot            *tgbotapi.BotAPI
	chatID         int64
	adminService   *service.AdminService
	productService *service.ProductService
}

func (h *Handler) SendMessage(t string) {
	msg := tgbotapi.NewMessage(h.chatID, t)
	h.bot.Send(msg)
}

func NewHandler(bot *tgbotapi.BotAPI, adminService *service.AdminService, productService *service.ProductService) *Handler {
	return &Handler{bot: bot,
		adminService:   adminService,
		productService: productService}
}

func (h *Handler) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		chatID := update.Message.Chat.ID
		text := update.Message.Text
		switch h.Admin(chatID) {
		case true:
			log.Printf("Admin ID: %d\nMessage: %s", chatID, text)
			h.HandlerMsgAdmin(text, chatID)

		case false:
			log.Printf("Chat ID: %d\nMessage: %s", chatID, text)
			h.HandlerMessageUser(text, chatID)
		}
	}

}

func (h *Handler) HandlerStart(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Салам! Это Hadia Parfums\nПрежде чем начать Выберите Язык:")
	msg.ReplyMarkup = menuLang()
	h.bot.Send(msg)

}
func (h *Handler) HandlerInfo(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Информация жайлы")
	h.bot.Send(msg)

}
func (h *Handler) KzInfo(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Кош келдыныз, Не калайсыз Шеф?")
	msg.ReplyMarkup = mainMenu()
	h.bot.Send(msg)

}
func (h *Handler) RuInfo(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Добро пожаловать, чего желаете ?")
	msg.ReplyMarkup = mainMenu()
	h.bot.Send(msg)

}

func (h *Handler) Catalog(chatID int64) {
}

func (h *Handler) HandlerMessageUser(text string, chatID int64) {
	switch text {
	case "/start":
		h.HandlerStart(chatID)
	case "/info":
		h.HandlerInfo(chatID)
	case "Казахский":
		h.KzInfo(chatID)
	case "Руский":
		h.RuInfo(chatID)
	case "Каталог":
		h.Catalog(chatID)
	case "012004adm":
		a := h.Admin(chatID)
		if a == true {
			h.AdmService(chatID)
		}
	}
}

// ADMIN
type AnswerAdminForProd struct {
	name string
}

func (h *Handler) Admin(chatID int64) bool {
	adm := h.adminService
	checkadm, _ := adm.IsAdmin(context.Background(), chatID)
	if checkadm == true {
		//h.SendMessage("Ты админ! Добро пожаловать")
		return true
	}
	err := adm.CreateAdmin(context.Background(), 969867088, "Nurken")
	if err != nil {
		return false
	}
	msg := tgbotapi.NewMessage(h.chatID, "Ты стал админом! Добро пожаловать")
	h.bot.Send(msg)
	return true
}

func (h *Handler) AdmService(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Как админ вы можете выбрать:")
	msg.ReplyMarkup = CommandForAdmin()
	h.bot.Send(msg)
	// 1. Товар либо Заказы
	// 2. (Товар) Добавить, Изменить, Удалить, Посмотреть

	//	prod := h.productService
}

var waiting = make(map[int64]string)
var tempData = make(map[int64]ProductInput)

var (
	StateWaitName        = "w_name"
	StateWaitCategory    = "w_category"
	StateWaitDescription = "w_descript"
)

func (h *Handler) AddProdMenu(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Выбирите вариант для Товара:")
	msg.ReplyMarkup = CommandForProdAdmin()
	h.bot.Send(msg)
}
func (h *Handler) AskProduct(chatID int64) {
	// инициализируем чистый структуру для этого чата
	tempData[chatID] = ProductInput{}
	waiting[chatID] = StateWaitName // Первый шаг
	msg := tgbotapi.NewMessage(chatID, "          Шаг 1/3\nНапишите имя товара: ")
	h.bot.Send(msg)
}

type ProductInput struct {
	name     string
	category string
	descrip  string
}

func (h *Handler) HandlerMsgAdmin(text string, chatID int64) {
	state, isWaiting := waiting[chatID]
	if isWaiting {
		h.handlerAdminSteps(state, text, chatID)
		return
	}
	switch text {
	case "/start":
		h.AdmService(chatID)

	case "Товар":
		h.AddProdMenu(chatID)
	case "Добавить":
		h.AskProduct(chatID)

	}
}
func (h *Handler) handlerAdminSteps(state, text string, chatID int64) {
	product := h.productService
	data := tempData[chatID] // Достаем то что уже успели записать

	switch state {
	case StateWaitName:
		// Записоваем Имя
		data.name = text
		tempData[chatID] = data

		// Переводим на след шаг
		waiting[chatID] = StateWaitCategory
		msg := tgbotapi.NewMessage(chatID, "          Шаг 2/3\nВыбери категорию товара: ")
		msg.ReplyMarkup = CommandForProduct()
		h.bot.Send(msg)

	case StateWaitCategory:
		// Записываем категорию
		data.category = text
		tempData[chatID] = data

		waiting[chatID] = StateWaitDescription
		msg := tgbotapi.NewMessage(chatID, "          Шаг 3/3\nНапишите описание товара:")
		h.bot.Send(msg)

	case StateWaitDescription:
		data.descrip = text

		// Все данные сохранены. Добавим в db
		ctx := context.Background()
		err := product.AddProduct(
			data.name,
			data.category,
			pgtype.Text{String: data.descrip, Valid: true},
			ctx,
		)

		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "❌ Ошибка при сохранении в БД: %v")
			h.bot.Send(msg)
			fmt.Println("СУУУКА")
		} else {
			reply := fmt.Sprintf("✅ Товар успешно добавлен!\n\n🏷 Название: %s\n📁 Категория: %s\n📝 Описание: %s", data.name, data.category, data.descrip)
			msg := tgbotapi.NewMessage(chatID, reply)
			h.bot.Send(msg)
			fmt.Println("БУУКА")
		}
		// Очищаем состояние и временные данные, чтобы админ мог снова пользоваться кнопками
		delete(waiting, chatID)
		delete(tempData, chatID)

	}
}
