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
	latest, _ := semver.NewVersion(current)

	if latest.GreaterThan(previous) {
		release = current
		obj.Payload = getRelease(string(data))
		return obj.Send()
	}

	return nil
}

func getRelease(data string) []byte {
	payload := fmt.Sprintf(`{
		"username": "GitHub",
		"embeds": [{
			"title": "%s",
			"description": "%s",
			"url": "%s",
			"color": %d,
			"author": {
				"name": "%s",
				"icon_url": "%s",
				"url": "%s"
			}
		}]
	}`,
		gjson.Get(data, "tag_name").String(),
		gjson.Get(data, "body").String(),
		gjson.Get(data, "html_url").String(),
		0x202225,
		gjson.Get(data, "author.login").String(),
		gjson.Get(data, "author.avatar_url").String(),
		gjson.Get(data, "author.html_url").String(),
	)

	return []byte(payload)
}
