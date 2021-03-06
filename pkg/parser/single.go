package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"
)

type singleStreamParser struct {
}

// required format
// https://i.imgur.com/vGlRVgj.png
// QUALITY | [ NAME ](LINK) | LANGUAGE | MISR | NSFW | DISABLE AD_BLOCK | CLICKS | MOBILE COMPATIBLE

// We are also implementing this weird kind of format used by apparently a lot of streamers
// QUALITY | [ NAME | LANGUAGE | MISR | NSFW | DISABLE AD_BLOCK | CLICKS | MOBILE COMPATIBLE] (LINK)

func (w *singleStreamParser) Parse(message string) []*soccerstreams.Stream {
	var streams []*soccerstreams.Stream

	for _, line := range strings.Split(message, "\n") {
		if stream := w.parseLine(line); stream.IsGood() {
			streams = append(streams, stream)
		}
	}

	return streams
}

func (w *singleStreamParser) parseLine(line string) *soccerstreams.Stream {
	var s soccerstreams.Stream

	newLine, link, mightHaveFragments, is := w.parseChannel(line)
	if is {
		if mightHaveFragments {
			subData := w.parseLine(newLine)

			if subData != nil {
				s = *subData
			}
			s.Channel = newLine[:strings.Index(newLine, "|")-1]
			s.Link = link
		}
	}

	for _, fragment := range strings.Split(line, "|") {
		fragment = strings.TrimSpace(fragment)

		if q, is := w.parseQuality(fragment); is {
			s.Quality = q
		}
		if !mightHaveFragments {
			if n, l, _, is := w.parseChannel(fragment); is {
				s.Link = l
				s.Channel = n
			}
		}
		if m, is := w.parseMISR(fragment); is {
			s.MISR = m
		}
		if c, is := w.parseClicks(fragment); is {
			s.Clicks = c
		}
		if mf, is := w.parseMobileFriendly(fragment); is {
			s.MobileFriendly = mf
		}

	}

	return &s
}

func (w *singleStreamParser) parseQuality(fragment string) (string, bool) {
	if fragment == "HD" || fragment == "**HD**" {
		return "HD", true
	} else if fragment == "SD" {
		return "SD", true
	} else if fragment == "520p" {
		return "520p", true
	}

	return "", false
}

func (w *singleStreamParser) parseChannel(fragment string) (string, string, bool, bool) {
	re := regexp.MustCompile(`\[\s?(.*?)\s?\]\s?\(\s?(.*?)\s?\)`)

	groups := re.FindStringSubmatch(fragment)
	if len(groups) > 0 {
		return groups[1], groups[2], strings.Contains(groups[1], "|"), true
	}

	return "", "", false, false
}

func (w *singleStreamParser) parseMISR(fragment string) (string, bool) {
	re := regexp.MustCompile(`MISR\s?:?\s?(.*)`)

	groups := re.FindStringSubmatch(fragment)
	if len(groups) > 0 {
		return groups[1], true
	}

	return "", false
}

func (w *singleStreamParser) parseClicks(fragment string) (int, bool) {
	re := regexp.MustCompile(`Clicks\s?:?\s?([0-9]+)`)

	groups := re.FindStringSubmatch(fragment)
	if len(groups) > 0 {
		num, err := strconv.Atoi(groups[1])
		if err != nil {
			return 0, false
		}
		return num, true
	}

	return 0, false
}

func (w *singleStreamParser) parseMobileFriendly(fragment string) (bool, bool) {
	re := regexp.MustCompile(`(?i)Mobile.*?:?\s?(yes|no)`)

	groups := re.FindStringSubmatch(fragment)
	if len(groups) > 0 {
		groups[1] = strings.ToLower(groups[1])
		if groups[1] == "yes" {
			return true, true
		} else if groups[1] == "no" {
			return false, true
		}
	}

	return false, false
}
