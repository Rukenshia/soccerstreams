package main

import (
	"context"
	"io/ioutil"
	"time"

	"cloud.google.com/go/datastore"

	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"
	raven "github.com/getsentry/raven-go"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)

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

	client, err := soccerstreams.NewDatastoreClient(context.Background())
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Minute)

	for t := range ticker.C {
		log.Debugf("Sweeping entries")
		query := datastore.NewQuery("matchthread").Filter("ExpiresAt <", t)

		threads, err := client.GetAll(query)
		if err != nil {
			raven.CaptureError(err, nil)
			log.Errorf("Could not get matchthreads: %v", err)
		}

		var ids []string
		for _, thread := range threads {
			log.WithField("post_id", thread.RedditID).
				WithField("team1", thread.Team1).
				WithField("team2", thread.Team2).
				WithField("kickoff", thread.Kickoff).Debugf("Matchthread %s has been selected for sweeping, it expired %ds ago", thread.DBKey(), int64(time.Since(thread.ExpiresAt).Seconds()))

			ids = append(ids, thread.DBKey())
		}

		if err := client.DeleteMulti(ids); err != nil {
			log.Errorf("Could not delete matchthreads: %v", err)
		}
	}
}
