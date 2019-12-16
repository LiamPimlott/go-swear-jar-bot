package main

import (
	"fmt"

	"github.com/turnage/graw/reddit"
)

var indent string

type swearBot struct {
	bot reddit.Bot
}

func (r *swearBot) DailyCount() error {
	// var posts []*reddit.Post
	var allComments []*reddit.Comment

	// hardcoded "after" (which is backwards i think its supposed to be before)
	// should store last 10 post fullnames in case some are deleted
	harvest, err := r.bot.Listing(subreddit, "t3_eavc0d")
	if err != nil {
		fmt.Printf("Failed to fetch %s: %s\n", subreddit, err)
		return err
	}

	for _, post := range harvest.Posts {
		var postComments []*reddit.Comment

		fmt.Printf("\n\n### POST ###\n\n")
		fmt.Printf("[%s] posted [%s]\n", post.Author, post.Title)

		p, err := r.bot.Thread(post.Permalink)
		if err != nil {
			fmt.Printf("Failed to fetch thread %s: %s\n", post.Permalink, err)
		}

		fmt.Printf("\n### COMMENTS ###\n\n")
		for _, topLevelComment := range p.Replies {
			indent = ""
			if err := appendCommentsDepthFirst(topLevelComment, &postComments); err != nil {
				fmt.Printf("Failed to get comments depth first for top level comment (%s): %s\n", topLevelComment.Name, err)
			}
		}

		fmt.Printf("\nTotal Post Comments: %v\n\n", len(postComments))
		allComments = append(allComments, postComments...)
	}

	fmt.Printf("\nTotal Comments found: %v\n", len(allComments))
	return nil
}

// printed tree matches thread but len of postComments is off by a bit in PGT

func appendCommentsDepthFirst(c *reddit.Comment, pc *[]*reddit.Comment) error {
	fmt.Printf("%s[%s] commented [%s]\n", indent, c.Author, c.Name)
	indent = fmt.Sprintf("%s---", indent)

	if len(c.Replies) > 0 {
		for _, sc := range c.Replies {
			if err := appendCommentsDepthFirst(sc, pc); err != nil {
				return err
			}

			*pc = append(*pc, sc)
		}
	} else {
		*pc = append(*pc, c)
	}

	return nil
}
