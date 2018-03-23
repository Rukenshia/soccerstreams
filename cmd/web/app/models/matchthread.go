package models

import (
	"sort"
	"time"

	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"
)

// FrontendMatchthread represents a soccerstream Matchthread with additional fields for frontend templates
type FrontendMatchthread struct {
	*soccerstreams.Matchthread

	GMTKickoff    string
	IsLive        bool
	NumAcestreams int
	NumWebstreams int
	NumStreams    int
	Webstreams    []*soccerstreams.Stream
	Acestreams    []*soccerstreams.Stream
}

// FrontendMatchthreads represents a slice of FrontendMatchthread
type FrontendMatchthreads []*FrontendMatchthread

// ByKickoff sorts FrontendMatchthreads by the time from now to Kickoff. Kickoffs before "now" are always preferred
type ByKickoff FrontendMatchthreads

func (b ByKickoff) Len() int { return len(b) }
func (b ByKickoff) Less(i, j int) bool {
	now := time.Now()
	if b[i].Kickoff.After(now) {
		return b[i].Kickoff.After(*b[j].Kickoff)
	}
	return true
}
func (b ByKickoff) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

// ByHasComments sorts FrontendMatchthreads by whether they have a stream or not.
type ByHasComments FrontendMatchthreads

func (b ByHasComments) Len() int { return len(b) }
func (b ByHasComments) Less(i, j int) bool {
	return b[i].NumAcestreams > 0 || b[i].NumWebstreams > 0
}
func (b ByHasComments) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

// ByCompetition returns a map of the FrontendMatchthreads grouped by their competition
func (f FrontendMatchthreads) ByCompetition() []*Competition {
	competitions := make(map[string]*Competition)

	for _, mt := range f {
		if mt.Competition == "" {
			mt.Competition = "unknown"
			mt.CompetitionName = "Unknown Competition"
		}

		if _, ok := competitions[mt.Competition]; !ok {
			competitions[mt.Competition] = &Competition{
				Name:         mt.CompetitionName,
				Identifier:   mt.Competition,
				Matchthreads: FrontendMatchthreads{mt},
			}
		} else {
			competitions[mt.Competition].Matchthreads = append(competitions[mt.Competition].Matchthreads, mt)
		}
	}

	var competitionsSlice []*Competition
	for _, c := range competitions {
		competitionsSlice = append(competitionsSlice, c)

		sort.Sort(ByKickoff(c.Matchthreads))
		sort.Sort(ByHasComments(c.Matchthreads))
	}

	sort.Sort(ByRelevance(competitionsSlice))

	return competitionsSlice
}
