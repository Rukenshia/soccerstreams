package app

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"time"

	"github.com/Rukenshia/soccerstreams/cmd/web/app/models"
	"github.com/Rukenshia/soccerstreams/cmd/web/metrics"
	"github.com/prometheus/common/log"

	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"

	"github.com/revel/log15"
	"github.com/revel/revel/logger"

	raven "github.com/getsentry/raven-go"

	"cloud.google.com/go/datastore"
	humanize "github.com/dustin/go-humanize"
	"github.com/revel/revel"
	"google.golang.org/api/option"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string

	// DB GCloud Datastore Client
	DB *datastore.Client
	// DBCtx Database Context
	DBCtx = context.Background()
)

func initDB() {
	db, err := datastore.NewClient(DBCtx, "soccerstreams-web", option.WithServiceAccountFile("/opt/soccerstreams/gcloud/gcloud-service-account.json"))
	if err != nil {
		revel.AppLog.Fatal(err.Error())
		return
	}

	DB = db
	revel.AppLog.Debugf("Created DB Client", DB)
}

type sentryHandler struct{}

func (s *sentryHandler) Log(r *log15.Record) error {
	ctx, err := s.getCtxFromArray(r.Ctx)
	if err == nil {
		raven.CaptureError(errors.New(r.Msg), ctx)
	} else {
		raven.CaptureError(errors.New(r.Msg), nil)
	}

	return nil
}

func (s *sentryHandler) getCtxFromArray(ctx []interface{}) (map[string]string, error) {
	if len(ctx)%2 != 0 {
		return nil, errors.New("Ctx array length is odd")
	}

	ctxMap := make(map[string]string, len(ctx)/2)

	for i, v := range ctx {
		if i%2 == 0 {
			key, ok := v.(string)
			if !ok {
				return nil, errors.New("Ctx key not a string")
			}
			str, ok := ctx[i+1].(string)
			if !ok {
				return nil, errors.New("Ctx value not a string")
			}

			ctxMap[key] = str
		}
	}

	return ctxMap, nil
}

func init() {
	raven.SetDSN(os.Getenv("SENTRY_DSN"))

	go func() {
		metrics.Register()
		log.Fatal(metrics.Serve())
	}()

	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		func(c *revel.Controller, fc []revel.Filter) {
			defer func() {
				res, ok := c.Result.(revel.ErrorResult)
				if ok {
					raven.CaptureError(res.Error, map[string]string{
						"Status":    fmt.Sprintf("%d", c.Response.Status),
						"ClientIP":  c.ClientIP,
						"UserAgent": c.Request.UserAgent(),
						"URL":       c.Request.URL.String(),
					})
					return
				}
			}()
			fc[0](c, fc[1:])
		},
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		func(c *revel.Controller, fc []revel.Filter) {
			fc[0](c, fc[1:])
			metrics.HTTPRequests.WithLabelValues(fmt.Sprintf("%d", c.Response.Status)).Inc()
		},
		revel.ActionInvoker, // Invoke the action.
	}

	// Add custom template functions
	revel.TemplateFuncs["relative"] = func(t time.Time) string {
		return humanize.Time(t)
	}
	revel.TemplateFuncs["truncate"] = func(s string) string {
		if len(s) < 20 {
			return s
		}

		return fmt.Sprintf("%s...", s[:20])
	}
	revel.TemplateFuncs["safeURL"] = func(s string) template.URL {
		return template.URL(s)
	}
	revel.TemplateFuncs["randomNumber"] = func(from, to int) int {
		return from + rand.Intn(to-from)
	}

	type CommentStreams struct {
		*models.FrontendComment
		Streams []*soccerstreams.Stream
	}
	revel.TemplateFuncs["overrideStreams"] = func(streams []*soccerstreams.Stream, comment *models.FrontendComment) CommentStreams {
		return CommentStreams{
			FrontendComment: comment,
			Streams:         streams,
		}
	}

	logger.LogFunctionMap["sentry"] = func(c *logger.CompositeMultiHandler, options *logger.LogOptions) {
		c.SetHandler(&sentryHandler{}, false, logger.LvlError)
	}

	revel.OnAppStart(initDB)
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	c.Response.Out.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

//func ExampleStartupScript() {
//	// revel.DevMod and revel.RunMode work here
//	// Use this script to check for dev mode and set dev/prod startup scripts here!
//	if revel.DevMode == true {
//		// Dev mode
//	}
//}
