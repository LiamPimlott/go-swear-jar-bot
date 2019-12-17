package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/turnage/graw/reddit"
)

type swearBot struct {
	bot   reddit.Bot
	db    *bolt.DB
	regex *regexp.Regexp
}

func (sb *swearBot) Post(p *reddit.Post) error {
	bodySwears := sb.regex.FindAllString(p.SelfText, -1)
	titleSwears := sb.regex.FindAllString(p.Title, -1)
	log.Printf("*** NEW POST *** id: %s  author: %s title_swears: %+v body_swears: %+v created_at: %d\n", p.Name, p.Author, titleSwears, bodySwears, p.CreatedUTC)

	total, err := sb.updateSwearCount(len(bodySwears)+len(titleSwears), p.Author)
	if err != nil {
		return fmt.Errorf("Failed to update swear count after post (%s): %s", err, p.Name)
	}
	log.Printf("*** USER BUCKET UPDATED *** key: %s  value: %s post: %s\n", p.Author, total, p.Name)

	return nil
}

func (sb *swearBot) Comment(c *reddit.Comment) error {
	bodySwears := sb.regex.FindAllString(c.Body, -1)
	fmt.Printf("*** NEW COMMENT *** id: %s author: %s total_swears: %v created_at: %d parent_id: %s\n", c.Name, c.Author, len(bodySwears), c.CreatedUTC, c.ParentID)

	total, err := sb.updateSwearCount(len(bodySwears), c.Author)
	if err != nil {
		return fmt.Errorf("Failed to update swear count after comment (%s): %s", err, c.Name)
	}
	log.Printf("*** USER BUCKET UPDATED *** key: %s  value: %s post: %s\n", c.Author, total, c.Name)

	return nil
}

func (sb *swearBot) updateSwearCount(newSwears int, user string) (string, error) {
	var newTotal []byte

	err := sb.db.Update(func(tx *bolt.Tx) error {
		var err error
		var currTotal int

		bucket := tx.Bucket([]byte(usersKey))
		key := []byte(user)
		currVal := string(bucket.Get(key))

		if currVal == "" {
			currTotal = 0
		} else {
			currTotal, err = strconv.Atoi(currVal)
			if err != nil {
				return fmt.Errorf("Failed convert string \"%s\"to int: %s", err, currVal)
			}
		}

		newTotal = []byte(strconv.Itoa(currTotal + newSwears))

		err = bucket.Put([]byte(key), newTotal)
		if err != nil {
			return fmt.Errorf("Failed to put value (%s) to bucket (%s): %s", newTotal, usersKey, err)
		}

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to execute batch: %s", err)
	}

	return string(newTotal), nil
}
