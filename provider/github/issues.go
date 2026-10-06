package github

import (
	"encoding/json"
	"fmt"
	"log"

	"qnotifier/discord"
	"qnotifier/provider"

	"github.com/tidwall/gjson"
)

var issue int64

func Issues() {
	notifier := provider.Notifier{
		API:     "https://api.github.com/repos/qtile/qtile/issues",
		Header:  header.Clone(),
		Webhook: discord.Webhook["issues"],
	}

	if issue == 0 {
		notifier.API += "?per_page=1"
	} else {
		notifier.API += fmt.Sprintf("/%d", issue+1)
	}

	data, err := notifier.Fetch(issue > 0)
	if err != nil {
		log.Printf("error (issues:fetch): %v", err)
		return
	}

	if data == nil {
		nextIssue := issue + 1
		if isDiscussion(nextIssue) && sendDiscussion(nextIssue) {
			issue += 1
			return
		}
		return
	}

	if issue == 0 {
		res := gjson.GetBytes(data, "0.number")
		if !res.Exists() {
			log.Print("error (issues:fetch): issue number could not be found")
			return
		}
		issue = res.Int()
		log.Printf("issues: #%d has been saved", issue)
		return
	}

	payload, patch := getIssuePayload(string(data))
	if payload == nil {
		return
	}

	if gjson.GetBytes(data, "pull_request").Exists() {
		notifier.Webhook = discord.Webhook["pulls"]
		if patch != nil {
			patch.Webhook = "pulls"
		}
	}

	notifier.Payload = payload
	id, err := notifier.Send(true)
	if err != nil {
		log.Printf("error (issues:send): %v", err)
		return
	}

	number := gjson.GetBytes(data, "number").String()
	if patch != nil {
		toPatch[id] = patch
		log.Printf("issues: a patch for #%s was created (ID: %s)", number, id)
	}

	issue += 1
	log.Printf("issues: #%s has been sent", number)
}

func getIssuePayload(data string) ([]byte, *savePatch) {
	number := gjson.Get(data, "number").String()
	title := "Issue opened: #" + number
	color := 0xeb6420
	itemType := "issues"

	if gjson.Get(data, "pull_request").Exists() {
		title = "Pull request opened: #" + number
		color = 0x7289da
		itemType = "pull"
	}

	updatedAt := gjson.Get(data, "updated_at").Time().Unix()
	openGraphURL := getOpenGraphURL(updatedAt, itemType, number)

	imageURL := openGraphURL
	if !isImage(openGraphURL) {
		log.Printf("issues (#%s): image could not be loaded, using fallback", number)
		imageURL = fallbackImage
	}

	payload := discord.Payload{
		Username: "GitHub",
		Embeds: []discord.Embed{
			{
				Title: title,
				URL:   gjson.Get(data, "html_url").String(),
				Color: color,
				Image: &discord.EmbedImage{
					URL: imageURL,
				},
				Author: &discord.EmbedAuthor{
					Name: gjson.Get(data, "user.login").String(),
					Icon: gjson.Get(data, "user.avatar_url").String(),
					URL:  gjson.Get(data, "user.html_url").String(),
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error (issues:json): %v", err)
		return nil, nil
	}

	if imageURL != fallbackImage {
		return body, nil
	}

	return body, &savePatch{
		Data:     body,
		ImageURL: openGraphURL,
		Webhook:  "issues",
	}
}
