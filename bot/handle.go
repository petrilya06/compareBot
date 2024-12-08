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
		user = &db.User{TgID: userID}
		if err := database.InsertUser(*user); err != nil {
			return
		}

		// подтверждение профиля
		SendMessage(bot, user, *confirmKeyboard, "Подтвердите свой профиль:")

	case "Выбрать аватарку":
		StopChecks(StopChannelPhoto)
		SendPhoto(bot, user, InlineKeyboardPhoto, fmt.Sprintf("Выберите подходяющую для вас аватарку\n"+
			"За данную аватарку будет выплачиться %s рублей", photoPrices[user.SelectPic]))
	}

}

func HandleCallback(bot *tg.Bot, update tg.Update) {
	switch update.CallbackQuery.Data {
	case "next", "no":
		user.SelectPic += 1
		if user.SelectPic == len(photoPrices) {
			user.SelectPic = 0 // если прокрутили до конца, то возвращаемся к первому элементу
		}
		if err := database.UpdateUser(*user); err != nil {
			return
		}

		EditPhotoKeyboard(bot, user, InlineKeyboardPhoto, fmt.Sprintf("За данную "+
			"аватарку будет выплачиться %s рублей", photoPrices[user.SelectPic]))

	case "yes":
		if err := database.UpdateUser(*user); err != nil {
			return
		}

		DeleteMessages(bot, user, []int{user.LastMessageID, user.LastPhotoID})
		SendMessage(bot, user, *chooseKeyboard, fmt.Sprintf("Спасибо за подтверждение! "+
			"Вам назначена выплата %s рублей", photoPrices[user.SelectPic]))

		StopChannelPhoto = make(chan struct{})
		go CheckPhotos(bot, user)

	default:
		if slices.Contains(photoPrices, update.CallbackQuery.Data) {
			EditPhotoKeyboard(bot, user, inlineKeyboardConfirm, "Хотите оставить аватарку?")
		}
	}
}

func HandleContact(bot *tg.Bot, update tg.Update) {
	userInfo, _ = bot.GetChat(&tg.GetChatParams{ChatID: tu.ID(user.TgID)})
	contact = update.Message.Contact

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
