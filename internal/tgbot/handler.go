package tgbot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

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

func (h *Handler) SendMessage(chatID int64, t string) {
	msg := tgbotapi.NewMessage(chatID, t)
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
			h.HandlerMsgAdmin(update.Message, chatID)

		case false:
			log.Printf("Chat ID: %d\nMessage: %s", chatID, text)
			h.HandlerMessageUser(update.Message, text, chatID)
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

func (h *Handler) HandlerMessageUser(msg *tgbotapi.Message, text string, chatID int64) {
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
			h.AdmService(chatID, msg.MessageID)
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

func (h *Handler) AdmService(chatID int64, userMsgID int) {
	msg := tgbotapi.NewMessage(chatID, "Как админ вы можете выбрать:")
	msg.ReplyMarkup = CommandForAdmin()
	h.saveMsg(msg, chatID, userMsgID)
	// 1. Товар либо Заказы
	// 2. (Товар) Добавить, Изменить, Удалить, Посмотреть

}

var waiting = make(map[int64]string)
var tempData = make(map[int64]ProductInput)

var (
	StateWaitName        = "w_name"
	StateWaitCategory    = "w_category"
	StateWaitDescription = "w_descript"

	StatusWaitIdProd = "w_id"

	StatusWaitListName = "w_listname"
)

func (h *Handler) AddProdMenu(chatID int64, userMsgID int) {
	msg := tgbotapi.NewMessage(chatID, "Выбирите вариант для Товара:")
	//SaveAndDel[chatID] = DeleteSessionMsg{msgToDel: []int{msg}}
	msg.ReplyMarkup = CommandForProdAdmin()
	h.saveMsg(msg, chatID, userMsgID)
}
func (h *Handler) AskProduct(chatID int64, msgUserID int) {
	// инициализируем чистый структуру для этого чата
	tempData[chatID] = ProductInput{}
	waiting[chatID] = StateWaitName // Первый шаг
	msg := tgbotapi.NewMessage(chatID, "          Шаг 1/3\nНапишите имя товара: ")
	h.saveMsg(msg, chatID, msgUserID)
}

//                Clear functions

type DeleteSessionMsg struct {
	msgToDel []int
}

var saveAndDel = make(map[int64]DeleteSessionMsg)

func (h *Handler) clearSessionMessages(chatID int64, messageIDs []int) {
	for _, msgID := range messageIDs {
		delMsg := tgbotapi.NewDeleteMessage(chatID, msgID)
		// Игнорим ошибку, так как сообщение могло быть уже удалено вручную
		_, _ = h.bot.Send(delMsg)
		delete(saveAndDel, chatID)
	}

}
func (h *Handler) saveMsg(msg tgbotapi.MessageConfig, chatID int64, userMsgID int) {
	// Ловим отправленное ботом сообщ (sentMsg)
	sentMsg, err := h.bot.Send(msg)
	if err != nil {
		fmt.Printf("Ошибка отправки сообщения бота: %v\n", err)
		return
	}
	delData := saveAndDel[chatID]

	// Добавляем ID админа
	if userMsgID > 0 {
		delData.msgToDel = append(delData.msgToDel, userMsgID)
		fmt.Printf("%d", userMsgID)
	}
	// Добавляем ID бота
	delData.msgToDel = append(delData.msgToDel, sentMsg.MessageID)

	// Перезаписываем мапу
	saveAndDel[chatID] = delData
}

//                  Main Functions

type ProductInput struct {
	name         string
	category     string
	descrip      string
	id           int32
	messageToDel []int // Массив для удаление сообщ одной сессий
}

func (h *Handler) HandlerMsgAdmin(msg *tgbotapi.Message, chatID int64) {
	state, isWaiting := waiting[chatID]
	if isWaiting {
		h.handlerAdminSteps(state, msg, chatID, msg.MessageID)
		return
	}

	switch msg.Text {
	case "/start":
		h.clearSessionMessages(chatID, saveAndDel[chatID].msgToDel)
		h.AdmService(chatID, msg.MessageID)
	case "Товар":
		h.AddProdMenu(chatID, msg.MessageID)
	case "Добавить":
		h.AskProduct(chatID, msg.MessageID)
	case "Удалить":
		h.DeleteProduct(chatID, msg.MessageID)
	case "Список":
		h.ListProdName(chatID, msg.MessageID)

	}
}
func (h *Handler) ListProdName(chatID int64, userMsgID int) {
	waiting[chatID] = StatusWaitListName

	msg := tgbotapi.NewMessage(chatID, "Напиши название товара для поиска: ")
	h.saveMsg(msg, chatID, userMsgID)
}

func (h *Handler) DeleteProduct(chatID int64, userMsgID int) {
	tempData[chatID] = ProductInput{}
	waiting[chatID] = StatusWaitIdProd
	msg := tgbotapi.NewMessage(chatID, "Напишите ID продукта для удаление: ")
	h.saveMsg(msg, chatID, userMsgID)
}

func (h *Handler) handlerAdminSteps(state string, msg *tgbotapi.Message, chatID int64, userMsgID int) {
	product := h.productService
	data := tempData[chatID] // Достаем то что уже успели записать
	delData := saveAndDel[chatID]

	text := msg.Text

	switch state {
	//                  AddProd

	case StateWaitName:

		delData.msgToDel = append(delData.msgToDel, msg.MessageID)
		// Записоваем Имя
		data.name = text
		tempData[chatID] = data

		// Переводим на след шаг
		waiting[chatID] = StateWaitCategory
		msg := tgbotapi.NewMessage(chatID, "          Шаг 2/3\nВыбери категорию товара: ")
		msg.ReplyMarkup = CommandForProduct()
		h.saveMsg(msg, chatID, userMsgID)

	case StateWaitCategory:
		delData.msgToDel = append(delData.msgToDel, msg.MessageID)
		// Записываем категорию
		data.category = text
		tempData[chatID] = data

		waiting[chatID] = StateWaitDescription
		msg := tgbotapi.NewMessage(chatID, "          Шаг 3/3\nНапишите описание товара:")
		h.saveMsg(msg, chatID, userMsgID)

	case StateWaitDescription:
		delData.msgToDel = append(delData.msgToDel, msg.MessageID)
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
			h.saveMsg(msg, chatID, userMsgID)
			fmt.Println("СУУУКА")
		} else {
			reply := fmt.Sprintf("✅ Товар успешно добавлен!\n\n🏷 Название: %s\n📁 Категория: %s\n📝 Описание: %s", data.name, data.category, data.descrip)
			msg := tgbotapi.NewMessage(chatID, reply)
			h.saveMsg(msg, chatID, userMsgID)
			fmt.Println("БУУКА")
		}
		h.clearSessionMessages(chatID, delData.msgToDel)

		// Очищаем состояние и временные данные, чтобы админ мог снова пользоваться кнопками
		delete(waiting, chatID)
		delete(tempData, chatID)

		// ------------------------------------------------

		// Delete Prod
	case StatusWaitIdProd:
		delData.msgToDel = append(delData.msgToDel, msg.MessageID)

		id, err := strconv.ParseInt(text, 10, 32)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "❌ Неверный формат ID. Пожалуйста, введите число:")
			h.saveMsg(msg, chatID, userMsgID)
		}
		data.id = int32(id)

		ctx := context.Background()
		err = product.DeleteProduct(ctx, int32(id))
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "❌ Ошибка при удалений")
			h.saveMsg(msg, chatID, userMsgID)
		} else {
			msg := tgbotapi.NewMessage(chatID, "✅ Товар успешно удален!")
			h.saveMsg(msg, chatID, userMsgID)
		}
		h.clearSessionMessages(chatID, delData.msgToDel)
		h.clearSessionMessages(chatID, saveAndDel[chatID].msgToDel)
		delete(waiting, chatID)
		delete(tempData, chatID)

	case StatusWaitListName:
		delData.msgToDel = append(delData.msgToDel, msg.MessageID)
		data.name = text
		Prods, err := product.ListProdByName(context.Background(), data.name)
		if err != nil {
			// Ловим сообщ об ошибке, чтобы оно стрело
			errMsg, _ := h.bot.Send(tgbotapi.NewMessage(chatID, "❌ Ошибка при поиске"))
			data.messageToDel = append(data.messageToDel, errMsg.MessageID)
			delete(waiting, chatID)
			return
		}
		// Если ничего не нашли
		if len(Prods) == 0 {
			notFoundMsg, _ := h.bot.Send(tgbotapi.NewMessage(chatID, "📭 Товары с таким названием не найдены."))
			data.messageToDel = append(data.messageToDel, notFoundMsg.MessageID)
			time.Sleep(3 * time.Second) // Чтобы админ прочитал даем 3 сек

			// Автоматический удаляем все промежуточные вопросы бота и сообщ user
			h.clearSessionMessages(chatID, delData.msgToDel)
			h.clearSessionMessages(chatID, saveAndDel[chatID].msgToDel)

			delete(waiting, chatID)
			delete(tempData, chatID)
		}
		// strings.Builder для красивой сборке текста списка
		var builder strings.Builder
		builder.WriteString("📋 **Результаты поиска товаров:**\n\n")

		for i, prod := range Prods {
			itemText := fmt.Sprintf("%d. 📦 **%s** (ID: `%d`)\n📁 Категория: %s\n📝 Описание: %s\n\n",
				i+1, prod.Name, prod.ID, prod.Category, prod.Description)

			builder.WriteString(itemText)
		}

		// Отправляем одно сообщение всех списков
		msg := tgbotapi.NewMessage(chatID, builder.String())
		msg.ParseMode = "Markdown"
		h.bot.Send(msg)

		// Автоматический удаляем все промежуточные вопросы бота и сообщ user
		//h.clearSessionMessages(chatID, data.messageToDel)
		h.clearSessionMessages(chatID, delData.msgToDel)
		h.clearSessionMessages(chatID, saveAndDel[chatID].msgToDel)

		delete(waiting, chatID)
		delete(tempData, chatID)

	}

}
