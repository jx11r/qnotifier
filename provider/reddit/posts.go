package reddit

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/jx11r/qnotifier/provider"
	"github.com/jx11r/qnotifier/utils"
	"github.com/tidwall/gjson"
)

var ids = make([]string, 5)

func Posts() error {
	obj := provider.Notifier{
		API:     "https://www.reddit.com/r/qtile/new.json?limit=1",
		Webhook: utils.Webhook["reddit"],
	}

	raw, err := obj.Fetch()
	if err != nil {
		return err
	}

	data := gjson.GetBytes(raw, "data.children.0.data").String()
	id := gjson.Get(data, "id").String()
	if ids[0] == "" {
		ids[0] = id
		return nil
	}

	for _, v := range ids {
		if v == id {
			return nil
		}
	}

	obj.Payload = getPost(data)
	copy(ids[1:], ids[:4])
	ids[0] = id

	return obj.Send()
}

func getPost(data string) []byte {
	const prefix string = "https://reddit.com"
	footer := "Flair: unspecified"
	author := gjson.Get(data, "author").String()

	flair := gjson.Get(data, "link_flair_text").String()
	if flair != "" {
		footer = "Flair: " + flair
	}

	thumbnail := gjson.Get(data, "thumbnail").String()
	if !strings.Contains(thumbnail, "https") {
		thumbnail = ""
	}

	url := gjson.Get(data, "url").String()
	if strings.Contains(url, "i.redd.it") {
		thumbnail = url
	}

	payload := fmt.Sprintf(`{
		"username": "Reddit",
		"embeds": [{
			"title": "%s",
			"url": "%s",
			"color": %d,
			"footer": {"text": "%s"},
			"thumbnail": {"url": "%s"},
			"author": {
				"name": "%s",
				"icon_url": "%s",
				"url": "%s"
			}
		}]
	}`,
		gjson.Get(data, "title").String(),
		prefix+gjson.Get(data, "permalink").String(),
		0xff4400,
		footer,
		thumbnail,
		author,
		fmt.Sprintf(
			"https://www.redditstatic.com/avatars/defaults/v2/avatar_default_%d.png",
			rand.Intn(8),
		),
		prefix+"/user/"+author,
	)

	return []byte(payload)
}
