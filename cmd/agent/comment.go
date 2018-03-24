package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Rukenshia/soccerstreams/cmd/agent/metrics"

	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"

	"github.com/Rukenshia/soccerstreams/pkg/monitoring"
	"github.com/Rukenshia/soccerstreams/pkg/parser"
	raven "github.com/getsentry/raven-go"
	"github.com/turnage/graw/reddit"

	log "github.com/sirupsen/logrus"
)

// Comment parses a new reddit comment
func (s *Agent) Comment(p *reddit.Comment) error {
	metrics.CommentsIngested.Inc()

	// We only care about top level comments
	if !p.IsTopLevel() {
		return nil
	}

	logger := log.WithField("comment_id", p.ID).
		WithField("post_id", p.ParentID).
		WithField("author", p.Author)

	if _, ok := s.guard[p.ParentID]; !ok {
		s.guard[p.ParentID] = &sync.Mutex{}
	}

	s.guard[p.ParentID].Lock()

	defer func() {
		s.guard[p.ParentID].Unlock()
	}()

	if p.Author == "AutoModerator" {
		handled, err := s.handleAutoModComment(p)
		if err != nil {
			logger.Errorf("Could not handle AutoModerator action: %v", err)
			return nil
		}

		if handled {
			metrics.CommentsDeleted.Inc()
			logger.Debugf("Backing off parsing comment, automoderator took action")
			return nil
		}
	}

	streams := parser.ParseComment(p.Body)
	if len(streams) > 0 {
		metrics.CommentsParsed.Inc()

		// Find the matchthread in datastore
		// We are not making this in a goroutine because there might be other
		// comments at the same time, breaking our Update for the datastore entry

		// post id starts with t3_, we need to get rid of that
		postID := p.ParentID[3:]

		mt, err := s.client.Get(postID)
		if err != nil {
			logger.Errorf("Could not grab from datastore, %v", err)
			logger.Debugf("Grabbing matchthread from database failed, attempting on-demand parse")

			// try to just grab it NOW
			post, err := s.bot.Thread(fmt.Sprintf("/r/soccerstreams/comments/%s", postID))
			if err != nil {
				logger.Errorf("Could not parse matchthread from Comment: %v, thread id %s", err, postID)

				raven.CaptureError(err, map[string]string{
					"PostID":    postID,
					"CommentID": p.ID,
				})
				return nil
			}
			mt = parser.ParsePost(post)
			if mt == nil {
				logger.Errorf("Unable to parse matchthread from comment")
				raven.Capture(
					logging.CreatePacket(raven.DEBUG, "Unable to parse matchthread from comment\nPermalink: https://reddit.com%s\nTitle: %s\nFlair: %s", post.Permalink, post.Title, post.LinkFlairText),
					map[string]string{
						"Author": post.Author,
						"PostID": post.ID,
					})

				return nil
			}

			mt.SetClient(s.client)
			FillMatchthreadInfo(mt, post)
			mt.UpdateExpiresAt()
			logger.Debugf("Parsed matchthread %s via comment", mt.DBKey())
		}

		s.StartPolling(mt)

		comment := &soccerstreams.Comment{
			Streams: streams,
		}

		FillCommentInfo(comment, p)
		comment.UpdateHash()

		if added := mt.AddComment(comment); !added {
			// Not adding duplicated comment
			return nil
		}

		if err := mt.Save(); err != nil {
			logger.Errorf("Could not update Matchthread: %v", err)

			raven.CaptureError(err, map[string]string{
				"PostID":    postID,
				"CommentID": p.ID,
			})
			return nil
		}

		return nil
	}

	logger.Debugf("parsing of comment was not successful")

	raven.Capture(
		logging.CreatePacket(raven.DEBUG, "Unable to parse comment\nPermalink: https://reddit.com%s\n\n--- Body ---\n%s", p.Permalink, p.Body),
		map[string]string{
			"PostID":    p.ParentID[3:],
			"CommentID": p.ID,
			"Author":    p.Author,
		})
	return nil

}

// handleAutoModComment First return value being true indicates that an action has been taken and no further processing should happen
func (s *Agent) handleAutoModComment(c *reddit.Comment) (bool, error) {
	if c.Author != "AutoModerator" {
		return false, nil
	}

	logger := log.WithField("comment_id", c.ID).
		WithField("post_id", c.ParentID).
		WithField("author", c.Author)

	if strings.Contains(c.Body, "Your match thread has been removed") {
		logger.Debugf("AutoModerator removed the matchthread, deleting it from our database")

		if err := s.client.Delete(c.ParentID[3:]); err != nil {
			logger.Errorf("Could not delete matchthread: %v", err)

			raven.CaptureError(err, map[string]string{
				"CommentID": c.ID,
				"PostID":    c.ParentID,
				"Trigger":   "AutoModRemoval",
			})
			return true, nil
		}
		return true, nil

	}

	return false, nil
}
