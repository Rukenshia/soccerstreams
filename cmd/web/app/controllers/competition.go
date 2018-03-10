package controllers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/revel/revel"
	"github.com/revel/revel/cache"
)

type Competition struct {
	*revel.Controller
}

// Image serves a competitions logo
func (c Competition) Image() revel.Result {
	path := strings.ToLower(c.Params.Route.Get("competition"))

	var imgData *InMemoryFile
	if err := cache.Get(fmt.Sprintf("competition_image_%s", path), &imgData); err != nil {
		if strings.Contains(path, "..") {
			return c.NotFound("Invalid path %s", path)
		}

		suffixes := []string{"svg", "png"}
		filename := ""
		suffix := "svg"

		for _, s := range suffixes {
			fp := filepath.Join(revel.AppPath[:len(revel.AppPath)-4], "assets", "img", "competitions", fmt.Sprintf("%s.%s", path, s))

			if _, err := os.Stat(fp); err == nil {
				filename = fp
				suffix = s
				break
			}
		}

		if filename == "" {
			filename = filepath.Join(revel.AppPath[:len(revel.AppPath)-4], "assets", "img", "competitions", "placeholder.png")
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
		go cache.Set(fmt.Sprintf("competition_image_%s", path), imgData, 30*time.Minute)
	}

	c.Log.Debugf("Returning image with ctype %s", imgData.ContentType)

	return imgData
}
