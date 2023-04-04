package utils

import "os"

var Token = map[string]string{
	"github": os.Getenv("GITHUB_TOKEN"),
}

var Webhook = map[string]string{
	"issues":   os.Getenv("WEBHOOK_ISSUES"),
	"pulls":    os.Getenv("WEBHOOK_PULLS"),
	"reddit":   os.Getenv("WEBHOOK_REDDIT"),
	"releases": os.Getenv("WEBHOOK_RELEASES"),
}
