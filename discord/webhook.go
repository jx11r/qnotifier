package discord

import "os"

var Webhook = map[string]string{
	"issues":      os.Getenv("WEBHOOK_ISSUES"),
	"discussions": os.Getenv("WEBHOOK_DISCUSSIONS"),
	"pulls":       os.Getenv("WEBHOOK_PULLS"),
	"reddit":      os.Getenv("WEBHOOK_REDDIT"),
	"releases":    os.Getenv("WEBHOOK_RELEASES"),
}
