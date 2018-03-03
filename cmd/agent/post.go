package main

import (
	"github.com/Rukenshia/soccerstreams/pkg/monitoring"
	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams/parser"
	raven "github.com/getsentry/raven-go"
	log "github.com/sirupsen/logrus"
	"github.com/turnage/graw/reddit"
)

func (s *SOCAgent) Post(p *reddit.Post) error {
	mt := parser.ParsePost(p)

	logger := log.WithField("post_id", p.ID).
		WithField("author", p.Author).
		WithField("title", p.Title)

	if mt != nil {
		mt.SetClient(s.client)
		mt.FillRedditInfo(p)

		logger = logger.WithField("team1", mt.Team1).
			WithField("team2", mt.Team2).
			WithField("kickoff", mt.Kickoff.Format("15:04"))

		logger.Debugf("Parsed matchthread")
		logger.Debugf("Saving to datastore")

		if err := mt.Save(); err != nil {
			raven.CaptureError(err, map[string]string{
				"PostID": p.ID,
			})
			log.Errorf("Could not save to datastore: %v", err)
		}

		logger.Debugf("Saved to datastore")
		return nil
	}

	logger.Debugf("could not parse post")

	raven.Capture(
		logging.CreatePacket(raven.DEBUG, "Unable to parse post\nPermalink: https://reddit.com%s\nTitle: %s\nFlair: %s", p.Permalink, p.Title, p.LinkFlairText),
		map[string]string{
			"Author": p.Author,
			"PostID": p.ID,
		})
	return nil
}
