package bot

import (
	"CompareBot/db"
	tg "github.com/mymmrac/telego"
)

func HandleMessage(update tg.Update, userID int64) {
	switch update.Message.Text {
	case "/start":
		newUser := db.User{
			TgID:          userID,
			SelectPic:     0,
			LastMessageID: 0,
			LastPhotoID:   0,
			CountCompare:  0,
		}

		if err := database.InsertUser(newUser); err != nil {
			return
		}
	}
}
