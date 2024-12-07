package bot

import (
	"CompareBot/db"
	"github.com/joho/godotenv"
	tg "github.com/mymmrac/telego"
	"log"
	"os"
)

func Bot() {
	var err error
	database, err = db.NewDatabase()
	if err != nil {
		log.Fatalf("error in database connection: %v", err)
	}

	if err = database.CreateTable(); err != nil {
		log.Fatalf("error in create table: %v", err)
	}

	// Загрузка переменных окружения из .env файла
	if err = godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, _ := tg.NewBot(os.Getenv("TOKEN"), tg.WithDefaultDebugLogger())
	updates, _ := bot.UpdatesViaLongPolling(nil)

	defer bot.StopLongPolling()
	defer database.CloseDatabase()

	for update := range updates {
		if update.Message != nil {
			userID := update.Message.From.ID
			HandleMessage(bot, update, userID)
		}

		if update.CallbackQuery != nil {
			HandleCallback(bot, update)
		}
	}
}
