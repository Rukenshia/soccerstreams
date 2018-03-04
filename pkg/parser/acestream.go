package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Rukenshia/soccerstreams/pkg/monitoring"
	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"
	raven "github.com/getsentry/raven-go"
	log "github.com/sirupsen/logrus"
)

type acestreamParser struct{}

// required format
// https://i.imgur.com/OzITykQ.png
// acestream://CONTENT-ID [QUALITY] [CHANNEL NAME] [LANGUAGE]

func (a *acestreamParser) Parse(comment string) []*soccerstreams.Stream {
	var streams []*soccerstreams.Stream

	for _, line := range strings.Split(comment, "\n") {
		if stream := a.parseLine(line); stream != nil {
			log.Debugf("Found acestream: %s", stream.Link)
			streams = append(streams, stream)
		}
	}

	return streams
}

func (a *acestreamParser) parseLine(line string) *soccerstreams.Stream {
	l, is := a.parseLink(line)

	if !is {
		return nil
	}

	s := &soccerstreams.Stream{}
	s.Link = fmt.Sprintf("acestream://%s", l)
	s.Channel = "Acestream"

	if q, matchedString, is := a.parseQuality(line); is {
		s.Quality = q

		// having quality means we can rule out one candidate for channel name
		// lets find the others
		rest := strings.Replace(line, matchedString, "", 1)
		re := regexp.MustCompile(`\[\s?(.{3,}?)\s?\]`)

		if _, matchedLanguage, is := a.parseLanguage(rest); is {
			rest = strings.Replace(rest, matchedLanguage, "", 1)
		}

		matches := re.FindAllStringSubmatch(rest, -1)

		var candidates []string

		for _, m := range matches {
			candidates = append(candidates, m[1])
		}

		if len(candidates) > 1 {
			// For the sake of having metrics for failed acestreams, directly report this to sentry
			raven.Capture(logging.CreatePacket(raven.DEBUG, "Too many channel name candidates\n Candidates are: ['%s']", strings.Join(candidates, "', '")), nil)
		} else if len(candidates) == 1 {
			s.Channel = candidates[0]
		}
	}

	return s
}

func (a *acestreamParser) parseLink(line string) (string, bool) {
	re := regexp.MustCompile(`acestream:\/\/([0-9a-z]+)`)

	groups := re.FindStringSubmatch(line)
	if len(groups) > 0 {
		return groups[1], true
	}

	return "", false
}

func (a *acestreamParser) parseQuality(line string) (string, string, bool) {
	re := regexp.MustCompile(`(?i)\[\s?(HD|SD|720p?|520p?|1920x1080|1080p?)\s?\]`)

	groups := re.FindStringSubmatch(line)
	if len(groups) > 0 {
		if groups[1][len(groups[1])-1] == 'P' {
			groups[1] = strings.Replace(groups[1], "P", "p", -1)
		}

		switch groups[1] {
		case "1920x1080", "1080p":
			groups[1] = "HD"
		}

		return groups[1], groups[0], true
	}

	return "", "", false
}

func (a *acestreamParser) parseLanguage(line string) (string, string, bool) {
	re := regexp.MustCompile(`(?i)\[\s?(EN(GLISH)?|RU(SSIAN)?|UKRANIAN|GERMAN|DE(UTSCH)?)\s?\]?`)

	groups := re.FindStringSubmatch(line)
	if len(groups) > 0 {
		switch strings.ToUpper(groups[1]) {
		case "ENGLISH":
			groups[1] = "EN"
		case "RUSSIAN":
			groups[1] = "RU"
		case "GERMAN", "DEUTSCH":
			groups[1] = "DE"
		}
		return groups[1], groups[0], true
	}

	return "", "", false
}
