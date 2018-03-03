package parser

import (
	"github.com/Rukenshia/soc-agent/soccerstream"
	"github.com/turnage/graw/reddit"
)

type CommentParser interface {
	Parse(string) []*soccerstream.Soccerstream
}

type PostParser interface {
	Parse(*reddit.Post) *soccerstream.Matchthread
}

func ParseComment(c *reddit.Comment) []*soccerstream.Soccerstream {
	parsers := []CommentParser{
		&singleStreamParser{},
		&acestreamParser{},
	}

	var s []*soccerstream.Soccerstream

	for _, parser := range parsers {
		s = append(s, parser.Parse(c.Body)...)
	}

	// fill comment info
	for _, stream := range s {
		stream.CommentLink = c.Permalink
		stream.Streamer = c.Author
	}

	return s
}

func ParsePost(p *reddit.Post) *soccerstream.Matchthread {
	var parser PostParser

	return parser.ParseThread(p)
}
