package soccerstreams

import (
	"crypto/md5"
	"fmt"

	"github.com/turnage/graw/reddit"
)

// A Stream is a link to media posted inside of a Matchthread. This struct contains all information that can be parsed from reddit.
// This mainly relies on the rules of /r/soccerstreams.
type Stream struct {
	CommentLink    string
	Streamer       string
	Quality        string
	Channel        string
	Link           string
	Clicks         int
	IsNSFW         bool
	MISR           string
	MobileFriendly bool

	Metadata struct {
		Hash             string
		ReliableStreamer bool
		Upvotes          int32
	}
}

// IsGood returns whether the Stream has sufficient information to be used further
func (s *Stream) IsGood() bool {
	if s.Link == "" || s.Channel == "" {
		return false
	}

	return true
}

// FillMetadata adds metadata information to the stream
func (s *Stream) FillMetadata(c *reddit.Comment) {
	if c.AuthorFlairCSSClass == "weekly" {
		s.Metadata.ReliableStreamer = true
	}
	s.Metadata.Hash = fmt.Sprintf("%x", md5.Sum([]byte(c.Body)))
	s.Metadata.Upvotes = c.Ups
}
