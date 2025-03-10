package main

import (
	"backend/notifyme"
	"fmt"
)

var telegramToken = "7640444898:AAGas29ZrF6vV8Wx8iSbyF-rNlR_QZdJShY"

func main() {
	tel := notifyme.Telegram(0, telegramToken)

	test, err := tel.HelperGetChatId()
	fmt.Println(test, err)

	tel = notifyme.Telegram(test, telegramToken)
	tel.Notify("Hello from Go!")
}
