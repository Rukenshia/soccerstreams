package models

import (
	"strings"
)

// Relevance of competitions based solely upon on my own preferences
var relevance = map[string]int{
	// lower = better
	"premierleague": 0,
	"bundesliga":    1,
	"laliga":        2,
}

// Competition represents a Football competition (e.g. Premier League)
type Competition struct {
	Name       string
	Identifier string

	Matchthreads FrontendMatchthreads
}

// ByRelevance sorts competitions by their relevance. If no relevancy is defined, they are compared lexicographically
type ByRelevance []*Competition

func (b ByRelevance) Len() int { return len(b) }
func (b ByRelevance) Less(i, j int) bool {
	iRel, iIsRel := relevance[b[i].Identifier]
	jRel, jIsRel := relevance[b[j].Identifier]

	if iIsRel && jIsRel {
		return iRel < jRel
	}

	if iIsRel {
		return true
	}

	if jIsRel {
		return false
	}

	return strings.Compare(b[i].Name, b[j].Name) < 1
}
func (b ByRelevance) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
