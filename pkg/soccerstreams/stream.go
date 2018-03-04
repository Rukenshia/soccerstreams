package soccerstreams

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

func (s *Stream) IsGood() bool {
	if s.Link == "" || s.Channel == "" {
		return false
	}

	return true
}
