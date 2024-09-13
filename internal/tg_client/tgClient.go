package tg_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"scrapNews/internal/models"
	"strconv"
)

const (
	RemoveKeyboard = true
)

type tgClient struct {
	host     string
	basePath string
}

type sendText struct {
	Chat_id string `json:"chat_id"`
	Text    string `json:"text"`
}

type KeyboardButton struct {
	Text string `json:"text"`
}

type ReplyMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard,omitempty"`
	ResizeKeyboard  bool               `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard bool               `json:"one_time_keyboard,omitempty"`
	RemoveKeyboard  bool               `json:"remove_keyboard,omitempty"`
}

type sendTextWithButtons struct {
	sendText
	ReplyMarkup ReplyMarkup `json:"reply_markup"`
}

func New(host string, token string) *tgClient {
	return &tgClient{
		host:     host,
		basePath: newBasePath(token),
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (tgClient *tgClient) Send(sendingMessage models.SendingMessage) error {
	sendText := sendText{
		Chat_id: sendingMessage.ID,
		Text:    sendingMessage.Text,
	}
	replyMarkup := createReplyMarkup(sendingMessage.Buttons)
	data := sendTextWithButtons{
		sendText,
		replyMarkup,
	}
	err := sendMessageToTelegram(tgClient, data)
	if err != nil {
		return err
	}
	return nil

}

func buildURL(tgClient *tgClient) string {
	fullPath := fmt.Sprintf("https://%s", path.Join(tgClient.host, tgClient.basePath, "sendMessage"))
	return fullPath
}

func sendMessageToTelegram(tgClient *tgClient, data interface{}) error {
	url := buildURL(tgClient)

	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return nil
}

func createReplyMarkup(buttons []models.Button) ReplyMarkup {
	if len(buttons) > 0 {
		keyboard := make([][]KeyboardButton, 0)

		for _, btn := range buttons {
			keyboard = append(keyboard, []KeyboardButton{{Text: btn.Text}})
		}

		return ReplyMarkup{
			Keyboard:        keyboard,
			ResizeKeyboard:  true,
			OneTimeKeyboard: true,
		}
	} else {
		return ReplyMarkup{
			RemoveKeyboard: true,
		}
	}

}

func (tgClient *tgClient) WebhookEvent(data map[string]interface{}, log *slog.Logger) (models.WebhookEvent, error) {
	webhookEvent := models.WebhookEvent{}
	message, ok := data["message"].(map[string]interface{})
	if !ok {
		log.Error("Error extracting 'message' field")
		return models.WebhookEvent{}, errors.New("Error extracting 'message' field")
	}

	// formattedJSON, err := json.MarshalIndent(data, "", "  ")
	// if err != nil {
	// log.Error("Error formatting JSON")
	// return models.WebhookEvent{}, errors.New("Error formatting JSON")
	// }

	// fmt.Println(string(formattedJSON))
	from, ok := message["from"].(map[string]interface{})
	if !ok {
		log.Error("Error extracting 'from' field")
		return models.WebhookEvent{}, errors.New("Error extracting 'from' field")
	}
	id, ok := from["id"].(float64) // id обычно приходит как float64
	if !ok {
		log.Error("Error extracting 'id' field")
		return models.WebhookEvent{}, errors.New("Error extracting 'id' field")
	}
	userIDStr := strconv.FormatFloat(id, 'f', -1, 64)
	text, ok := message["text"].(string)
	if !ok {
		log.Error("Error extracting 'text' field")
		return models.WebhookEvent{}, errors.New("Error extracting 'text' field")
	}
	webhookEvent.Text = text
	webhookEvent.ChatId = userIDStr

	switch text {
	case models.OpenMenu:
		webhookEvent.EventType = "clickButton"
		webhookEvent.Data.IsOpenMenu = true
		webhookEvent.Data.IsCloseMenu = false
		webhookEvent.Data.IsChangeFollow = false
	case models.Return:
		webhookEvent.EventType = "clickButton"
		webhookEvent.Data.IsOpenMenu = false
		webhookEvent.Data.IsCloseMenu = true
		webhookEvent.Data.IsChangeFollow = false
	case models.FollowSport:
		webhookEvent.EventType = "clickButton"
		webhookEvent.Data.IsOpenMenu = false
		webhookEvent.Data.IsCloseMenu = false
		webhookEvent.Data.IsChangeFollow = true
		followData := models.FollowData{
			NameSite: models.Sport,
			IsFollow: true,
		}
		webhookEvent.Data.FollowData = followData
	case models.FollowKommersant:
		webhookEvent.EventType = "clickButton"
		webhookEvent.Data.IsOpenMenu = false
		webhookEvent.Data.IsCloseMenu = false
		webhookEvent.Data.IsChangeFollow = true
		followData := models.FollowData{
			NameSite: models.Kommersant,
			IsFollow: true,
		}
		webhookEvent.Data.FollowData = followData
	case models.UnFollowKommersant:
		webhookEvent.EventType = "clickButton"
		webhookEvent.Data.IsOpenMenu = false
		webhookEvent.Data.IsCloseMenu = false
		webhookEvent.Data.IsChangeFollow = true
		followData := models.FollowData{
			NameSite: models.Kommersant,
			IsFollow: false,
		}
		webhookEvent.Data.FollowData = followData
	case models.UnFollowSport:
		webhookEvent.EventType = "clickButton"
		webhookEvent.Data.IsOpenMenu = false
		webhookEvent.Data.IsCloseMenu = false
		webhookEvent.Data.IsChangeFollow = true
		followData := models.FollowData{
			NameSite: models.Sport,
			IsFollow: false,
		}
		webhookEvent.Data.FollowData = followData
	case models.Start:
		webhookEvent.EventType = "botCommand"
	default:
		webhookEvent.EventType = "text"
	}

	fmt.Println("webhookEvent %v", webhookEvent) //!!!!!!!!!

	return webhookEvent, nil
}
