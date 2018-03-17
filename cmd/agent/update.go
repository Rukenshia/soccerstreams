package main

import (
	"fmt"

	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"
	raven "github.com/getsentry/raven-go"
	log "github.com/sirupsen/logrus"
	"github.com/turnage/graw/reddit"
)

// StartPolling starts polling a Matchthread for updates. If we are already polling a Matchthread, no action will be taken.
func (s *Agent) StartPolling(mt *soccerstreams.Matchthread) {
	for _, p := range s.polling {
		if p == mt.RedditID {
			// We are already polling this thread; no further action required
			return
		}
	}

	logger := log.WithField("post_id", mt.RedditID).
		WithField("polling", true)

	go func() {
		defer func() {
			// Delete entry from polling array
			for idx, p := range s.polling {
				if p == mt.RedditID {
					s.polling = append(s.polling[:idx], s.polling[idx+1:]...)
					break
				}
			}

			logger.Debugf("Stopped polling")
		}()

		poller := NewThreadPoller(fmt.Sprintf("https://reddit.com/r/soccerstreams/comments/%s", mt.RedditID), *mt.Kickoff, s.bot)

		s.polling = append(s.polling, mt.RedditID)

		for update := range poller.Poll() {
			if update == nil {
				break
			}

			keepOpen, err := s.HandleUpdate(mt.RedditID, update)
			if err != nil {
				logger.Debugf("Stopping update polling, error: %v", err)
				break
			}

			if !keepOpen {
				logger.Debugf("Stopping update polling, keepOpen is false. No error")
				break
			}
		}
	}()
}

// HandleUpdate handles polling updates for Matchthreads
func (s *Agent) HandleUpdate(postID string, post *reddit.Post) (bool, error) {
	// Get the Matchthread from the database, otherwise fail
	mt, err := s.client.Get(postID)
	if err != nil {
		return false, err
	}

	logger := log.WithField("post_id", postID).
		WithField("polling", true)

	logger.Debugf("received update")

	if post.Deleted || post.Hidden {
		logger.Debugf("Thread has been deleted or hidden, removing from database")
		if err := mt.Delete(); err != nil {
			logger.Debugf("Could not delete matchthread: %v", err)

			raven.CaptureError(err, map[string]string{
				"PostID":  postID,
				"Polling": "yes",
			})
			return false, nil
		}
	}

	return true, nil
}
