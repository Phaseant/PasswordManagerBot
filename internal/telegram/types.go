package telegram

type Update struct {
	UpdateId int              `json:"update_id"`
	Message  *IncomingMessage `json:"message"`
}

type UpdatesResponse struct {
	Ok      bool     `json:"ok"`
	Updates []Update `json:"result"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}

type ReplyKeyboardMarkup struct {
	Keyboard     [][]KeyboardButton `json:"keyboard"`
	IsPersistent bool               `json:"is_persistent"`
}

type KeyboardButton struct {
	Text string `json:"text"`
}
