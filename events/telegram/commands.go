package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"
	"telegramBot/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (d *Dispatcher) doCmd(text string, chatId int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got a new command '%s' from '%s'", text, username)
	if isAddCmd(text) {
		return d.savePage(chatId, text, username)
	}
	// add a page: http://...
	// rnd page: /rnd
	//help: /help
	// start: /start: hi + help
	switch text {

	case RndCmd:
		return d.sendRandom(chatId, username)
	case HelpCmd:
		return d.sendHelp(chatId)
	case StartCmd:
		return d.sendHello(chatId)
	default:
		return d.tg.SendMessage(chatId, msgUnknownCommand)

	}
}

func (d *Dispatcher) savePage(chatId int, pageUrl string, username string) (err error) {
	page := &storage.Page{
		URL:      pageUrl,
		UserName: username,
	}
	isExists, err := d.storage.IsExists(page)
	if err != nil {
		return errors.New("existing error")
	}
	if isExists {
		return d.tg.SendMessage(chatId, msgAlreadyExists)
	}
	if err := d.storage.Save(page); err != nil {
		return err //errors.New("saving error")
	}
	if err := d.tg.SendMessage(chatId, msgSaved); err != nil {
		return errors.New("sending error")
	}
	return nil
}

func (d *Dispatcher) sendRandom(chatId int, username string) (err error) {
	page, err := d.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return d.tg.SendMessage(chatId, msgNoSavedPages)
	}
	if err := d.tg.SendMessage(chatId, page.URL); err != nil {
		return errors.New("couldnt send a page")
	}
	return d.storage.Remove(page)
}

func (d *Dispatcher) sendHelp(chatId int) error {
	return d.tg.SendMessage(chatId, msgHelp)
}
func (d *Dispatcher) sendHello(chatId int) error {
	return d.tg.SendMessage(chatId, msgHello)
}

func isAddCmd(text string) bool {
	return isUrl(text)

}

func isUrl(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
