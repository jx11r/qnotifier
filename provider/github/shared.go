package github

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"qnotifier/provider"
)

const fallbackImage = "https://raw.githubusercontent.com/jx11r/src/assets/qnotifier/not_found.png"

const userAgent = "qnotifier/1.0"

var token = map[string]string{
	"github": os.Getenv("GITHUB_TOKEN"),
}

var header = http.Header{
	"User-Agent":           {userAgent},
	"Accept":               {"application/vnd.github+json"},
	"Authorization":        {"token " + token["github"]},
	"X-GitHub-Api-Version": {"2026-03-10"},
}

func getOpenGraphURL(updatedAt int64, itemType, id string) string {
	return fmt.Sprintf("https://opengraph.githubassets.com/%d/qtile/qtile/%s/%s",
		updatedAt,
		itemType,
		id,
	)
}

func isImage(url string) bool {
	for i := range 3 {
		req, err := http.NewRequest(http.MethodHead, url, nil)
		if err != nil {
			log.Printf("error (github:image): %v", err)
			return false
		}

		req.Header.Set("User-Agent", userAgent)

		resp, err := provider.Client.Do(req)
		if err == nil {
			contentType := resp.Header.Get("Content-Type")
			_ = resp.Body.Close()

			if resp.StatusCode == http.StatusOK && strings.HasPrefix(contentType, "image/") {
				return true
			}

			log.Printf("error (github:image): %s returned %s", url, resp.Status)
		}

		if i < 2 {
			time.Sleep(5 * time.Second)
		}
	}

	return false
}
