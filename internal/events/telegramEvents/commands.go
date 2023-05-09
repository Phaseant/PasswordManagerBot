package telegramEvents

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Phaseant/PasswordManagerBot/internal/events"
)

const (
	START  = "/start"
	HELP   = "/help"
	ADD    = "/add"
	DELETE = "/delete"
	GET    = "/get"
)

var ErrNotFound = "no documents found for this service"

// Message router
func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)
	log.Printf("got command: %s", text)

	switch text {
	case START:
		return p.tg.SendMessage(chatID, StartMessage)
	case HELP:
		return p.tg.SendMessage(chatID, HelpMessage)
	case ADD:
		if err := p.doAdd(chatID, username); err != nil {
			return p.tg.SendMessage(chatID, AddErrorMessage)
		}
		return p.tg.SendMessage(chatID, AddSuccessMessage)
	case DELETE:
		if err := p.doDelete(chatID, username); err != nil {
			return p.tg.SendMessage(chatID, DeleteErrorMessage)
		}
		return p.tg.SendMessage(chatID, DeleteSuccessMessage)
	case GET:
		if err := p.doGet(chatID, username); err != nil {
			return p.tg.SendMessage(chatID, GetErrorMessage)
		}
		return nil
	default:
		return p.tg.SendMessage(chatID, unknownMessage)
	}
}

// /add command
func (p *Processor) doAdd(chatID int, userID string) error {
	const inputService = "Введите название сервиса"
	const inputUsername = "Введите логин"
	const inputPassword = "Введите пароль"
	var service, username, password string
	service, err := p.sendMsgAndProcess(inputService, chatID)
	if err != nil {
		return err
	}
	username, err = p.sendMsgAndProcess(inputUsername, chatID)
	if err != nil {
		return err
	}
	password, err = p.sendMsgAndProcess(inputPassword, chatID)
	if err != nil {
		return err
	}
	err = p.repo.Add(service, username, password, userID)
	if err != nil {
		return err
	}
	return nil
}

// /delete command
func (p *Processor) doDelete(chatID int, userID string) error {
	const inputService = "Введите название сервиса"
	const inputUsername = "Введите логин"
	service, err := p.sendMsgAndProcess(inputService, chatID)
	if err != nil {
		return err
	}
	username, err := p.sendMsgAndProcess(inputUsername, chatID)
	if err != nil {
		return err
	}
	err = p.repo.Delete(service, username, userID)
	if err != nil {
		return err
	}
	return nil
}

// /get command
func (p *Processor) doGet(chatID int, userID string) error {
	const inputService = "Введите название сервиса"
	service, msgID, err := p.sendMsgAndProcessWithMsgID(inputService, chatID)
	if err != nil {
		return err
	}
	logsAndPass, err := p.repo.Get(service, userID)
	if err != nil {
		if err.Error() == ErrNotFound {
			return p.tg.SendMessage(chatID, NotFoundMessage)
		}
		return err
	}
	if len(logsAndPass) < 1 {
		return p.tg.SendMessage(chatID, GetErrorMessage)
	}
	for _, logAndPass := range logsAndPass {
		result := fmt.Sprintf("Логин: %s\nПароль: %s", logAndPass.Login, logAndPass.Password)
		msgID += 1 //move delete cursor to our message
		err = p.tg.SendMessage(chatID, result)
		go p.deleteMessage(chatID, msgID)
	}
	if err != nil {
		return err
	}
	msgID += 1 //move delete cursor to our message
	p.tg.SendMessage(chatID, "(Сообщение удалится через 30 секунд)")
	go p.deleteMessage(chatID, msgID)

	return nil

}

// delete message after 30 seconds
func (p *Processor) deleteMessage(chatID, msgID int) {
	c := time.NewTimer(time.Second * 30)
	<-c.C
	p.tg.DeleteMessage(chatID, msgID)
}

// send message and process response, needs to handle nested commands
func (p *Processor) sendMsgAndProcessWithMsgID(msg string, chatID int) (string, int, error) {
	p.tg.SendMessage(chatID, msg)
	got := false
	for !got {
		eventss, _ := p.Fetch(10)
		for _, event := range eventss {
			switch event.Type {
			case events.Message:
				variable := event.Text
				msgID := event.MessageID
				got = true
				return variable, msgID, nil
			}
		}
	}
	return "", 0, nil
}

// wrapper for sendMsgAndProcessWithMsgID without ID of message
func (p *Processor) sendMsgAndProcess(msg string, chatID int) (string, error) {
	variable, _, err := p.sendMsgAndProcessWithMsgID(msg, chatID)
	return variable, err
}
