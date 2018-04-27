package main

import (
	"context"
	"os"
	"time"

	"github.com/Rukenshia/graw/reddit"
	"github.com/Rukenshia/soccerstreams/cmd/agent/metrics"
	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"

	raven "github.com/getsentry/raven-go"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
	raven.SetDSN(os.Getenv("SENTRY_DSN"))
}

func main() {
	client, err := soccerstreams.NewDatastoreClient(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	bot, err := reddit.NewBotFromAgentFile("/opt/soccerstreams/graw/graw", time.Second*15)
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
