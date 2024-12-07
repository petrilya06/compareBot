package bot

import (
	"CompareBot/db"
	"fmt"
	tg "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"io"
	"log"
	"net/http"
	"os"
)

func SendMessage(bot *tg.Bot, user *db.User, keyboard tg.ReplyKeyboardMarkup, message string) {
	msg, _ := bot.SendMessage(tu.Message(tu.ID(user.TgID), message).WithReplyMarkup(&keyboard))

	user.LastMessageID = msg.MessageID
	if err := database.UpdateUser(*user); err != nil {
		log.Println(err)
	}
}

func SendPhoto(bot *tg.Bot, user *db.User, keyboard tg.InlineKeyboardMarkup, message string) {
	SendMessage(bot, user, emptyKeyboard, message)
	msg, _ := bot.SendPhoto(tu.Photo(
		tu.ID(user.TgID),
		tu.File(OpenFile(fmt.Sprintf("src/%d.jpg", index))),
	).WithReplyMarkup(&keyboard))

	user.LastPhotoID = msg.MessageID
	if err := database.UpdateUser(*user); err != nil {
		return
	}
}

func EditPhotoKeyboard(bot *tg.Bot, user *db.User, keyboard tg.InlineKeyboardMarkup, message string) {
	_, _ = bot.EditMessageText(&tg.EditMessageTextParams{
		ChatID:    tu.ID(user.TgID),
		MessageID: user.LastMessageID,
		Text:      message,
	})

	_, _ = bot.EditMessageMedia(&tg.EditMessageMediaParams{
		ChatID:    tu.ID(user.TgID),
		MessageID: user.LastPhotoID,
		Media:     tu.MediaPhoto(tu.File(OpenFile(fmt.Sprintf("src/%d.jpg", index)))),
	})

	_, _ = bot.EditMessageReplyMarkup(&tg.EditMessageReplyMarkupParams{
		ChatID:      tu.ID(user.TgID),
		MessageID:   user.LastPhotoID,
		ReplyMarkup: &keyboard,
	})
}

func DeleteMessages(bot *tg.Bot, user *db.User, message []int) {
	_ = bot.DeleteMessages(&tg.DeleteMessagesParams{
		ChatID:     tu.ID(user.TgID),
		MessageIDs: message,
	})
}

func downloadPhoto(bot *tg.Bot, user *db.User) {
	profilePhotos, err := bot.GetUserProfilePhotos(&tg.GetUserProfilePhotosParams{
		UserID: user.TgID,
		Limit:  1,
		Offset: 0,
	})
	if err != nil {
		fmt.Println("error in get user photo:", err)
		return
	}

	if len(profilePhotos.Photos) == 0 {
		SendMessage(bot, user, emptyKeyboard, "У вас нет фотографий!")
		return
	}

	photo := profilePhotos.Photos[0][2]
	file, err := bot.GetFile(&tg.GetFileParams{FileID: photo.FileID})
	if err != nil {
		fmt.Println("error in get file:", err)
		return
	}

	response, err := http.Get("https://api.telegram.org/file/bot" + os.Getenv("TOKEN") + "/" + file.FilePath)
	if err != nil {
		fmt.Println("error in download:", err)
		return
	}
	defer response.Body.Close()

	out, err := os.Create(fmt.Sprintf("src/photos/%d.jpg", user.TgID))
	if err != nil {
		fmt.Println("error in make a file:", err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		fmt.Println("error in write to file:", err)
		return
	}
}
