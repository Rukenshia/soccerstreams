package parser

import (
	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"
	"github.com/turnage/graw/reddit"
)

type CommentParser interface {
	Parse(string) []*soccerstreams.Soccerstream
}

type PostParser interface {
	Parse(*reddit.Post) *soccerstreams.Matchthread
}

func ParseComment(c *reddit.Comment) []*soccerstreams.Soccerstream {
	parsers := []CommentParser{
		&singleStreamParser{},
		&acestreamParser{},
	}

	var s []*soccerstreams.Soccerstream

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

func ParsePost(p *reddit.Post) *soccerstreams.Matchthread {
	var parser PostParser

	return parser.Parse(p)
}
