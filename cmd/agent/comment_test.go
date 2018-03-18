package main

import (
	"testing"
	"time"

	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"

	mockdb "github.com/Rukenshia/soccerstreams/pkg/soccerstreams/testing"
	"github.com/stretchr/testify/assert"
	"github.com/turnage/graw/reddit"
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

	assert.Len(t, mt.Streams, 1)

	s := mt.Streams[0]

	assert.Equal(t, "ENGLISH TSN4 1080p", s.Channel)
	assert.Equal(t, "HD", s.Quality)
	assert.Equal(t, "3mbps", s.MISR)
	assert.Equal(t, 2, s.Clicks)
	assert.Equal(t, "testuser", s.Streamer)
	assert.True(t, s.MobileFriendly)
	assert.True(t, s.Metadata.ReliableStreamer)
	assert.Equal(t, int32(1), s.Metadata.Upvotes)
}
