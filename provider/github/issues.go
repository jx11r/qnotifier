package github

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jx11r/qnotifier/provider"
	"github.com/jx11r/qnotifier/utils"
	"github.com/tidwall/gjson"
)

const fallback = "https://raw.githubusercontent.com/jx11r/src/!/img/github.png"

var issue int

func Issues() error {
	obj := provider.Notifier{
		API:     "https://api.github.com/repos/qtile/qtile/issues?per_page=1",
		Webhook: utils.Webhook["issues"],
	}

	data, err := obj.Fetch()
	if err != nil {
		return err
	}

	new := int(gjson.GetBytes(data, "0.number").Int())
	if issue == 0 {
		issue = new
		return nil
	}

	if new > issue {
		issue = new
		if gjson.GetBytes(data, "0.pull_request").Exists() {
			obj.Webhook = utils.Webhook["pulls"]
		}
		obj.Payload = getIssue(string(data))
		return obj.Send()
	}

	return nil
}

func getIssue(data string) []byte {
	number := gjson.Get(data, "0.number").String()
	title := "Issue opened: #" + number
	color := 0xeb6420
	url := gjson.Get(data, "0.html_url").String()

	if gjson.Get(data, "0.pull_request").Exists() {
		title = "Pull request opened: #" + number
		color = 0x7289da
	}

	imageURL, exists := getImage(url)
	if exists {
		if !isImage(imageURL) {
			imageURL = fallback
		}
	}

	payload := fmt.Sprintf(`{
		"username": "GitHub",
		"embeds": [{
			"title": "%s",
			"url": "%s",
			"color": %d,
			"image": {"url": "%s"},
			"author": {
				"name": "%s",
				"icon_url": "%s",
				"url": "%s"
			}
		}]
	}`,
		title,
		url,
		color,
		imageURL,
		gjson.Get(data, "0.user.login").String(),
		gjson.Get(data, "0.user.avatar_url").String(),
		gjson.Get(data, "0.user.html_url").String(),
	)

	return []byte(payload)
}

func getImage(url string) (string, bool) {
	resp, err := http.Get(url)
	if err != nil {
		return fallback, false
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fallback, false
	}

	content, exists := doc.Find("meta[property='og:image']").Attr("content")
	if !exists {
		return fallback, false
	}

	return content, true
}

func isImage(url string) bool {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	for i := 0; i < 3; i++ {
		resp, err := client.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		content := resp.Header.Get("Content-Type")
		if strings.HasPrefix(content, "image/") {
			return true
		}

		time.Sleep(time.Second * 30)
	}

	return false
}
