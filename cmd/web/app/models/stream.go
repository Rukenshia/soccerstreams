package models

import "github.com/Rukenshia/soccerstreams/pkg/soccerstreams"

// ByUpvotes sorts the Streams by number of Upvotes
type ByUpvotes []*soccerstreams.Stream

func (b ByUpvotes) Len() int { return len(b) }
func (b ByUpvotes) Less(i, j int) bool {
	return b[i].Metadata.Upvotes > b[j].Metadata.Upvotes
}
func (b ByUpvotes) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
