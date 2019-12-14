package main

import (
	"fmt"
	"os"
	"regexp"
	// "strings"
	// "time"

	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

const (
	envSubreddit  = "SWEARBOT_SUBREDDIT"
	envAgentFile  = "SWEARBOT_AGENT_FILE"
	envSwearRegex = "SWEARBOT_REGEX"
)

var (
	subreddit  string
	agentFile  string
	swearRegex string
	re         *regexp.Regexp
)

func init() {
	subreddit = os.Getenv(envSubreddit)
	agentFile = os.Getenv(envAgentFile)
	swearRegex = os.Getenv(envSwearRegex)
	re = regexp.MustCompile(swearRegex)
}

func main() {
	if bot, err := reddit.NewBotFromAgentFile(agentFile, 0); err != nil {
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
