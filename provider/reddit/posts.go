package reddit

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strings"
	"time"

	"qnotifier/discord"
	"qnotifier/provider"
)

var lastPublishedAt time.Time

func Posts() {
	notifier := provider.Notifier{
		API:     "https://www.reddit.com/r/qtile/new.rss?limit=1",
		Webhook: discord.Webhook["reddit"],
		Header: http.Header{
			"User-Agent": {"script:qnotifier:v1.0 (by /u/jx11r)"},
		},
	}

	data, err := notifier.Fetch()
	if err != nil {
		log.Printf("error (reddit:fetch): %v", err)
		return
	}

	var feed Feed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		log.Printf("error (reddit:xml): %v", err)
		return
	}

	if len(feed.Entries) == 0 {
		log.Println("reddit: no entries found in feed")
		return
	}
	entry := feed.Entries[0]

	if lastPublishedAt.IsZero() {
		lastPublishedAt, _ = time.Parse(time.RFC3339, entry.Published)
		log.Printf("reddit: latest post has been saved (%s)", entry.Link.Href)
		return
	}

	publishedAt, _ := time.Parse(time.RFC3339, entry.Published)
	if !publishedAt.After(lastPublishedAt) {
		return
	}

	payload := getPayload(entry)
	if payload == nil {
		return
	}

	notifier.Payload = payload
	_, err = notifier.Send(false)
	if err != nil {
		log.Printf("error (reddit:send): %v", err)
		return
	}

	lastPublishedAt = publishedAt
	log.Printf("reddit: a new post has been sent (%s)", entry.Link.Href)
}

func getPayload(entry Entry) []byte {
	icon_url := fmt.Sprintf(
		"https://www.redditstatic.com/avatars/defaults/v2/avatar_default_%d.png",
		rand.IntN(8),
	)

	payload := discord.Payload{
		Username: "Reddit",
		Embeds: []discord.Embed{
			{
				Title: entry.Title,
				URL:   entry.Link.Href,
				Color: 0xff4400,
				Author: &discord.EmbedAuthor{
					Name: strings.TrimPrefix(entry.Author.Name, "/u/"),
					Icon: icon_url,
					URL:  entry.Author.URI,
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error (reddit:json): %v", err)
		return nil
	}

	return body
}
