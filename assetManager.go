package hhot

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"os"

	l "github.com/SamHennessy/hlive"
	"github.com/rs/zerolog"
)

func NewAssetManager(sl ServiceLocator) *AssetManager {
	am := &AssetManager{
		logger: sl.Logger(),
		config: sl.Config(),
	}

	if contents, err := os.ReadFile("./assets/dist/js/app.js"); err != nil {
		am.logger.Err(err).Msg("asset manager: read app.js")
	} else {
		am.js = contents
		am.hashJS = fmt.Sprintf("%x", sha1.Sum(contents))
	}

	if contents, err := os.ReadFile("./assets/dist/css/app.css"); err != nil {
		am.logger.Err(err).Msg("asset manager: read app.css")
	} else {
		am.css = contents
		am.hashCSS = fmt.Sprintf("%x", sha1.Sum(contents))
	}

	return am
}

type AssetManager struct {
	logger  *zerolog.Logger
	config  Config
	js      []byte
	hashJS  string
	css     []byte
	hashCSS string
}

func (am *AssetManager) path(path string) string {
	return am.config.BasePath() + path
}

func (am *AssetManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == am.path("css/app.ClassBool") {
		w.Header().Add("Content-Type", "text/css")
		_, err := w.Write(am.css)
		if err != nil {
			am.logger.Err(err).Msg("serve app.ClassBool")
		}

		return
	}

	if r.URL.Path == am.path("js/app.js") {
		w.Header().Add("Content-Type", "text/javascript")

		_, err := w.Write(am.js)
		if err != nil {
			am.logger.Err(err).Msg("serve app.js")
		}

		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func (am *AssetManager) Tags() *l.NodeGroup {
	return l.G(
		l.T("link", l.Attrs{"rel": "stylesheet", "href": am.path("css/app.ClassBool?v=" + am.hashCSS)}),
		l.T("script", l.Attrs{"src": am.path("js/app.js?v=" + am.hashJS), "defer": ""}),
	)
}
