package main

import (
	"log"
	"sync"

	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"

	"github.com/Rukenshia/graw"
	"github.com/Rukenshia/graw/reddit"
	logrus "github.com/sirupsen/logrus"
)

type Agent struct {
	bot    reddit.Bot
	client soccerstreams.DBClient

	guard   map[string]*sync.Mutex
	polling []string
}

func NewAgent(bot reddit.Bot, client soccerstreams.DBClient) *Agent {
	return &Agent{bot, client, make(map[string]*sync.Mutex), make([]string, 0)}
}

func (s *Agent) Run() error {
	logger := log.New(logrus.StandardLogger().Out, "graw", log.LstdFlags)

	cfg := graw.Config{
		Subreddits:        []string{"soccerstreams"},
		SubredditComments: []string{"soccerstreams"},
		Logger:            logger,
	}

	_, wait, err := graw.Run(s, s.bot, cfg)
	if err != nil {
		return err
	}
	return wait()
}
