package agent

import (
	"log"
	"os"

	"github.com/Rukenshia/soc-agent/soccerstream"

	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

type SOCAgent struct {
	bot    reddit.Bot
	client soccerstream.DBClient
}

func NewSOCAgent(bot reddit.Bot, client soccerstream.DBClient) *SOCAgent {
	return &SOCAgent{bot, client}
}

func (s *SOCAgent) Run() error {
	cfg := graw.Config{Subreddits: []string{"soccerstreams"}, SubredditComments: []string{"soccerstreams"}, Logger: log.New(os.Stdout, "graw", 0)}

	_, wait, err := graw.Run(s, s.bot, cfg)
	if err != nil {
		return err
	}
	return wait()
}
