package models

import "github.com/Rukenshia/soccerstreams/pkg/soccerstreams"

// ByStreamRelevance combines ByUpvotes and StreamerOfTheWeek sorts
type ByStreamRelevance []*soccerstreams.Stream

func (b ByStreamRelevance) Len() int { return len(b) }

// Less determines which stream is more relevant. The most relevant stream is a SOTW with the most upvotes.
// All other SOTW follow, afterwards non-SOTW Streams follow sorted by decreasing Upvotes.
func (b ByStreamRelevance) Less(i, j int) bool {
	if b[i].Metadata.ReliableStreamer {
		if ByUpvotes(b).Less(j, i) && b[j].Metadata.ReliableStreamer {
			return false
		}
		return true
	}
	return ByUpvotes(b).Less(i, j)
}
func (b ByStreamRelevance) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

// ByUpvotes sorts the Streams by number of Upvotes
type ByUpvotes []*soccerstreams.Stream

func (b ByUpvotes) Len() int { return len(b) }
func (b ByUpvotes) Less(i, j int) bool {
	return b[i].Metadata.Upvotes > b[j].Metadata.Upvotes
}
func (b ByUpvotes) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
