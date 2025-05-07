package telegram

import (
	"errors"
	telegram "github.com/sol1corejz/read-adviser-bot/internal/clients"
	"github.com/sol1corejz/read-adviser-bot/lib/e"
	"github.com/sol1corejz/read-adviser-bot/storage"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatId int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command: %s from %s", text, username)

	if isAddCmd(text) {
		return p.savePage(chatId, text, username)
	}

	switch text {
	case HelpCmd:
		return p.sendHelp(chatId)
	case StartCmd:
		return p.sendHello(chatId)
	case RndCmd:
		return p.sendRandom(chatId, username)
	default:
		return p.tg.SendMessage(chatId, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageUrl string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can`t do command: save page", err) }()

	sendMsg := NewMessageSender(chatID, p.tg)

	page := &storage.Page{
		URL:      pageUrl,
		UserName: username,
	}

	isExists, err := p.storage.IsExist(page)
	if err != nil {
		return err
	}

	if isExists {
		return sendMsg(msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := sendMsg(msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can`t do command: can`t send random page", err) }()

	sendMsg := NewMessageSender(chatID, p.tg)

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return sendMsg(msgNoSavedPages)
	}

	if err := sendMsg(page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) (err error) {
	sendMsg := NewMessageSender(chatID, p.tg)

	return sendMsg(msgHelp)
}

func (p *Processor) sendHello(chatID int) (err error) {
	sendMsg := NewMessageSender(chatID, p.tg)

	return sendMsg(msgHello)
}

func NewMessageSender(cahtId int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(cahtId, msg)
	}
}

func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
