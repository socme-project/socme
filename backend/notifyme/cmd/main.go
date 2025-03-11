package main

import (
	"backend/notifyme"
	"fmt"
)

var telegramToken = "7640444898:AAGas29ZrF6vV8Wx8iSbyF-rNlR_QZdJShY"

func main() {
	tel := notifyme.Telegram(0, telegramToken)

	test, err := tel.HelperGetChatIds()
	if err != nil {
		fmt.Println(err)
	}

	chatid := test["hadrienaka"]
	if chatid == 0 {
		fmt.Println("ChatId not found")
	}

	tel = notifyme.Telegram(chatid, telegramToken)

	notifiers := notifyme.Notifiers{tel}

	errs := notifiers.Notify("Test", "Hello World")
	if errs != nil {
		fmt.Println(errs)
	}
}
