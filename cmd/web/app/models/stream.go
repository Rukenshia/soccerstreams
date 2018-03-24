package models

import (
	"strings"

	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"
)

var qualities = map[string]int{
	// lower = better
	"HD":    0,
	"1080p": 1,
	"720":   2,
	"720p":  3,
	"SD":    4,
}

// ByQuality sorts streams by their quality. If no relevancy is defined, they are compared lexicographically
type ByQuality []*soccerstreams.Stream

func (b ByQuality) Len() int { return len(b) }
func (b ByQuality) Less(i, j int) bool {
	iRel, iIsRel := qualities[b[i].Quality]
	jRel, jIsRel := qualities[b[j].Quality]

	if iIsRel && jIsRel {
		return iRel < jRel
	}

	if iIsRel {
		return true
	}

	if jIsRel {
		return false
	}

	return strings.Compare(b[i].Quality, b[j].Quality) < 1
}
func (b ByQuality) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
