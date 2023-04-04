package provider

import (
	"bytes"
	"io"
	"net/http"

	"github.com/jx11r/qnotifier/utils"
)

type Base interface {
	Send() error
	Fetch() ([]byte, error)
}

type Notifier struct {
	Base
	API     string
	Payload []byte
	Webhook string
}

func (n *Notifier) Send() error {
	req, err := http.NewRequest("POST", n.Webhook, bytes.NewBuffer(n.Payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (n *Notifier) Fetch() ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", n.API, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"User-Agent":    {"jx11r"},
		"Authorization": {"token " + utils.Token["github"]},
		"Accept":        {"application/vnd.github.v3+json"},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
