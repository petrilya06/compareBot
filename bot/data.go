package bot

import (
	"CompareBot/db"
	tg "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

var userInfo *tg.ChatFullInfo
var StopChannel = make(chan struct{})
var database *db.Database
var user *db.User
var index = 0
var photoPrices = []string{"1000", "2000"}
var phrasesPrices = []string{"150", "300"}
var phrases = []string{
	"Бу! Испугался? Не бойся",
	"Я РУССКИЙ!",
}

var chooseKeyboard = tu.Keyboard(
	tu.KeyboardRow(
		tu.KeyboardButton("Выбрать аватарку"),
		tu.KeyboardButton("Выбрать описание"),
	)).WithResizeKeyboard().WithOneTimeKeyboard()

var confirmKeyboard = tu.Keyboard(
	tu.KeyboardRow(
		tu.KeyboardButton("Подтвердить номер телефон").WithRequestContact(),
	)).WithResizeKeyboard()

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
				CallbackData: photoPrices[index],
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
