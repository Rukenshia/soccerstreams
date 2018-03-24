package controllers

import (
	"fmt"
	"sort"
	"time"

	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"

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

	thread.NumAcestreams = 0
	thread.NumWebstreams = 0

	for _, s := range thread.Comments {
		thread.NumAcestreams += len(s.Acestreams)
		thread.NumWebstreams += len(s.Webstreams)
	}

	thread.NumStreams = thread.NumAcestreams + thread.NumWebstreams

	sort.Sort(models.ByCommentRelevance(thread.Comments))
	sort.Sort(models.ByCommentRelevance(thread.Comments))
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

	var backendThreads []*soccerstreams.Matchthread

	if err := cache.Get("matchthreads", &threads); err != nil {
		query := datastore.NewQuery("matchthread")
		if _, err := app.DB.GetAll(app.DBCtx, query, &backendThreads); err != nil {
			return c.handleDbError(err, revel.NewErrorFromPanic(err))
		}

		for _, backendThread := range backendThreads {
			fc := &models.FrontendMatchthread{
				Matchthread: backendThread,
			}

			var frontendComments []*models.FrontendComment
			for _, com := range backendThread.Comments {
				frontendComments = append(frontendComments, models.NewFrontendComment(com))
			}

			fc.Comments = frontendComments

			processThread(fc)

			threads = append(threads, fc)

			go cache.Set(fmt.Sprintf("matchthread_%s", fc.DBKey()), fc, 30*time.Second)
		}

		go cache.Set("matchthreads", threads, 30*time.Second)
	}

	moreStyles := []string{"css/soc.css"}

	competitions := threads.ByCompetition()
	return c.Render(competitions, threads, moreStyles)
}

// Details renders a matchthread's details page with links to the individual streams
func (c App) Details() revel.Result {
	var thread *models.FrontendMatchthread

	cacheKey := fmt.Sprintf("matchthread_%s", c.Params.Route.Get("thread"))

	if err := cache.Get(cacheKey, &thread); err != nil {
		var backendThread soccerstreams.Matchthread
		if err := app.DB.Get(app.DBCtx, datastore.NameKey("matchthread", c.Params.Route.Get("thread"), nil), &backendThread); err != nil {
			return c.handleDbError(err, revel.NewErrorFromPanic(err))
		}

		var comments []*models.FrontendComment
		for _, com := range backendThread.Comments {
			comments = append(comments, models.NewFrontendComment(com))
		}

		thread = &models.FrontendMatchthread{
			Matchthread: &backendThread,
			Comments:    comments,
		}

		processThread(thread)

		go cache.Set(cacheKey, thread, 30*time.Second)
	}

	c.Log.Debugf("%v", thread)

	moreStyles := []string{"css/soc.css"}

	return c.Render(thread, moreStyles)
}
