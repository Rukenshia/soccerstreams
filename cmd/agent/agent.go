package main

import (
	"log"
	"os"

	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"

	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

type Agent struct {
	bot    reddit.Bot
	client soccerstreams.DBClient
}

func NewAgent(bot reddit.Bot, client soccerstreams.DBClient) *Agent {
	return &Agent{bot, client}
}

func (s *Agent) Run() error {
	cfg := graw.Config{Subreddits: []string{"soccerstreams"}, SubredditComments: []string{"soccerstreams"}, Logger: log.New(os.Stdout, "graw", 0)}

	_, wait, err := graw.Run(s, s.bot, cfg)
	if err != nil {
		return err
	}
	return wait()
}
