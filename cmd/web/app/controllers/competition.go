package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/revel/revel"
)

type Competition struct {
	*revel.Controller
}

// Image serves a competitions logo
func (c Competition) Image() revel.Result {
	path := strings.ToLower(c.Params.Route.Get("competition"))

	if strings.Contains(path, "..") {
		return c.NotFound("Invalid path %s", path)
	}

	suffixes := []string{"svg", "png"}
	filename := ""

	for _, s := range suffixes {
		fp := filepath.Join(revel.AppPath[:len(revel.AppPath)-4], "assets", "img", "competitions", fmt.Sprintf("%s.%s", path, s))

		if _, err := os.Stat(fp); err == nil {
			filename = fp
			break
		}
	}

	if filename == "" {
		filename = filepath.Join(revel.AppPath[:len(revel.AppPath)-4], "assets", "img", "clubs", "placeholder.png")
	}

	return c.RenderFileName(filename, revel.NoDisposition)
}
