package controllers

import (
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/Rukenshia/soccerstreams/cmd/web/app"
	"github.com/Rukenshia/soccerstreams/cmd/web/app/models"
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
)

// App is the main controller for soc-web
type App struct {
	*revel.Controller
}

func processThread(thread *models.FrontendMatchthread) {
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

// Index renders the main site (dashboard) with a list of matchthreads
func (c App) Index() revel.Result {
	// get all Matchthreads
	var threads models.FrontendMatchthreads

	if err := cache.Get("matchthreads", &threads); err != nil {
		query := datastore.NewQuery("matchthread")
		if _, err := app.DB.GetAll(app.DBCtx, query, &threads); err != nil {
			return c.handleDbError(err, revel.NewErrorFromPanic(err))
		}

		for _, thread := range threads {
			processThread(thread)

			go cache.Set(fmt.Sprintf("matchthread_%s", thread.DBKey()), thread, 2*time.Minute)
		}

		go cache.Set("matchthreads", threads, 2*time.Minute)
	}

	moreStyles := []string{"css/soc.css"}

	competitions := threads.ByCompetition()
	return c.Render(competitions, threads, moreStyles)
}

// Details renders a matchthread's details page with links to the individual streams
func (c App) Details() revel.Result {
	thread := &models.FrontendMatchthread{}

	cacheKey := fmt.Sprintf("matchthread_%s", c.Params.Route.Get("thread"))

	if err := cache.Get(cacheKey, &thread); err != nil {
		if err := app.DB.Get(app.DBCtx, datastore.NameKey("matchthread", c.Params.Route.Get("thread"), nil), thread); err != nil {
			return c.handleDbError(err, revel.NewErrorFromPanic(err))
		}
		processThread(thread)

		go cache.Set(cacheKey, thread, 2*time.Minute)
	}

	moreStyles := []string{"css/soc.css"}

	return c.Render(thread, moreStyles)
}
