package bot

import (
	"CompareBot/db"
	"fmt"
	tg "github.com/mymmrac/telego"
	"slices"
)

func HandleMessage(bot *tg.Bot, update tg.Update, userID int64) {
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

		user, _ = database.GetDataUser(newUser.TgID)
		SendMessage(bot, user, *chooseKeyboard,
			"Привет! Помогу вам подзаработать очень простым способом.")

	case "Выбрать аватарку":
		StopChecks()
		SendPhoto(bot, user, inlineKeyboard, fmt.Sprintf("Выберите подходяющую для вас аватарку\n"+
			"За данную аватарку будет выплачиться %s рублей", prices[index]))
	}
}

func HandleCallback(bot *tg.Bot, update tg.Update) {
	switch update.CallbackQuery.Data {
	case "next", "no":
		index += 1
		if index == len(prices) {
			index = 0 // если прокрутили до конца, то возвращаемся к первому элементу
		}

		EditPhotoKeyboard(bot, user, inlineKeyboard, fmt.Sprintf("За данную "+
			"аватарку будет выплачиться %s рублей", prices[index]))

	case "yes":
		user.SelectPic = index
		if err := database.UpdateUser(*user); err != nil {
			return
		}

		DeleteMessages(bot, user, []int{user.LastMessageID, user.LastPhotoID})
		SendMessage(bot, user, *chooseKeyboard, fmt.Sprintf("Спасибо за подтверждение! "+
			"Вам назначена выплата %s рублей", prices[index]))

		StopChannel = make(chan struct{})
		go CheckPhotos(bot, user)

	default:
		if slices.Contains(prices, update.CallbackQuery.Data) {
			EditPhotoKeyboard(bot, user, inlineKeyboardConfirm, "Хотите оставить аватарку?")
		}
	}
}
