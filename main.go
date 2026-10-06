package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"qnotifier/provider/github"
	"qnotifier/provider/reddit"
)

func tasks() {
	if github.PendingPatches() {
		github.ApplyPatches()
	}
	github.Issues()
	github.Releases()
	reddit.Posts()
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	tasks()

	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down...")
			return
		case <-ticker.C:
			tasks()
		}
	}
}
