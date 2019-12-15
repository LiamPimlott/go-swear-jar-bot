package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/turnage/graw"
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

	swearRegex = os.Getenv("SWEARBOT_AGENT_FILE")
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
	if bot, err := reddit.NewBot(swearBotConfig); err != nil {
		fmt.Println("Failed to create bot handle: ", err)
	} else {
		cfg := graw.Config{
			Subreddits:        []string{subreddit},
			SubredditComments: []string{subreddit},
		}
		handler := &swearBot{bot: bot}
		if _, wait, err := graw.Run(handler, bot, cfg); err != nil {
			fmt.Println("Failed to start graw run: ", err)
		} else {
			fmt.Println("graw run failed: ", wait())
		}
	}
}
