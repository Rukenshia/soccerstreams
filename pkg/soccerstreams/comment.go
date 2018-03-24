package soccerstreams

import (
	"crypto/md5"
	"fmt"
)

// A Comment is a Reply to a Matchthread that contains one or more Streams. It holds Metadata to the Streamer like them being SOTW
// and the amount of Upvotes the comment received.
type Comment struct {
	Streams          []*Stream
	RedditID         string
	Permalink        string
	Streamer         string
	ReliableStreamer bool
	Upvotes          int32

	Body string

	// BodyHash is the Comment Hash (md5) from when the Body was parsed
	BodyHash string
}

// UpdateHash calculates the md5 hash of the Comment body.
func (c *Comment) UpdateHash() {
	c.BodyHash = fmt.Sprintf("%x", md5.Sum([]byte(c.Body)))
}
