package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/sol1corejz/read-adviser-bot/lib/e"
	"io"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	SendAll(userName string) ([]*Page, error)
	Remove(p *Page) error
	IsExist(p *Page) (bool, error)
}

var ErrNoSavedPages = errors.New("no saved pages")

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can`t calculate hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can`t calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
