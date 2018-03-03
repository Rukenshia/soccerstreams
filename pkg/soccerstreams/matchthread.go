package soccerstreams

import (
	"time"

	"github.com/turnage/graw/reddit"
)

type Matchthread struct {
	Team1   string
	Team2   string
	Kickoff *time.Time
	Streams []*Soccerstream

	UpdatedAt time.Time
	ExpiresAt time.Time
	RedditID  string

	client DBClient
}

func NewMatchthread(client DBClient) *Matchthread {
	return &Matchthread{
		client: client,
	}
}

func (m *Matchthread) SetClient(d DBClient) {
	m.client = d
}

func (m *Matchthread) DBKey() string {
	return m.RedditID
}

func (m *Matchthread) SetExpiresAt() {
	m.ExpiresAt = m.Kickoff.Add(time.Hour*2 + time.Minute*30)
}

func (m *Matchthread) Save() error {
	m.UpdatedAt = time.Now()

	return m.client.Upsert(m)
}

func (m *Matchthread) Delete() error {
	return m.client.Delete(m.DBKey())
}

func (m *Matchthread) FillRedditInfo(p *reddit.Post) {
	m.RedditID = p.ID
	m.SetExpiresAt()
}
