package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	service "github.com/nurkenti/hadiaParfums/internal/core/service/admin"
	db "github.com/nurkenti/hadiaParfums/internal/repository/sqlc"
	tgbot "github.com/nurkenti/hadiaParfums/internal/tgbot"
)

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" {
		dbHost = "localhost"
	}

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		log.Fatal("DB_SOURCE не заполнен в .env")
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("Ошибка с токеном")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	tgbot.SetMenuCommands(bot)

	bot.Debug = false // Для max ditales

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Подключение к базе данных
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal("db isn't connect: ", err)
	}
	defer conn.Close(context.Background())

	// Подключаем все к клиенту и админу
	store := db.New(conn)
	adminService := service.NewAdminService(store)
	productService := service.NewProductService(store)

	handler := tgbot.NewHandler(bot, adminService, productService)
	handler.Start()

}
