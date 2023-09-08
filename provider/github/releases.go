package github

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/jx11r/qnotifier/provider"
	"github.com/jx11r/qnotifier/utils"
	"github.com/tidwall/gjson"
)

var release string

func Releases() error {
	obj := provider.Notifier{
		API:     "https://api.github.com/repos/qtile/qtile/releases/latest",
		Webhook: utils.Webhook["releases"],
	}

	data, err := obj.Fetch()
	if err != nil {
		return err
	}

	current := gjson.GetBytes(data, "tag_name").String()
	if release == "" {
		release = current
		return nil
	}

	previous, _ := semver.NewVersion(release)
	latest, err := semver.NewVersion(current)
	if err != nil {
		return nil
	}

	if latest.GreaterThan(previous) {
		release = current
		obj.Payload = getRelease(string(data))
		return obj.Send()
	}

	return nil
}

func getRelease(data string) []byte {
	url := gjson.Get(data, "html_url").String()

	imageURL, exists := getImage(url)
	if exists {
		if !isImage(imageURL) {
			imageURL = fallback
		}
	}

	payload := fmt.Sprintf(`{
		"username": "GitHub",
		"embeds": [{
			"title": "New release published: %s",
			"url": "%s",
			"color": %d,
			"image": {"url": "%s"}
		}]
	}`,
		gjson.Get(data, "tag_name").String(),
		url,
		0x60c67b,
		imageURL,
	)

	return []byte(payload)
}
