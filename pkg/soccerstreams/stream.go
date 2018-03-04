package soccerstreams

// A Stream is a link to media posted inside of a Matchthread. This struct contains all information that can be parsed from reddit.
// This mainly relies on the rules of /r/soccerstreams.
type Stream struct {
	CommentLink    string
	Streamer       string
	Quality        string
	Channel        string
	Link           string
	Clicks         int
	IsNSFW         bool
	MISR           string
	MobileFriendly bool
}

// IsGood returns whether the Stream has sufficient information to be used further
func (s *Stream) IsGood() bool {
	if s.Link == "" || s.Channel == "" {
		return false
	}

	return true
}
