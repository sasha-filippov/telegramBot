package telegram

import (
	"errors"
	"telegramBot/clients/telegramClient"
	"telegramBot/events"
	"telegramBot/lib/e"
	"telegramBot/storage"
)

type Dispatcher struct {
	tg      *telegramClient.Client
	offset  int
	storage storage.Storage
}

var ErrUnknownEventType = errors.New("Unknown event type")
var ErrUnknownMetaType = errors.New("Unknown meta type")

func New(client *telegramClient.Client, storage storage.Storage) *Dispatcher {
	return &Dispatcher{
		tg:      client,
		storage: storage,
	}
}

type Meta struct {
	ChatId   int
	UserName string
}

func (d *Dispatcher) Fetcher(limit int) ([]events.Event, error) {
	updates, err := d.tg.Updates(d.offset, limit)
	if err != nil {
		return nil, err
	}
	if len(updates) == 0 {
		return nil, nil
	}
	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}
	d.offset = updates[len(updates)-1].ID + 1

	return res, nil

}

func (d *Dispatcher) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return d.processMessage(event)
	default:
		return e.Wrap("can't process a message", ErrUnknownEventType)
	}
}

func (d *Dispatcher) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}
	if err := d.doCmd(event.Text, meta.ChatId, meta.UserName); err != nil {
		return e.Wrap("can't process message", err)
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get Meta", ErrUnknownMetaType)
	}
	return res, nil

}

func event(upd telegramClient.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}
	// chatID, userName
	if updType == events.Message {
		res.Meta = Meta{
			ChatId:   upd.Message.Chat.ID,
			UserName: upd.Message.From.UserName,
		}
	}

	return res

}

func fetchText(upd telegramClient.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd telegramClient.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}
