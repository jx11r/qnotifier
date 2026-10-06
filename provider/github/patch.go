package github

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"qnotifier/discord"
	"qnotifier/provider"

	"github.com/tidwall/sjson"
)

type savePatch struct {
	Data     []byte
	ImageURL string
	Webhook  string
}

var toPatch = map[string]*savePatch{}

func PendingPatches() bool {
	return len(toPatch) > 0
}

func ApplyPatches() {
	for id, patch := range toPatch {
		if !isImage(patch.ImageURL) {
			log.Printf("error (patch:%s): image could not be loaded, skipping", id)
			continue
		}

		payload, err := sjson.SetBytes(patch.Data, "embeds.0.image.url", patch.ImageURL)
		if err != nil {
			log.Printf("error (patch:%s): %v", id, err)
			continue
		}

		err = patchWebhook(discord.Webhook[patch.Webhook], id, payload)
		if err != nil {
			log.Printf("error (patch:%s): %v", id, err)
			continue
		}

		delete(toPatch, id)
		log.Printf("patch: %s was successfully updated", id)
		time.Sleep(3 * time.Second)
	}
}

func patchWebhook(webhook, id string, payload []byte) error {
	url := fmt.Sprintf("%s/messages/%s", webhook, id)
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := provider.Client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("failed to patch webhook, discord returned %s", resp.Status)
}
