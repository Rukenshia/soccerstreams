package main

import (
	"context"
	"os"
	"time"

	"cloud.google.com/go/datastore"

	"github.com/Rukenshia/soccerstreams/cmd/sweeper/metrics"
	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"
	raven "github.com/getsentry/raven-go"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	raven.SetDSN(os.Getenv("SENTRY_DSN"))
}

func main() {
	log.Debugf("Starting sweeper process")

	client, err := soccerstreams.NewDatastoreClient(context.Background())
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.Fatal(err)
	}

	go func() {
		metrics.Register()
		log.Error(metrics.Serve())
	}()

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

		metrics.MatchthreadsDeleted.Add(float64(len(threads)))
	}
}
