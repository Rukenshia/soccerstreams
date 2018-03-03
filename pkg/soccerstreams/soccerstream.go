package soccerstreams

type Soccerstream struct {
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

func (s *Soccerstream) IsGood() bool {
	if s.Link == "" || s.Channel == "" {
		return false
	}

	return true
}
