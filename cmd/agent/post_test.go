package main

import (
	"testing"

	"github.com/Rukenshia/graw/reddit"
	mockdb "github.com/Rukenshia/soccerstreams/pkg/soccerstreams/testing"
	"github.com/stretchr/testify/assert"
)

type TestBot struct{}

func (t *TestBot) Reply(string, string) error {
	return nil
}
func (t *TestBot) SendMessage(string, string, string) error {
	return nil
}
func (t *TestBot) PostSelf(string, string, string) error {
	return nil
}
func (t *TestBot) PostLink(string, string, string) error {
	return nil
}
func (t *TestBot) Thread(string) (*reddit.Post, error) {
	return nil, nil
}
func (t *TestBot) Listing(string, string) (reddit.Harvest, error) {
	return reddit.Harvest{}, nil
}
func (t *TestBot) ListingWithParams(string, map[string]string) (reddit.Harvest, error) {
	return reddit.Harvest{}, nil
}

func TestPostSuccessfulParsing(t *testing.T) {
	db := mockdb.NewMockDBClient()
	a := NewAgent(&TestBot{}, db)

	assert.NoError(t, a.Post(&reddit.Post{
		ID:                "testing",
		Author:            "testuser",
		Title:             "[20:00 GMT] Team A vs Team B",
		LinkFlairCSSClass: "premierleague",
		LinkFlairText:     "English Premier League",
	}))

	mt, ok := db.Threads["testing"]
	assert.Equal(t, true, ok)
	assert.Equal(t, "Team A", mt.Team1)
	assert.Equal(t, "Team B", mt.Team2)
	assert.Equal(t, "English Premier League", mt.CompetitionName)
	assert.Equal(t, "premierleague", mt.Competition)
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
