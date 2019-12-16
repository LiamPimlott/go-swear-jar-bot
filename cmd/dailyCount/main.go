package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/turnage/graw/reddit"
)

var (
	subreddit      string
	swearRegex     string
	swearBotConfig reddit.BotConfig
	re             *regexp.Regexp
)

func init() {
	subreddit = os.Getenv("SWEARBOT_SUBREDDIT")

	swearRegex = os.Getenv("SWEARBOT_REGEX")
	re = regexp.MustCompile(swearRegex)

	swearBotConfig = reddit.BotConfig{
		Agent: os.Getenv("SWEARBOT_USER_AGENT"),
		App: reddit.App{
			ID:       os.Getenv("SWEARBOT_APP_ID"),
			Secret:   os.Getenv("SWEARBOT_APP_SECRET"),
			Username: os.Getenv("SWEARBOT_APP_USERNAME"),
			Password: os.Getenv("SWEARBOT_APP_PASSWORD"),
		},
	}
}

func main() {
	start := time.Now()

	if bot, err := reddit.NewBot(swearBotConfig); err != nil {
		fmt.Println("Failed to create bot handle: ", err)
	} else {
		handler := &swearBot{bot: bot}
		handler.DailyCount()
	}

	fmt.Printf("\n\nExecution time: %s", time.Since(start))
}
