package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/revel/revel"
	"github.com/revel/revel/cache"
)

type Club struct {
	*revel.Controller
}

var simpleClubMappings = map[string]string{
	"milan":                  "ac milan",
	"atlético madrid":        "atletico madrid",
	"atalanta":               "atalanta bc",
	"sampdoria":              "uc sampdoria",
	"barcelona":              "fc barcelona",
	"inter milano":           "inter mailand",
	"bayern münchen":         "bayern munchen",
	"beşiktaş":               "besiktas",
	"bayern munich":          "bayern munchen",
	"brighton & hove albion": "brighton and hove albion",
}

var regexClubMappings = map[string]func([]string) string{
	`^(.*) u[0-9]{2}$`: func(m []string) string { return m[1] },
}

type InMemoryFile struct {
	Data        []byte
	ContentType string
}

func (r InMemoryFile) Apply(req *revel.Request, resp *revel.Response) {
	resp.WriteHeader(http.StatusOK, r.ContentType)
	resp.GetWriter().Write(r.Data)
}

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

		if perm, ok := simpleClubMappings[path]; ok {
			path = perm
		} else {
			for reStr, res := range regexClubMappings {
				re := regexp.MustCompile(reStr)
				if groups := re.FindStringSubmatch(path); len(groups) > 0 {
					path = res(groups)
				}
			}
		}

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
