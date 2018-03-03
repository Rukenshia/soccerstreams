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
	client, err := soccerstream.NewDatastoreClient()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := reddit.NewBotFromAgentFile("graw", 0)
	if err != nil {
		log.Fatal(err)
	}

	agent := NewSOCAgent(bot, client)

	// agent.Comment(&reddit.Comment{
	// 	ParentID: "t3_7yxtok",
	// 	Author:   "me",
	// 	Body:     "**HD** | [ENGLISH TSN4 1080p] (http://foundationsports.com/crvcsk/) | MISR : 3mbps | Ad Overlay : 1 | Clicks : 2  | Mobile : Yes.",
	// })
	err = agent.Run()

	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.Fatal(err)
	}
}
