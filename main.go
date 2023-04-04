package main

import (
	"fmt"
	"time"

	"github.com/jx11r/qnotifier/provider/github"
	"github.com/jx11r/qnotifier/provider/reddit"
)

func task(id string, fn func() error) {
	for {
		if err := fn(); err != nil {
			fmt.Printf("%s: %s\n", id, err.Error())
			time.Sleep(time.Minute)
		}
		time.Sleep(time.Minute)
	}
}

func main() {
	go task("issues", github.Issues)
	go task("releases", github.Releases)
	task("reddit", reddit.Posts)
}
