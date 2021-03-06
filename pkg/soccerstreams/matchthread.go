package soccerstreams

import (
	"time"
)

// Matchthread represents a reddit thread containing Streams.
// They are usually posted ~1 hour before Kickoff (GMT) and contain basic info (Team 1 vs Team 2)
type Matchthread struct {
	Team1           string
	Team2           string
	Competition     string
	CompetitionName string
	Kickoff         *time.Time
	Comments        []*Comment

	UpdatedAt time.Time
	ExpiresAt time.Time
	RedditID  string

	client DBClient
}

// NewMatchthread creates a new Matchthread that can be persisted
func NewMatchthread(client DBClient) *Matchthread {
	return &Matchthread{
		client: client,
	}
}

// SetClient sets the database client for the Matchthread
func (m *Matchthread) SetClient(d DBClient) {
	m.client = d
}

// DBKey returns the database key(id) that should be used for persistence
func (m *Matchthread) DBKey() string {
	return m.RedditID
}

// AddComment adds a Comment to a Matchthread if it does not exist yet
func (m *Matchthread) AddComment(c *Comment) bool {
	for _, ec := range m.Comments {
		if ec.RedditID == c.RedditID {
			return false
		}
	}

	m.Comments = append(m.Comments, c)
	return true
}

// UpdateExpiresAt sets the Matchthreads expiry time. A football game usually takes 1 1/2 hours so we give it a buffer of another hour
// after kickoff
func (m *Matchthread) UpdateExpiresAt() {
	m.ExpiresAt = m.Kickoff.Add(time.Hour*2 + time.Minute*30)
}

// Save saves the matchthread
func (m *Matchthread) Save() error {
	m.UpdatedAt = time.Now()

	return m.client.Upsert(m)
}

// Delete deletes the matchthread
func (m *Matchthread) Delete() error {
	return m.client.Delete(m.DBKey())
}
