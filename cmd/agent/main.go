package main

import (
	"io/ioutil"

	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"
	raven "github.com/getsentry/raven-go"
	log "github.com/sirupsen/logrus"
	"github.com/turnage/graw/reddit"
)

func init() {
	log.SetLevel(log.DebugLevel)

	sentryb, err := ioutil.ReadFile("sentry")
	if err != nil {
		log.Fatal(err)
	}
	sentryb = sentryb[:len(sentryb)-1]

	log.Debugf("Using sentry DSN: %s", string(sentryb))

	raven.SetDSN(string(sentryb))
}

func main() {
	client, err := soccerstreams.NewDatastoreClient()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := reddit.NewBotFromAgentFile("graw", 0)
	if err != nil {
		log.Fatal(err)
	}

	agent := NewSOCAgent(bot, client)

	if err := agent.Run(); err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.Fatal(err)
	}
}
