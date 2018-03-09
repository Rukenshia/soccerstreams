package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Rukenshia/soccerstreams/cmd/web/app/controllers/club"
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
)

type Club struct {
	*revel.Controller
}

// InMemoryFile is probably a bad idea
type InMemoryFile struct {
	Data        []byte
	ContentType string
}

// Apply implements revel.Response and writes the memory file to the HTTP response
func (r InMemoryFile) Apply(req *revel.Request, resp *revel.Response) {
	resp.WriteHeader(http.StatusOK, r.ContentType)
	resp.GetWriter().Write(r.Data)
}

// Image serves a clubs logo
func (c Club) Image() revel.Result {
	path := strings.ToLower(c.Params.Route.Get("club"))

	var imgData *InMemoryFile
	if err := cache.Get(fmt.Sprintf("club_image_%s", path), &imgData); err != nil {
		if strings.Contains(path, "..") {
			return c.NotFound("Invalid path %s", path)
		}

		suffixes := []string{"svg", "png"}
		filename := ""
		suffix := "svg"
		path = club.NormaliseName(path)

		for _, s := range suffixes {
			fp := filepath.Join(revel.AppPath[:len(revel.AppPath)-4], "assets", "img", "clubs", fmt.Sprintf("%s.%s", path, s))

			if _, err := os.Stat(fp); err == nil {
				filename = fp
				suffix = s
				break
			}
		}

		if filename == "" {
			filename = filepath.Join(revel.AppPath[:len(revel.AppPath)-4], "assets", "img", "clubs", "placeholder.png")
			suffix = "png"
		}

		f, err := ioutil.ReadFile(filename)
		if err != nil {
			re := revel.NewErrorFromPanic(err)
			if re == nil {
				re = &revel.Error{
					Description: err.Error(),
				}
			}
			re.Title = "File error"
			return revel.ErrorResult{Error: re}
		}

		ctype := "image/png"

		switch suffix {
		case "svg":
			ctype = "image/svg+xml"
		}

		imgData = &InMemoryFile{ContentType: ctype, Data: f}
		go cache.Set(fmt.Sprintf("club_image_%s", path), imgData, 30*time.Minute)
	}

	c.Log.Debugf("Returning image with ctype %s", imgData.ContentType)

	return imgData
}
