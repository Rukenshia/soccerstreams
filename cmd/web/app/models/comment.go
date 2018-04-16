package models

import (
	"sort"
	"strings"

	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"
)

// FrontendComment represents a enhanced Comment for the views
type FrontendComment struct {
	*soccerstreams.Comment

	Acestreams []*soccerstreams.Stream
	Webstreams []*soccerstreams.Stream
}

// NewFrontendComment creates a FrontendComment and fills the struct with additional data for the views
func NewFrontendComment(c *soccerstreams.Comment) *FrontendComment {
	var ace []*soccerstreams.Stream
	var web []*soccerstreams.Stream

	for _, s := range c.Streams {
		if strings.Contains(s.Link, "acestream://") {
			ace = append(ace, s)

			continue
		}

		web = append(web, s)
	}

	sort.Sort(ByQuality(ace))
	sort.Sort(ByQuality(web))

	return &FrontendComment{
		Comment:    c,
		Acestreams: ace,
		Webstreams: web,
	}
}

// ByCommentRelevance combines ByUpvotes and StreamerOfTheWeek sorts
type ByCommentRelevance []*FrontendComment

func (b ByCommentRelevance) Len() int { return len(b) }

// Less determines which stream is more relevant. The most relevant stream is a SOTW with the most upvotes.
// All other SOTW follow, afterwards non-SOTW Streams follow sorted by decreasing Upvotes.
func (b ByCommentRelevance) Less(i, j int) bool {
	if b[i].ReliableStreamer {
		if ByUpvotes(b).Less(j, i) && b[j].ReliableStreamer {
			return false
		}
		return true
	}
	return ByUpvotes(b).Less(i, j) && !b[j].ReliableStreamer
}
func (b ByCommentRelevance) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

// ByUpvotes sorts the Comments by number of Upvotes
type ByUpvotes []*FrontendComment

func (b ByUpvotes) Len() int { return len(b) }
func (b ByUpvotes) Less(i, j int) bool {
	return b[i].Upvotes > b[j].Upvotes
}
func (b ByUpvotes) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
