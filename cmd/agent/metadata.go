package main

import (
	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"
	"github.com/turnage/graw/reddit"
)

// FillMatchthreadInfo adds metadata to a matchthread from a reddit post
func FillMatchthreadInfo(m *soccerstreams.Matchthread, p *reddit.Post) {
	m.RedditID = p.ID
	m.Competition = p.LinkFlairCSSClass
	m.CompetitionName = p.LinkFlairText
}

// FillCommentInfo adds metadata information to the comment
func FillCommentInfo(c *soccerstreams.Comment, rc *reddit.Comment) {
	c.Body = rc.Body
	c.Streamer = rc.Author
	c.Permalink = rc.Permalink
	c.RedditID = rc.ID

	if rc.AuthorFlairCSSClass == "weekly" {
		c.ReliableStreamer = true
	}
	c.Upvotes = rc.Ups
}
