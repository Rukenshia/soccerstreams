package main

import (
	"context"
	"os"

	"github.com/Rukenshia/soccerstreams/cmd/agent/metrics"
	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"
	raven "github.com/getsentry/raven-go"
	log "github.com/sirupsen/logrus"
	"github.com/turnage/graw/reddit"
)

func init() {
	log.SetLevel(log.DebugLevel)
	raven.SetDSN(os.Getenv("SENTRY_DSN"))
}

func main() {
	client, err := soccerstreams.NewDatastoreClient(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	bot, err := reddit.NewBotFromAgentFile("/opt/soccerstreams/graw/graw", 0)
	if err != nil {
		log.Fatal(err)
	}

	agent := NewAgent(bot, client)

	go func() {
		metrics.Register()
		log.Error(metrics.Serve())
	}()

	if err := agent.Run(); err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.Fatal(err)
	}
}
