package parser

import (
	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"
	"github.com/turnage/graw/reddit"
)

type CommentParser interface {
	Parse(string) []*soccerstreams.Stream
}

type PostParser interface {
	Parse(*reddit.Post) *soccerstreams.Matchthread
}

func ParseComment(comment string) []*soccerstreams.Stream {
	parsers := []CommentParser{
		&singleStreamParser{},
		&acestreamParser{},
	}

	var s []*soccerstreams.Stream

	for _, parser := range parsers {
		s = append(s, parser.Parse(comment)...)
	}

	return s
}

func ParsePost(p *reddit.Post) *soccerstreams.Matchthread {
	var parser threadParser

	return parser.Parse(p)
}
