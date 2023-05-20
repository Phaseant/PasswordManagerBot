package telegramEvents

import (
	"errors"

	"github.com/Phaseant/PasswordManagerBot/internal/events"
	"github.com/Phaseant/PasswordManagerBot/internal/repository"
	"github.com/Phaseant/PasswordManagerBot/internal/telegram"
)

type Meta struct {
	ChatID   int
	Username string
}

type Processor struct {
	tg     *telegram.Client
	offset int
	repo   *repository.Repository
}

func New(tg *telegram.Client, repo *repository.Repository) *Processor {
	return &Processor{
		tg:     tg,
		offset: 0,
		repo:   repo,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, errors.New("failed to fetch updates: " + err.Error())
	}

	res := make([]events.Event, len(updates))

	for _, u := range updates {
		event := event(u)
		if u.Message != nil {
			event.MessageID = u.Message.MessageID //for msg delete
		}
		res = append(res, event)
	}

	if len(updates) > 0 {
		p.offset = updates[len(updates)-1].UpdateId + 1
	} else {
		p.offset = 0
	}

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return errors.New("unknown event type")

	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return errors.New("failed to get meta: " + err.Error())
	}
	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return errors.New("failed to process message: " + err.Error())
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	meta, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, errors.New("failed to cast meta to Meta")
	}

	return meta, nil
}

func event(upd telegram.Update) events.Event {
	event := events.Event{
		Type: fetchType(upd),
		Text: fetchText(upd),
	}

	if event.Type == events.Message {
		event.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}
	return event
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}
