package main

import (
	"fmt"

	"github.com/turnage/graw/reddit"
)

type swearBot struct {
	bot reddit.Bot
}

func (r *swearBot) Post(p *reddit.Post) error {
	bodySwears := re.FindAllString(p.SelfText, -1)
	titleSwears := re.FindAllString(p.Title, -1)
	fmt.Printf("*** NEW POST *** id: %s  author: %s title_swears: %+v body_swears: %+v created_at: %d\n", p.Name, p.Author, titleSwears, bodySwears, p.CreatedUTC)
	return nil
}

func (r *swearBot) Comment(p *reddit.Comment) error {
	bodySwears := re.FindAllString(p.Body, -1)
	fmt.Printf("*** NEW COMMENT *** id: %s author: %s swears: %+v created_at: %d parent_id: %s\n", p.Name, p.Author, bodySwears, p.CreatedUTC, p.ParentID)
	return nil
}
