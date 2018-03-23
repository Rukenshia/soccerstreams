package models

import "github.com/Rukenshia/soccerstreams/pkg/soccerstreams"

// ByCommentRelevance combines ByUpvotes and StreamerOfTheWeek sorts
type ByCommentRelevance []*soccerstreams.Comment

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
	return ByUpvotes(b).Less(i, j)
}
func (b ByCommentRelevance) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

// ByUpvotes sorts the Comments by number of Upvotes
type ByUpvotes []*soccerstreams.Comment

func (b ByUpvotes) Len() int { return len(b) }
func (b ByUpvotes) Less(i, j int) bool {
	return b[i].Upvotes > b[j].Upvotes
}
func (b ByUpvotes) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
