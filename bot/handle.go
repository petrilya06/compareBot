package bot

import (
	"CompareBot/db"
	"fmt"
	tg "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"slices"
)

func HandleMessage(bot *tg.Bot, update tg.Update, userID int64) {
	switch update.Message.Text {
	case "/start":
		newUser := db.User{
			TgID:              userID,
			Confirm:           false,
			SelectPic:         0,
			LastMessageID:     0,
			LastPhotoID:       0,
			CountPhotoCompare: 0,
			CountTextCompare:  0,
		}

		if err := database.InsertUser(newUser); err != nil {
			return
		}

		user, _ = database.GetDataUser(newUser.TgID)

		// подтверждение профиля
		SendMessage(bot, user, *confirmKeyboard, "Подтвердите свой профиль:")

	case "Выбрать аватарку":
		StopChecks()
		SendPhoto(bot, user, inlineKeyboard, fmt.Sprintf("Выберите подходяющую для вас аватарку\n"+
			"За данную аватарку будет выплачиться %s рублей", photoPrices[index]))
	}
}

func HandleCallback(bot *tg.Bot, update tg.Update) {
	switch update.CallbackQuery.Data {
	case "next", "no":
		index += 1
		if index == len(photoPrices) {
			index = 0 // если прокрутили до конца, то возвращаемся к первому элементу
		}

		EditPhotoKeyboard(bot, user, inlineKeyboard, fmt.Sprintf("За данную "+
			"аватарку будет выплачиться %s рублей", photoPrices[index]))

	case "yes":
		user.SelectPic = index
		if err := database.UpdateUser(*user); err != nil {
			return
		}

		DeleteMessages(bot, user, []int{user.LastMessageID, user.LastPhotoID})
		SendMessage(bot, user, *chooseKeyboard, fmt.Sprintf("Спасибо за подтверждение! "+
			"Вам назначена выплата %s рублей", photoPrices[index]))

		StopChannel = make(chan struct{})
		go CheckPhotos(bot, user)

	default:
		if slices.Contains(photoPrices, update.CallbackQuery.Data) {
			EditPhotoKeyboard(bot, user, inlineKeyboardConfirm, "Хотите оставить аватарку?")
		}
	}
}

func HandleContact(bot *tg.Bot, update tg.Update) {
	userInfo, _ = bot.GetChat(&tg.GetChatParams{
		ChatID: tu.ID(user.TgID),
	})
	contact := update.Message.Contact

	if contact.FirstName == userInfo.FirstName && contact.LastName == userInfo.LastName {
		user.Confirm = true
		if err := database.UpdateUser(*user); err != nil {
			return
		}

		SendMessage(bot, user, *chooseKeyboard, "Привет! Ты успешно прошел проверку. "+
			"Помогу вам подзаработать очень простым способом.")
	} else {
		SendMessage(bot, user, *confirmKeyboard, "Это не ты! Прошу использовать настоящий номер")
	}
}
