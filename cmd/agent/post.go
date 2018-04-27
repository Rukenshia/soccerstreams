package main

import (
	"sync"
	"time"

	"github.com/Rukenshia/soccerstreams/cmd/agent/metrics"

	"github.com/Rukenshia/graw/reddit"
	"github.com/Rukenshia/soccerstreams/pkg/monitoring"
	"github.com/Rukenshia/soccerstreams/pkg/parser"
	raven "github.com/getsentry/raven-go"
	log "github.com/sirupsen/logrus"
)

func (s *Agent) Post(p *reddit.Post) error {
	metrics.GrawEvents.WithLabelValues("stream_post").Inc()

	timeCreated := time.Unix(int64(p.CreatedUTC), 0)
	metrics.GrawEventDiff.Observe(float64(time.Now().Sub(timeCreated) / time.Second))

	mt := parser.ParsePost(p)

	logger := log.WithField("post_id", p.ID).
		WithField("author", p.Author).
		WithField("title", p.Title).
		WithField("component", "agent")
	logger.Debugf("handling post")

	if _, ok := s.guard[p.ID]; !ok {
		s.guard[p.ID] = &sync.Mutex{}
	}

	s.guard[p.ID].Lock()

	defer func() {
		s.guard[p.ID].Unlock()
	}()

	if mt != nil {
		mt.SetClient(s.client)
		FillMatchthreadInfo(mt, p)
		mt.UpdateExpiresAt()

		logger = logger.WithField("team1", mt.Team1).
			WithField("team2", mt.Team2).
			WithField("kickoff", mt.Kickoff.Format("15:04"))

		metrics.PostsIngested.WithLabelValues("successful").Inc()

		if err := mt.Save(); err != nil {
			raven.CaptureError(err, map[string]string{
				"PostID": p.ID,
			})
			log.Errorf("Could not save to datastore: %v", err)
		}

		logger.Debugf("Saved to datastore")

		logger.Debugf("Start polling")
		s.StartPolling(mt)
		return nil
	}

	metrics.PostsIngested.WithLabelValues("failed").Inc()
	logger.Debugf("could not parse post")

	raven.Capture(
		logging.CreatePacket(raven.DEBUG, "Unable to parse post\nPermalink: https://reddit.com%s\nTitle: %s\nFlair: %s", p.Permalink, p.Title, p.LinkFlairText),
		map[string]string{
			"Author": p.Author,
			"PostID": p.ID,
		})
	return nil
}
