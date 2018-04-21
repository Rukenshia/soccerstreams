package main

import (
	"testing"
	"time"

	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"

	"github.com/Rukenshia/graw/reddit"
	mockdb "github.com/Rukenshia/soccerstreams/pkg/soccerstreams/testing"
	"github.com/stretchr/testify/assert"
)

func TestCommentSingleWebstream(t *testing.T) {
	db := mockdb.NewMockDBClient()
	a := NewAgent(&TestBot{}, db)

	mt := soccerstreams.NewMatchthread(db)
	mt.RedditID = "testmt"

	now := time.Now()
	mt.Kickoff = &now

	db.Add(mt)

	assert.NoError(t, a.Comment(&reddit.Comment{
		ParentID:            "t3_testmt",
		ID:                  "testcomment",
		Author:              "testuser",
		AuthorFlairCSSClass: "weekly",
		Ups:                 1,
		Body:                "**HD** | [ENGLISH TSN4 1080p] (http://foundationsports.com/crvcsk/) | MISR : 3mbps | Ad Overlay : 1 | Clicks : 2  | Mobile : Yes",
	}))

	assert.Len(t, mt.Comments, 1)
	assert.Len(t, mt.Comments[0].Streams, 1)

	c := mt.Comments[0]
	s := mt.Comments[0].Streams[0]

	assert.Equal(t, "ENGLISH TSN4 1080p", s.Channel)
	assert.Equal(t, "HD", s.Quality)
	assert.Equal(t, "3mbps", s.MISR)
	assert.Equal(t, 2, s.Clicks)
	assert.Equal(t, "testuser", c.Streamer)
	assert.True(t, s.MobileFriendly)
	assert.True(t, c.ReliableStreamer)
	assert.Equal(t, int32(1), c.Upvotes)
}
