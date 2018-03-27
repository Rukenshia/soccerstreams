package main

import (
	"crypto/md5"
	"fmt"

	"github.com/Rukenshia/soccerstreams/cmd/agent/metrics"

	"github.com/Rukenshia/soccerstreams/pkg/parser"
	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"
	raven "github.com/getsentry/raven-go"
	log "github.com/sirupsen/logrus"
	"github.com/turnage/graw/reddit"
)

// StartPolling starts polling a Matchthread for updates. If we are already polling a Matchthread, no action will be taken.
func (s *Agent) StartPolling(mt *soccerstreams.Matchthread) {

	logger := log.WithField("post_id", mt.RedditID).
		WithField("polling", true)

	for _, p := range s.polling {
		if p == mt.RedditID {
			// We are already polling this thread; no further action required
			return
		}
	}

	metrics.PostsPolling.Inc()

	go func() {
		defer func() {
			// Delete entry from polling array
			for idx, p := range s.polling {
				if p == mt.RedditID {
					s.polling = append(s.polling[:idx], s.polling[idx+1:]...)
					break
				}
			}

			metrics.PostsPolling.Dec()
		}()

		poller := NewThreadPoller(fmt.Sprintf("/r/soccerstreams/comments/%s", mt.RedditID), *mt.Kickoff, s.bot)

		s.polling = append(s.polling, mt.RedditID)

		pollChannel := poller.Poll()
		for {
			update, ok := <-pollChannel
			if !ok {
				logger.Debugf("Poll thread reached EOL")
				break
			}

			keepOpen, err := s.HandleUpdate(mt.RedditID, update)
			if err != nil {
				logger.Warnf("Error while handling update: %v", err)
				if !keepOpen {
					logger.Warnf("-> Stop handling updates")
					close(pollChannel)
					break
				}
			}

			if !keepOpen {
				logger.Debugf("Stopping update polling, keepOpen is false. No error")
				close(pollChannel)
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

	if post.Deleted || post.Hidden || post.SelfText == "[removed]" {
		metrics.PostsDeleted.Inc()

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

	updated := false
	var comments []*soccerstreams.Comment

	for _, c := range post.Replies {
		streamLogger := logger.WithField("comment_id", c.ID).
			WithField("author", c.Author)

		// Find comment in matchthread
		var existing *soccerstreams.Comment
		for _, ec := range mt.Comments {
			if ec.RedditID == c.ID {
				existing = ec
			}
		}

		if existing == nil {
			// We don't know this comment so let's not bother
			continue
		}

		if c.Deleted || c.Body == "[removed]" {
			metrics.CommentsDeleted.Inc()

			streamLogger.Debugf("Comment was deleted. Removing streams")
			continue
		}

		// Check comment for changes
		hash := fmt.Sprintf("%x", md5.Sum([]byte(c.Body)))
		if hash == existing.BodyHash {
			comments = append(comments, existing)

			// Check if Upvotes changed to mark the Matchthread as updated
			if existing.Upvotes != c.Ups {
				FillCommentInfo(existing, c)
				updated = true
			}
			continue
		}

		updated = true

		streamLogger.WithFields(log.Fields{
			"old": existing.BodyHash,
			"new": hash,
		}).Debugf("Comment hash changed, updating all streams")

		existing.Streams = parser.ParseComment(c.Body)
		FillCommentInfo(existing, c)
		existing.Body = c.Body
		existing.UpdateHash()

		comments = append(comments, existing)
	}

	if updated {
		metrics.CommentsChanged.Inc()

		logger.WithFields(log.Fields{
			"comments_before": len(mt.Comments),
			"comments_after":  len(comments),
		}).Debugf("Updating matchthread after polling update")

		mt.Comments = comments
		if err := mt.Save(); err != nil {
			logger.Errorf("Could not save matchthread: %v", err)
			return true, err
		}
	}

	return true, nil
}
