package discord

import "os"

var Webhook = map[string]string{
	"issues":      os.Getenv("WEBHOOK_ISSUES"),
	"discussions": os.Getenv("WEBHOOK_DISCUSSIONS"),
	"pulls":       os.Getenv("WEBHOOK_PULLS"),
	"releases":    os.Getenv("WEBHOOK_RELEASES"),
}
