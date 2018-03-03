package main

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"cloud.google.com/go/datastore"

	"github.com/Rukenshia/soc-agent/soccerstream"
	raven "github.com/getsentry/raven-go"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)

	f, err := os.OpenFile("soc-sweeper.log", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)

	sentryb, err := ioutil.ReadFile("sentry")
	if err != nil {
		log.Fatal(err)
	}
	sentryb = sentryb[:len(sentryb)-1]

	log.Debugf("Using sentry DSN: %s", string(sentryb))

	if err := raven.SetDSN(string(sentryb)); err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Debugf("Starting sweeper process")

	ctx := context.Background()
	client, err := soccerstream.NewClient()
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Minute)

	for t := range ticker.C {
		log.Debugf("Sweeping entries")
		query := datastore.NewQuery("matchthread").Filter("ExpiresAt <", t)

		var threads []*soccerstream.Matchthread

		keys, err := client.GetAll(ctx, query, &threads)
		if err != nil {
			raven.CaptureError(err, nil)
			log.Errorf("Could not get matchthreads: %v", err)
		}

		for _, thread := range threads {
			log.WithField("post_id", thread.RedditID).
				WithField("team1", thread.Team1).
				WithField("team2", thread.Team2).
				WithField("kickoff", thread.Kickoff).Debugf("Matchthread %s has been selected for sweeping, it expired %ds ago", thread.DBKey(), int64(time.Since(thread.ExpiresAt).Seconds()))
		}

		if err := client.DeleteMulti(ctx, keys); err != nil {
			log.Errorf("Could not delete matchthreads: %v", err)
		}
	}
}
