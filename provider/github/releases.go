package github

import (
	"encoding/json"
	"fmt"
	"log"

	"qnotifier/discord"
	"qnotifier/provider"

	"github.com/Masterminds/semver/v3"
	"github.com/tidwall/gjson"
)

var release string

func Releases() {
	notifier := provider.Notifier{
		API:     "https://api.github.com/repos/qtile/qtile/releases/latest",
		Header:  header.Clone(),
		Webhook: discord.Webhook["releases"],
	}

	data, err := notifier.Fetch()
	if err != nil {
		log.Printf("error (releases:fetch): %v", err)
		return
	}

	current := gjson.GetBytes(data, "tag_name").String()
	if release == "" {
		release = current
		log.Printf("releases: latest version has been saved (%s)", current)
		return
	}

	previous, _ := semver.NewVersion(release)
	latest, err := semver.NewVersion(current)
	if err != nil {
		log.Printf("error (releases): %v", err)
		return
	}

	if !latest.GreaterThan(previous) {
		return
	}

	payload, patch := getReleasePayload(string(data))
	if payload == nil {
		return
	}

	notifier.Payload = payload
	id, err := notifier.Send(true)
	if err != nil {
		log.Printf("error (releases:send): %v", err)
		return
	}

	if patch != nil {
		toPatch[id] = patch
		log.Printf("releases: a patch for %s was created (ID: %s)", current, id)
	}

	release = current
	log.Printf("releases: %s has been sent", current)
}

func getReleasePayload(data string) ([]byte, *savePatch) {
	version := gjson.Get(data, "tag_name").String()
	updatedAt := gjson.Get(data, "updated_at").Time().Unix()
	openGraphURL := getOpenGraphURL(updatedAt, "releases/tag", version)

	imageURL := openGraphURL
	if !isImage(openGraphURL) {
		log.Printf("releases (%s): image could not be loaded, using fallback", version)
		imageURL = fallbackImage
	}

	payload := discord.Payload{
		Username: "GitHub",
		Embeds: []discord.Embed{
			{
				Title: fmt.Sprintf("New release published: %s", version),
				URL:   gjson.Get(data, "html_url").String(),
				Color: 0x60c67b,
				Image: &discord.EmbedImage{
					URL: imageURL,
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error (releases:json): %v", err)
		return nil, nil
	}

	if imageURL != fallbackImage {
		return body, nil
	}

	return body, &savePatch{
		Data:     body,
		ImageURL: openGraphURL,
		Webhook:  "releases",
	}
}
