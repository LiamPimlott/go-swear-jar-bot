package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/boltdb/bolt"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

const usersKey = "users"

var (
	subreddit      string
	re             *regexp.Regexp
	swearBotConfig reddit.BotConfig
	grawConfig     graw.Config
)

func init() {
	subreddit = os.Getenv("SWEARBOT_SUBREDDIT")
	re = regexp.MustCompile(os.Getenv("SWEARBOT_REGEX"))

	swearBotConfig = reddit.BotConfig{
		Agent: os.Getenv("SWEARBOT_USER_AGENT"),
		App: reddit.App{
			ID:       os.Getenv("SWEARBOT_APP_ID"),
			Secret:   os.Getenv("SWEARBOT_APP_SECRET"),
			Username: os.Getenv("SWEARBOT_APP_USERNAME"),
			Password: os.Getenv("SWEARBOT_APP_PASSWORD"),
		},
	}

	grawConfig = graw.Config{
		Subreddits:        []string{subreddit},
		SubredditComments: []string{subreddit},
	}
}

func main() {
	// Open the swearbot.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("swearbot.db", 0600, nil)
	if err != nil {
		log.Fatalf("Failed to open db: %s", err)
	}
	defer db.Close()

	// Init swearbot.db buckets if they dont exist
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(usersKey))
		if err != nil {
			return fmt.Errorf("Failed to create bucket \"%s\": %s", usersKey, err)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to initialize db: %s", err)
	}

	bot, err := reddit.NewBot(swearBotConfig)
	if err != nil {
		log.Fatalf("Failed to create bot handle: %s", err)
	}

	handler := &swearBot{bot: bot, db: db, regex: re}
	_, wait, err := graw.Run(handler, bot, grawConfig)
	if err != nil {
		log.Fatalf("Failed to start graw run: %s", err)
	}

	log.Fatal("graw run failed: ", wait())
}
