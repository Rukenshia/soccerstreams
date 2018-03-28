package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Rukenshia/soccerstreams/cmd/web/app/controllers/club"
	"github.com/Rukenshia/soccerstreams/cmd/web/metrics"
	"github.com/revel/revel"
)

type Club struct {
	*revel.Controller
}

// Image serves a clubs logo
func (c Club) Image() revel.Result {
	path := strings.ToLower(c.Params.Route.Get("club"))

	if strings.Contains(path, "..") {
		return c.NotFound("Invalid path %s", path)
	}

	suffixes := []string{"svg", "png"}
	filename := ""
	path = club.NormaliseName(path)

	for _, s := range suffixes {
		fp := filepath.Join(revel.AppPath[:len(revel.AppPath)-4], "assets", "img", "clubs", fmt.Sprintf("%s.%s", path, s))

		if _, err := os.Stat(fp); err == nil {
			filename = fp
			break
		}
	}

	if filename == "" {
		filename = filepath.Join(revel.AppPath[:len(revel.AppPath)-4], "assets", "img", "clubs", "placeholder.png")
		metrics.ImageNotFound.WithLabelValues("club").Inc()
	}

	return c.RenderFileName(filename, revel.NoDisposition)
}
