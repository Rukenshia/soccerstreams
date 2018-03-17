package main

import (
	"crypto/md5"
	"fmt"

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

		poller := NewThreadPoller(fmt.Sprintf("/r/soccerstreams/comments/%s", mt.RedditID), *mt.Kickoff, s.bot)

		s.polling = append(s.polling, mt.RedditID)

		pollChannel := poller.Poll()
		for {
			update := <-pollChannel
			if update == nil {
				logger.Warnf("Stopping update polling, received nil update")
				break
			}

			keepOpen, err := s.HandleUpdate(mt.RedditID, update)
			if err != nil {
				logger.Warnf("Error while handling update: %v", err)
				if !keepOpen {
					logger.Warnf("-> stop handling updates")
					break
				}
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

	var updated = false
	var streams []*soccerstreams.Stream

	for _, c := range post.Replies {
		streamLogger := logger.WithField("comment_id", c.ID).
			WithField("author", c.Author)

		if c.Deleted {
			streamLogger.Debugf("Comment was deleted. Removing streams")
			continue
		}

		// Find all streams for this comment
		var commentStreams []*soccerstreams.Stream

		for _, stream := range mt.Streams {
			if stream.CommentLink == c.Permalink {
				commentStreams = append(commentStreams, stream)
			}
		}

		// Next, if we have any streams, check if the Hash changed.
		// FIXME: This is where a design flaw of having flat Streams in a Matchthread shows.
		// 		  When we detect the first stream with a Hash change, we will have to get rid of
		//		  all streams and parse them again. I don't like this and maybe we should switch
		//		  to using Matchthread -> Comment -> Streams.
		if len(commentStreams) == 0 {
			continue
		}

		// Update upvotes
		if c.Ups != commentStreams[0].Metadata.Upvotes {
			updated = true

			for _, cs := range commentStreams {
				cs.Metadata.Upvotes = c.Ups
			}
		}

		// Check stream for changes
		hash := fmt.Sprintf("%x", md5.Sum([]byte(c.Body)))
		if hash == commentStreams[0].Metadata.Hash {
			streams = append(streams, commentStreams...)
			continue
		}

		updated = true

		streamLogger.WithFields(log.Fields{
			"old": commentStreams[0].Metadata.Hash,
			"new": hash,
		}).Debugf("Comment hash changed, updating all streams")

		newStreams := parser.ParseComment(c)

		for _, ns := range newStreams {
			ns.FillMetadata(c)
		}

		streamLogger.WithFields(log.Fields{
			"old": len(commentStreams),
			"new": len(newStreams),
		}).Debugf("Adding new streams after comment hash change")
		streams = append(streams, newStreams...)
	}

	if updated {
		logger.WithFields(log.Fields{
			"streams_before": len(mt.Streams),
			"streams_after":  len(streams),
		}).Debugf("Updating matchthread after polling update")

		mt.Streams = streams
		if err := mt.Save(); err != nil {
			logger.Errorf("Could not save matchthread: %v", err)
			return true, err
		}
	}

	return true, nil
}
