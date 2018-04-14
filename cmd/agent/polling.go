package main

import (
	"time"

	"github.com/Rukenshia/soccerstreams/cmd/agent/metrics"
	log "github.com/sirupsen/logrus"
	"github.com/turnage/graw/reddit"
)

// ThreadPoller represents a Matchthread that is being checked for updates
type ThreadPoller struct {
	Permalink string
	Kickoff   time.Time

	bot reddit.Bot
}

// NewThreadPoller creates a new instance of ThreadPoller. It needs to be started using Poll before actually polling.
func NewThreadPoller(permalink string, kickoff time.Time, bot reddit.Bot) *ThreadPoller {
	return &ThreadPoller{
		permalink,
		kickoff,
		bot,
	}
}

// Poll starts polling the thread for updates.
func (t *ThreadPoller) Poll() chan *reddit.Post {
	ticker := time.NewTicker(time.Minute * 5)
	updates := make(chan *reddit.Post)

	if time.Now().After(t.Kickoff.Add(time.Minute * 60)) {
		close(updates)
		return updates
	}

	go func() {
		defer func() {
			recover()
		}()

		for _ = range ticker.C {
			if time.Now().After(t.Kickoff.Add(time.Minute * 60)) {
				close(updates)
				break
			}

			post, err := t.bot.Thread(t.Permalink)
			metrics.GrawEvents.WithLabelValues("manual_thread").Inc()
			if err != nil {
				metrics.GrawEvents.WithLabelValues("manual_thread_error").Inc()
				log.Warnf("Could not poll thread: %v", err)
				close(updates)
				break
			}

			updates <- post
		}
	}()

	return updates
}
