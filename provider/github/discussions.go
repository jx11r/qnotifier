package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"qnotifier/discord"
	"qnotifier/provider"

	"github.com/tidwall/gjson"
)

type GraphQLRequest struct {
	Query     string    `json:"query"`
	Variables Variables `json:"variables"`
}

type Variables struct {
	Owner  string `json:"owner"`
	Repo   string `json:"repo"`
	Number int64  `json:"number"`
}

const query = `query($number: Int!) {
	repository(owner: "qtile", name: "qtile") {
		discussion(number: $number) {
			title url updatedAt author { login url avatarUrl }
		}
	}
}`

func isDiscussion(number int64) bool {
	url := fmt.Sprintf("https://github.com/qtile/qtile/discussions/%d", number)
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		log.Printf("error (github): %v", err)
		return false
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := provider.Client.Do(req)
	if err != nil {
		log.Printf("error (github): %v", err)
		return false
	}

	if resp.StatusCode == http.StatusOK {
		return true
	}

	return false
}

func fetchDiscussion(number int64) ([]byte, error) {
	api := "https://api.github.com/graphql"

	payload := GraphQLRequest{
		Query: query,
		Variables: Variables{
			Number: number,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, api, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header = header.Clone()

	resp, err := provider.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s returned %s", api, resp.Status)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func sendDiscussion(number int64) bool {
	notifier := provider.Notifier{
		Webhook: discord.Webhook["discussions"],
	}

	data, err := fetchDiscussion(number)
	if err != nil {
		log.Printf("error (discussions:fetch): %v", err)
		return false
	}

	payload, patch := getDiscussionPayload(string(data), number)
	if payload == nil {
		return false
	}

	notifier.Payload = payload
	id, err := notifier.Send(true)
	if err != nil {
		log.Printf("error (discussions:send): %v", err)
		return false
	}

	if patch != nil {
		toPatch[id] = patch
		log.Printf("discussions: a patch for #%d was created (ID: %s)", number, id)
	}

	log.Printf("discussions: #%d has been sent", number)
	return true
}

func getDiscussionPayload(data string, number int64) ([]byte, *savePatch) {
	discussion := gjson.Get(data, "data.repository.discussion")
	updatedAt := discussion.Get("updatedAt").Time().Unix()
	openGraphURL := getOpenGraphURL(updatedAt, "discussions", strconv.FormatInt(number, 10))

	imageURL := openGraphURL
	if !isImage(openGraphURL) {
		log.Printf("discussions (#%d): image could not be loaded, using fallback", number)
		imageURL = fallbackImage
	}

	payload := discord.Payload{
		Username: "GitHub",
		Embeds: []discord.Embed{
			{
				Title: fmt.Sprintf("Discussion opened: #%d", number),
				URL:   discussion.Get("url").String(),
				Color: 0xe68d60,
				Image: &discord.EmbedImage{
					URL: imageURL,
				},
				Author: &discord.EmbedAuthor{
					Name: discussion.Get("author.login").String(),
					Icon: discussion.Get("author.avatarUrl").String(),
					URL:  discussion.Get("author.url").String(),
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error (discussions:json): %v", err)
		return nil, nil
	}

	if imageURL != fallbackImage {
		return body, nil
	}

	return body, &savePatch{
		Data:     body,
		ImageURL: openGraphURL,
		Webhook:  "discussions",
	}
}
