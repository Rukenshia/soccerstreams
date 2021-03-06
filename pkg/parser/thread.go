package parser

import (
	"regexp"
	"strconv"
	"time"

	"github.com/Rukenshia/graw/reddit"
	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"
	log "github.com/sirupsen/logrus"
)

type threadParser struct {
	logger *log.Entry
}

func (t *threadParser) Parse(p *reddit.Post) *soccerstreams.Matchthread {
	var mt soccerstreams.Matchthread

	t.logger = log.WithField("post_id", p.ID).
		WithField("author", p.Author).
		WithField("title", p.Title)

	if kickoff, is := t.parseTime(p.Title); is {
		mt.Kickoff = kickoff
	} else {
		t.logger.Debug("could not parse kickoff, bailing")
		return nil
	}

	if t1, t2, is := t.parseTeams(p.Title); is {
		mt.Team1 = t1
		mt.Team2 = t2
	} else {
		t.logger.Debug("could not parse teams, bailing")
		return nil
	}

	return &mt
}

func (t *threadParser) parseTime(title string) (*time.Time, bool) {
	re := regexp.MustCompile(`\[([0-2]?[0-9])(?:\:|\.)([0-5]?[0-9])\s?GMT\]`)

	groups := re.FindStringSubmatch(title)
	if len(groups) > 0 {
		hour, err := strconv.Atoi(groups[1])
		if err != nil {
			return nil, false
		}

		min, err := strconv.Atoi(groups[2])
		if err != nil {
			return nil, false
		}

		gmt, err := time.LoadLocation("GMT")
		if err != nil {
			return nil, false
		}
		now := time.Now().In(gmt)

		kickoff := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, gmt)

		return &kickoff, true
	}

	return nil, false
}

func (t *threadParser) parseTeams(title string) (string, string, bool) {
	re := regexp.MustCompile(`(?i)\]\s?(.*?) (?:vs?\.?|x|-) (.*)`)

	groups := re.FindStringSubmatch(title)
	if len(groups) > 0 {
		return groups[1], groups[2], true
	}

	return "", "", false
}
