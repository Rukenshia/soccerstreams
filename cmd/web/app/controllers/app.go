package controllers

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/Rukenshia/soccerstreams/cmd/web/app"
	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
)

type App struct {
	*revel.Controller
}

type FrontendMatchthread struct {
	*soccerstreams.Matchthread

	GMTKickoff    string
	IsLive        bool
	NumAcestreams int
	NumWebstreams int
	Acestreams    []*soccerstreams.Stream
	Webstreams    []*soccerstreams.Stream
}

type ByKickoff []*FrontendMatchthread

func (b ByKickoff) Len() int { return len(b) }
func (b ByKickoff) Less(i, j int) bool {
	now := time.Now()
	return b[i].Kickoff.After(now) && b[i].Kickoff.Before(*b[j].Kickoff)
}
func (b ByKickoff) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

type ByHasStreams []*FrontendMatchthread

func (b ByHasStreams) Len() int { return len(b) }
func (b ByHasStreams) Less(i, j int) bool {
	return len(b[i].Streams) > 0
}
func (b ByHasStreams) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

func processThread(thread *FrontendMatchthread) {
	if time.Now().After(*thread.Kickoff) {
		thread.IsLive = true
	}

	gmt, _ := time.LoadLocation("GMT")
	gmtKickoff := thread.Kickoff.In(gmt)
	thread.GMTKickoff = gmtKickoff.Format("15:04")
	thread.Kickoff = &gmtKickoff

	for _, stream := range thread.Streams {
		if strings.Contains(stream.Link, "acestream://") {
			thread.NumAcestreams++

			thread.Acestreams = append(thread.Acestreams, stream)
		} else {
			thread.NumWebstreams++

			thread.Webstreams = append(thread.Webstreams, stream)
		}
	}
}

func (c App) handleDbError(err error, baseErr *revel.Error) revel.Result {
	if err == datastore.ErrNoSuchEntity {
		return c.NotFound("Invalid match thread. It might have expired")
	}
	baseErr.Title = "Database Error"
	return revel.ErrorResult{Error: baseErr}
}

func (c App) Index() revel.Result {
	// get all Matchthreads
	var threads []*FrontendMatchthread

	if err := cache.Get("matchthreads", &threads); err != nil {
		query := datastore.NewQuery("matchthread")
		if _, err := app.DB.GetAll(app.DBCtx, query, &threads); err != nil {
			return c.handleDbError(err, revel.NewErrorFromPanic(err))
		}

		for _, thread := range threads {
			processThread(thread)

			go cache.Set(fmt.Sprintf("matchthread_%s", thread.DBKey()), thread, 2*time.Minute)
		}

		sort.Sort(ByKickoff(threads))
		sort.Sort(ByHasStreams(threads))

		go cache.Set("matchthreads", threads, 2*time.Minute)
	}

	moreStyles := []string{"css/soc.css"}

	return c.Render(threads, moreStyles)
}

func (c App) Details() revel.Result {
	// get all Matchthreads
	thread := &FrontendMatchthread{}

	cacheKey := fmt.Sprintf("matchthread_%s", c.Params.Route.Get("thread"))

	if err := cache.Get(cacheKey, &thread); err != nil {
		if err := app.DB.Get(app.DBCtx, datastore.NameKey("matchthread", c.Params.Route.Get("thread"), nil), thread); err != nil {
			return c.handleDbError(err, revel.NewErrorFromPanic(err))
		}

		c.Log.Debugf("%v", thread)

		processThread(thread)

		go cache.Set(cacheKey, thread, 2*time.Minute)
	}

	moreStyles := []string{"css/soc.css"}

	return c.Render(thread, moreStyles)
}
