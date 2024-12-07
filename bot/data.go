package bot

import (
	"CompareBot/db"
	tg "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

var StopChannel = make(chan struct{})
var database *db.Database
var user *db.User
var index = 0
var prices = []string{"1000", "2000"}

var chooseKeyboard = tu.Keyboard(
	tu.KeyboardRow(tu.KeyboardButton("Выбрать аватарку")),
).WithResizeKeyboard().WithOneTimeKeyboard()

var inlineKeyboardConfirm = tg.InlineKeyboardMarkup{
	InlineKeyboard: [][]tg.InlineKeyboardButton{
		{
			{
				Text:         "Да",
				CallbackData: "yes",
			},
			{
				Text:         "Нет",
				CallbackData: "no",
			},
		},
	},
}

var inlineKeyboard = tg.InlineKeyboardMarkup{
	InlineKeyboard: [][]tg.InlineKeyboardButton{
		{
			{
				Text:         "Выбрать",
				CallbackData: prices[index],
			},
			{
				Text:         "Следующая",
				CallbackData: "next",
			},
		},
	},
}

var emptyKeyboard = tg.ReplyKeyboardMarkup{
	Keyboard: [][]tg.KeyboardButton{},
}
