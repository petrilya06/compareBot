package bot

import (
	"CompareBot/db"
	"fmt"
	tg "github.com/mymmrac/telego"
	"github.com/vitali-fedulov/images3"
	"log"
	"os"
	"time"
)

func OpenFile(filename string) *os.File {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return file
}

func comparePhotos(path1, path2 string) bool {
	img1, _ := images3.Open(path1)
	img2, _ := images3.Open(path2)
	icon1 := images3.Icon(img1, path1)
	icon2 := images3.Icon(img2, path2)

	if images3.Similar(icon1, icon2) {
		return true
	} else {
		return false
	}
}

func CheckPhotos(bot *tg.Bot, user *db.User) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			downloadPhoto(bot, user)
			path1 := fmt.Sprintf("src/%d.jpg", user.SelectPic)
			path2 := fmt.Sprintf("src/photos/%d.jpg", user.TgID)

			if comparePhotos(path1, path2) {
				user.CountPhotoCompare += 1
				fmt.Println(user)

				if user.CountPhotoCompare >= 24 {
					// TODO сделать свою систему выплату
					SendMessage(bot, user, emptyKeyboard, "Тебе назначена выплата!")
					user.CountPhotoCompare = 0
				}
			} else {
				if user.CountPhotoCompare != 0 {
					SendMessage(bot, user, emptyKeyboard, "Ваша аватарка не совпадает!")
				}

				user.CountPhotoCompare = 0
			}

			if err := database.UpdateUser(*user); err != nil {
				log.Println(err)
			}

		case <-StopChannel:
			return
		}
	}
}

func StopChecks() {
	close(StopChannel)
}
