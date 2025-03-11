package notifyme

import (
	"errors"
	"io"
	"net/http"
	"strings"
)

type ntfysh struct {
	url   string // Full URL to the ntfysh server (http://localhost:8080)
	topic string
}

func Ntfysh(url, topic string) *ntfysh {
	// parse the url
	return &ntfysh{topic: topic, url: url}
}

func (n *ntfysh) Notify(title, message string) error {
	req, _ := http.NewRequest("POST", n.url+"/"+n.topic,
		strings.NewReader(message))
	if title != "" {
		req.Header.Set("Title", title)
	}
	// req.Header.Set("Priority", "urgent")
	// req.Header.Set("Tags", "warning,skull")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return errors.New("Error: " + string(body))
	}
	return nil
}
