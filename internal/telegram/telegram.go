package telegram

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
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

// get updates from telegram
func (c *Client) Updates(offset, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doReqWithTimeout(q)
	if err != nil {
		return nil, errors.New("failed to do request: " + err.Error())
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, errors.New("failed to unmarshal response: " + err.Error())
	}

	return res.Updates, nil

}

// func to do request with exponential backoff
func (c *Client) doReqWithTimeout(q url.Values) ([]byte, error) {
	maxRetries := 5
	retryDelay := 1 * time.Second
	var data []byte
	var err error

	for i := 0; i < maxRetries; i++ {
		data, err = c.doRequest("getUpdates", q) //do request
		if err == nil {
			return data, nil
		}
		log.Print("failed to do request: " + err.Error() + ", retrying...")
		delay := time.Duration(1<<uint(i)) * retryDelay        //delay
		jitter := time.Duration(rand.Int63n(int64(delay / 2))) //jitter to avoid synchronized requests
		sleepTime := delay + jitter
		time.Sleep(sleepTime)
	}

	return nil, errors.New("failed to do request: " + err.Error())
}

// base func to send requests
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

// send message to telegram
func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest("sendMessage", q)
	if err != nil {
		return errors.New("failed to send message: " + err.Error())
	}
	return nil
}

// delete message from telegram
func (c *Client) DeleteMessage(chatID int, msgId int) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("message_id", strconv.Itoa(msgId))

	_, err := c.doRequest("deleteMessage", q)
	if err != nil {
		return errors.New("failed to delete message: " + err.Error())
	}
	return nil
}
