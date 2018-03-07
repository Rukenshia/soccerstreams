package main

import (
	"fmt"
	"strings"

	"github.com/Rukenshia/soccerstreams/pkg/monitoring"
	"github.com/Rukenshia/soccerstreams/pkg/parser"
	raven "github.com/getsentry/raven-go"
	"github.com/turnage/graw/reddit"

	log "github.com/sirupsen/logrus"
)

// Comment parses a new reddit comment
func (s *Agent) Comment(p *reddit.Comment) error {
	// We only care about top level comments
	if !p.IsTopLevel() {
		return nil
	}

	logger := log.WithField("comment_id", p.ID).
		WithField("post_id", p.ParentID).
		WithField("author", p.Author)

	if p.Author == "AutoModerator" {
		handled, err := s.handleAutoModComment(p)
		if err != nil {
			logger.Errorf("Could not handle AutoModerator action: %v", err)
			return nil
		}

		if handled {
			logger.Debugf("Backing off parsing comment, automoderator took action")
			return nil
		}
	}

	streams := parser.ParseComment(p)
	if len(streams) > 0 {
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
			if post, err := s.bot.Thread(fmt.Sprintf("/r/soccerstreams/comments/%s", postID)); err != nil {
				logger.Errorf("Could not parse matchthread from Comment: %v, thread id %s", err, postID)

				raven.CaptureError(err, map[string]string{
					"PostID":    postID,
					"CommentID": p.ID,
				})
				return nil
			} else {
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
				mt.FillRedditInfo(post)
				logger.Debugf("Parsed matchthread %s via comment", mt.DBKey())
			}
		}

		logger.Debugf("Found %d new streams", len(streams))

		mt.Streams = append(mt.Streams, streams...)

		logger.Debugf("saving updated matchthread")

		if err := mt.Save(); err != nil {
			logger.Errorf("Could not update Matchthread: %v", err)

			raven.CaptureError(err, map[string]string{
				"PostID":    postID,
				"CommentID": p.ID,
			})
			return nil
		}

		logger.Debugf("saved updated matchthread")

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
