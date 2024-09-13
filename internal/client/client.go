package client

import (
	"fmt"
	"log/slog"
	"scrapNews/internal/models"
	"scrapNews/internal/tg_client"
)

type ClientMessanger interface {
	Send(sendingMessage models.SendingMessage) error
	WebhookEvent(data map[string]interface{}, log *slog.Logger) (models.WebhookEvent, error)
}

type Client struct {
	ClientMessangers map[string]ClientMessanger
	Log              *slog.Logger
}

func New(tghost string, tgtoken string, log *slog.Logger) *Client {
	c := &Client{
		ClientMessangers: make(map[string]ClientMessanger),
	}
	telegramClient := tg_client.New(tghost, tgtoken)
	c.ClientMessangers[models.Telegram] = telegramClient
	log.Info("start client")
	return c
}

func (c *Client) Send(sendingMessage models.SendingMessage) error {
	sender, ok := c.ClientMessangers[sendingMessage.Messenger.Name]
	if !ok {
		return fmt.Errorf("messenger %s not found", sendingMessage.Messenger)
	}

	err := sender.Send(sendingMessage)
	if err != nil {
		c.Log.Info("failed to send message: %w", err)
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (c *Client) WebhookEvent(data map[string]interface{}, messenger models.Messenger) (models.WebhookEvent, error) {
	recepient, ok := c.ClientMessangers[messenger.Name]
	if !ok {

		return models.WebhookEvent{}, fmt.Errorf("messenger %s not found", messenger)
	}
	webhookEvent, err := recepient.WebhookEvent(data, c.Log)
	if err != nil {
		return models.WebhookEvent{}, fmt.Errorf("messenger %s cant process the request ", messenger)
	}
	return webhookEvent, nil
}
