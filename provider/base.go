package provider

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

var (
	Client        = &http.Client{Timeout: 5 * time.Second}
	WebhookClient = &http.Client{Timeout: 10 * time.Second}
)

type Notifier struct {
	API     string
	Header  http.Header
	Payload []byte
	Webhook string
}

func (n *Notifier) Fetch(allowNotFound ...bool) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, n.API, nil)
	if err != nil {
		return nil, err
	}

	req.Header = n.Header.Clone()

	resp, err := Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if len(allowNotFound) > 0 && allowNotFound[0] {
		if resp.StatusCode == http.StatusNotFound {
			return nil, nil
		}
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		reset := resp.Header.Get("X-Ratelimit-Reset")
		return nil, fmt.Errorf("%s returned %s (retry in %ss)", n.API, resp.Status, reset)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s returned %s", n.API, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (n *Notifier) Send(returnID bool) (string, error) {
	if returnID {
		n.Webhook += "?wait=true"
	}

	req, err := http.NewRequest(http.MethodPost, n.Webhook, bytes.NewBuffer(n.Payload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := WebhookClient.Do(req)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusNoContent && !returnID {
		return "", nil
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to send webhook, discord returned %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if !gjson.GetBytes(body, "id").Exists() {
		return "", fmt.Errorf("message id doesn't exist")
	}

	return gjson.GetBytes(body, "id").String(), nil
}
