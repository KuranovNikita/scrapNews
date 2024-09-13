package models

const (
	Telegram           = "telegram"
	OpenMenu           = "Открыть меню"
	Return             = "Назад"
	FollowSport        = "Подписаться Спорт"
	UnFollowSport      = "Отписаться Спорт"
	FollowKommersant   = "Подписаться Коммерсант"
	UnFollowKommersant = "Отписаться Коммерсант"
	Start              = "/start"
	Kommersant         = "Коммерсант"
	Sport              = "СпортЭкспресс"
)

type SendingMessage struct {
	ID        string
	Text      string
	Messenger Messenger
	Buttons   []Button
}

type Messenger struct {
	Name string
}

type Button struct {
	Text string
}

type WebhookEvent struct {
	EventType string
	ChatId    string
	Text      string
	Data      ClickBtnEventData
}

// type EventData interface {
// 	Type() string
// }

type ClickBtnEventData struct {
	IsOpenMenu     bool
	IsCloseMenu    bool
	IsChangeFollow bool
	FollowData     FollowData
}

// func (cbd *ClickBtnEventData) Type() string { return "clickButton" }

// type FollowDataInterface interface {
// 	GetNameSite() string
// 	GetIsFollow() bool
// }

type FollowData struct {
	NameSite string
	IsFollow bool
}

// func (fd *FollowData) GetNameSite() string {
// 	return fd.NameSite
// }

// func (fd *FollowData) GetIsFollow() bool {
// 	return fd.IsFollow
// }
