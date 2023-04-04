package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jx11r/qnotifier/provider/github"
	"github.com/jx11r/qnotifier/provider/reddit"
)

func task(id string, fn func() error) {
	for {
		if err := fn(); err != nil {
			fmt.Printf("%s: %s\n", id, err.Error())
			time.Sleep(time.Minute * 2)
			continue
		}
		time.Sleep(time.Minute)
	}
}

func main() {
	go task("issues", github.Issues)
	go task("releases", github.Releases)
	go task("reddit", reddit.Posts)

	fmt.Println("Press Enter to exit...")
	if _, err := fmt.Scanln(); err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}
}
