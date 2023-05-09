package telegram

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest("getUpdates", q)
	if err != nil {
		return nil, errors.New("failed to get updates: " + err.Error())
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, errors.New("failed to unmarshal response: " + err.Error())
	}

	return res.Updates, nil

}

func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.New("failed to do request: " + err.Error())
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.New("failed to do request: " + err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read response body: " + err.Error())
	}

	return body, nil
}

func (c *Client) SendMessage(chatID int, text string, keyboard ReplyKeyboardMarkup) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)
	if keyboard.Keyboard != nil {
		keyboardJSON, err := json.Marshal(keyboard)
		if err != nil {
			return errors.New("failed to marshal keyboard: " + err.Error())
		}
		q.Add("reply_markup", string(keyboardJSON))
	}

	_, err := c.doRequest("sendMessage", q)
	if err != nil {
		return errors.New("failed to send message: " + err.Error())
	}
	return nil
}

func (c *Client) SendPicture(chatID int, imageURL string, keyboard ReplyKeyboardMarkup) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("photo", imageURL)
	if keyboard.Keyboard != nil {
		keyboardJSON, err := json.Marshal(keyboard)
		if err != nil {
			return errors.New("failed to marshal keyboard: " + err.Error())
		}
		q.Add("reply_markup", string(keyboardJSON))
	}

	_, err := c.doRequest("sendPhoto", q)
	if err != nil {
		return errors.New("failed to send picture: " + err.Error())
	}
	return nil
}
