package main

import (
	"testing"

	mockdb "github.com/Rukenshia/soccerstreams/pkg/soccerstreams/testing"
	"github.com/stretchr/testify/assert"
	"github.com/turnage/graw/reddit"
)

func TestPostSuccessfulParsing(t *testing.T) {
	db := mockdb.NewMockDBClient()
	a := NewAgent(nil, db)

	assert.NoError(t, a.Post(&reddit.Post{
		ID:     "testing",
		Author: "testuser",
		Title:  "[20:00 GMT] Team A vs Team B",
	}))

	mt, ok := db.Threads["testing"]
	assert.Equal(t, true, ok)
	assert.Equal(t, "Team A", mt.Team1)
	assert.Equal(t, "Team B", mt.Team2)
	assert.Equal(t, 20, mt.Kickoff.Hour())
	assert.Equal(t, 0, mt.Kickoff.Minute())
}

func TestPostUnsuccessfulParsing(t *testing.T) {
	db := mockdb.NewMockDBClient()
	a := NewAgent(nil, db)

	assert.NoError(t, a.Post(&reddit.Post{
		ID:     "testing",
		Author: "testuser",
		Title:  "[20:00 GMT] Bundesliga Konferenz",
	}))

	_, ok := db.Threads["testing"]
	assert.Equal(t, false, ok)
}
