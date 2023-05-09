package telegramEvents

import (
	"log"
	"strings"

	"github.com/Phaseant/PasswordManagerBot/internal/telegram"
)

const (
	START   = "/start"
	HELP    = "/help"
	JOKE    = "Шутка 🤡"
	POEM    = "Стих 📖"
	SONG    = "Песня 🎵"
	PICTURE = "Картинка 🖼️"

	RUSSONG = "🇷🇺"
	ENGSONG = "🇺🇸"

	POEM1 = "Пушкин"
	POEM2 = "Есенин"

	JOKE1 = "Смешная"
	JOKE2 = "Детская"

	PIC1 = "Гофер"
	PIC2 = "Спанч-Боб"

	JOKECOMMAND    = "/joke"
	POEMCOMMAND    = "/poem"
	SONGCOMMAND    = "/song"
	PICTURECOMMAND = "/picture"

	gopherImageURL    = "https://telegra.ph/file/dd33d46494b0f3e45a997.jpg"
	spongeBobImageURL = "https://telegra.ph/file/368bb8b80b34c017fcb8c.jpg"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	JokeBtn := telegram.KeyboardButton{Text: JOKE}
	SongBtn := telegram.KeyboardButton{Text: SONG}
	PoemBtn := telegram.KeyboardButton{Text: POEM}
	PictureBtn := telegram.KeyboardButton{Text: PICTURE}

	defKeyboard := telegram.ReplyKeyboardMarkup{
		Keyboard: [][]telegram.KeyboardButton{
			[]telegram.KeyboardButton{JokeBtn, SongBtn},
			[]telegram.KeyboardButton{PoemBtn, PictureBtn}},
		IsPersistent: true}

	SongKeyboard := buildKeyboard(RUSSONG, ENGSONG)
	PoemKeyboard := buildKeyboard(POEM1, POEM2)
	ImgKeyboard := buildKeyboard(PIC1, PIC2)
	JokeKeyboard := buildKeyboard(JOKE1, JOKE2)

	text = strings.TrimSpace(text)

	log.Printf("got command: %s", text)

	switch text {
	case START:
		return p.tg.SendMessage(chatID, StartMessage, defKeyboard)
	case HELP:
		return p.tg.SendMessage(chatID, HelpMessage, defKeyboard)
	//categories
	case JOKE, JOKECOMMAND:
		return p.tg.SendMessage(chatID, "Выбери категорию", JokeKeyboard)
	case POEM, POEMCOMMAND:
		return p.tg.SendMessage(chatID, "Выбери автора", PoemKeyboard)
	case SONG, SONGCOMMAND:
		return p.tg.SendMessage(chatID, "Выбери язык песни", SongKeyboard)
	case PICTURE, PICTURECOMMAND:
		return p.tg.SendMessage(chatID, "Выбери категорию", ImgKeyboard)

	//subcategories
	case RUSSONG:
		return p.tg.SendMessage(chatID, RusSongMessage, defKeyboard)
	case ENGSONG:
		return p.tg.SendMessage(chatID, EngSongMessage, defKeyboard)

	case POEM1:
		return p.tg.SendMessage(chatID, PoemMessage1, defKeyboard)
	case POEM2:
		return p.tg.SendMessage(chatID, PoemMessage2, defKeyboard)

	case JOKE1:
		return p.tg.SendMessage(chatID, JokeMessage1, defKeyboard)
	case JOKE2:
		return p.tg.SendMessage(chatID, JokeMessage2, defKeyboard)

	case PIC1:
		return p.tg.SendPicture(chatID, gopherImageURL, defKeyboard)
	case PIC2:
		return p.tg.SendPicture(chatID, spongeBobImageURL, defKeyboard)

	default:
		return p.tg.SendMessage(chatID, unknownMessage, defKeyboard)
	}
}

func buildKeyboard(btn1text, btn2text string) telegram.ReplyKeyboardMarkup {
	btn1 := telegram.KeyboardButton{Text: btn1text}
	btn2 := telegram.KeyboardButton{Text: btn2text}

	return telegram.ReplyKeyboardMarkup{
		Keyboard: [][]telegram.KeyboardButton{
			[]telegram.KeyboardButton{btn1, btn2}},
		IsPersistent: true}
}
