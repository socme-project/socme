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

func (t *telegram) Notify(title, message string) error {
	chatId := fmt.Sprintf("%d", t.ChatId)
	msg := ""
	if title != "" {
		msg += "*" + title + ":*\n"
	}
	msg += message
	response, err := http.Post(
		"https://api.telegram.org/bot"+t.Token+"/sendMessage",
		"application/json",
		bytes.NewBufferString(
			`{"chat_id": "`+chatId+`", "text": "`+msg+`", "parse_mode": "markdown"}`,
		),
	)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return err
	}

	return nil
}

func (t *telegram) HelperGetChatIds() (map[string]int, error) {
	response, err := http.Get("https://api.telegram.org/bot" + t.Token + "/getUpdates")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	type responseType struct {
		Result []struct {
			Message struct {
				From struct {
					Username string `json:"username"`
				}
				Chat struct {
					Id int `json:"id"`
				}
			}
		}
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var result responseType
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if len(result.Result) == 0 {
		return nil, fmt.Errorf("No chat ID found")
	}

	var chatIds map[string]int = make(map[string]int)

	for _, r := range result.Result {
		chatIds[r.Message.From.Username] = r.Message.Chat.Id
	}

	return chatIds, nil
}

func (t *telegram) HelperGetChatIdByUsername(username string) ([]int, error) {
	chatIds, err := t.HelperGetChatIds()
	if err != nil {
		return nil, err
	}

	var chatIdsInt []int

	for k, v := range chatIds {
		if k == username {
			chatIdsInt = append(chatIdsInt, v)
		}
	}

	return chatIdsInt, nil
}
