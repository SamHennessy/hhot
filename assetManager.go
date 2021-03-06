package hhot

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"net/http"
	"os"
	"strings"

	l "github.com/SamHennessy/hlive"
	"github.com/rs/zerolog"
)

func NewAssetManager(sl ServiceLocator) *AssetManager {
	am := &AssetManager{
		logger: sl.Logger(),
		config: sl.Config(),
	}

	// TODO: stop this and use file server as it manages cache checks
	if contents, err := os.ReadFile("./assets/dist/js/app.js"); err != nil {
		am.logger.Warn().Err(err).Msg("asset manager: read app.js")
	} else {
		am.js = contents
		am.hashJS = fmt.Sprintf("%x", sha1.Sum(contents))
	}

	if contents, err := os.ReadFile("./assets/dist/css/app.css"); err != nil {
		am.logger.Err(err).Msg("asset manager: read app.css")
	} else {
		am.hashCSS = fmt.Sprintf("%x", sha1.Sum(contents))
	}

	return am
}

type AssetManager struct {
	logger  *zerolog.Logger
	config  Config
	js      []byte
	hashJS  string
	hashCSS string
}

func (am *AssetManager) path(path string) string {
	return am.config.BasePath() + path
}

func (am *AssetManager) pathFav(path string) string {
	return am.config.BasePath() + "img/favicons/" + path
}

func (am *AssetManager) pathFont(path string) string {
	return am.config.BasePath() + "fonts/" + path
}

func (am *AssetManager) pathImg(path string) string {
	return am.config.BasePath() + "img/" + path
}

func (am *AssetManager) pathJs(path string) string {
	return am.config.BasePath() + "js/" + path
}

func (am *AssetManager) pathCSS(path string) string {
	return am.config.BasePath() + "css/" + path
}

// https://www.emergeinteractive.com/insights/detail/the-essentials-of-favicons/
func (am *AssetManager) favicons(w http.ResponseWriter, r *http.Request) bool {
	if r.URL.Path == am.pathFav("site.webmanifest") {
		w.Header().Add("Content-Type", "application/manifest+json")

		fileContents, err := os.ReadFile("./assets/dist/img/favicons/site.webmanifest")
		if err != nil {
			am.logger.Err(err).Msg("read site.webmanifest")
		} else {
			fileContents = bytes.ReplaceAll(fileContents, []byte("[iconPath]"), []byte(am.pathFav("")))

			_, _ = w.Write(fileContents)
		}

		return true
	}

	if strings.HasPrefix(r.URL.Path, am.pathFav("")) {
		file := strings.TrimPrefix(r.URL.Path, am.pathFav(""))

		http.ServeFile(w, r, "./assets/dist/img/favicons/"+file)

		return true
	}

	return false
}

func (am *AssetManager) Tags() *l.NodeGroup {
	g := l.G(
		l.T("link", l.Attrs{"rel": "stylesheet", "href": am.path("css/app.css?v=" + am.hashCSS)}),
		l.T("link", l.Attrs{"rel": "apple-touch-icon", "sizes": "180x180", "href": am.pathFav("apple-touch-icon.png")}),
		l.T("link", l.Attrs{"rel": "icon", "type": "image/png", "sizes": "32x32", "href": am.pathFav("favicon-32x32.png")}),
		l.T("link", l.Attrs{"rel": "icon", "type": "image/png", "sizes": "16x16", "href": am.pathFav("favicon-16x16.png")}),
		l.T("link", l.Attrs{"rel": "manifest", "href": am.pathFav("site.webmanifest")}),
	)

	if len(am.js) > 0 {
		g.Add(l.T("script", l.Attrs{"src": am.path("js/app.js?v=" + am.hashJS), "defer": ""}))
	}

	return g
}

func (am *AssetManager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == am.path("css/app.css") {
			w.Header().Add("Content-Type", "text/css")

			if r.URL.Query().Get("v") != "" {
				w.Header().Set("Cache-Control", "public, max-age=604800, immutable")
			}

			file := strings.TrimPrefix(r.URL.Path, am.pathCSS(""))
			http.ServeFile(w, r, "./assets/dist/css/"+file)

			return
		}

		if r.URL.Path == am.path("js/app.js") {
			w.Header().Add("Content-Type", "text/javascript")

			if r.URL.Query().Get("v") != "" {
				w.Header().Set("Cache-Control", "public, max-age=604800, immutable")
			}

			_, err := w.Write(am.js)
			if err != nil {
				am.logger.Err(err).Msg("serve app.js")
			}

			return
		}

		if r.URL.Path == am.path("js/app.js.map") {
			w.Header().Add("Content-Type", "text/javascript")

			file := strings.TrimPrefix(r.URL.Path, am.pathJs(""))
			http.ServeFile(w, r, "./assets/dist/js/"+file)

			return
		}

		if am.favicons(w, r) {
			return
		}

		if am.fonts(w, r) {
			return
		}

		if am.images(w, r) {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (am *AssetManager) fonts(w http.ResponseWriter, r *http.Request) bool {
	if strings.HasPrefix(r.URL.Path, am.pathFont("")) {
		// Assume fonts don't change
		w.Header().Set("Cache-Control", "public, max-age=604800, immutable")

		file := strings.TrimPrefix(r.URL.Path, am.pathFont(""))
		http.ServeFile(w, r, "./assets/dist/fonts/"+file)

		return true
	}

	return false
}

func (am *AssetManager) images(w http.ResponseWriter, r *http.Request) bool {
	if strings.HasPrefix(r.URL.Path, am.pathImg("")) {
		if r.URL.Query().Get("v") != "" {
			w.Header().Set("Cache-Control", "public, max-age=604800, immutable")
		}

		file := strings.TrimPrefix(r.URL.Path, am.pathImg(""))
		http.ServeFile(w, r, "./assets/dist/img/"+file)

		return true
	}

	return false
}

// TODO: add hash?
func (am *AssetManager) ImageSrc(path string) string {
	return am.config.BasePath() + "img/" + path
}
