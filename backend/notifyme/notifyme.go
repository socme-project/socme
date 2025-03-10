package notifyme

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type telegram struct {
	ChatId int
	Token  string
}

func Telegram(chatId int, token string) *telegram {
	return &telegram{ChatId: chatId, Token: token}
}

func (t *telegram) Notify(message string) error {
	chatId := fmt.Sprintf("%d", t.ChatId)
	response, err := http.Post(
		"https://api.telegram.org/bot"+t.Token+"/sendMessage",
		"application/json",
		bytes.NewBufferString(`{"chat_id": "`+chatId+`", "text": "`+message+`"}`),
	)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return err
	}
	return nil
}

func (t *telegram) HelperGetChatId() (int, error) {
	// BOT_TOKEN="YOURBOTTOKEN" curl -s "https://api.telegram.org/bot$BOT_TOKEN/getUpdates" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2
	response, err := http.Get("https://api.telegram.org/bot" + t.Token + "/getUpdates")
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	type responseType struct {
		Result []struct {
			Message struct {
				Chat struct {
					Id int `json:"id"`
				}
			}
		}
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	var result responseType
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}
	if len(result.Result) == 0 {
		return 0, fmt.Errorf("No chat ID found")
	}

	return responseType(result).Result[0].Message.Chat.Id, nil
}
