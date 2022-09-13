package site

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/SamHennessy/hhot/hhotui/site/pages"
	"github.com/SamHennessy/hlive"
	"github.com/SamHennessy/hlive/hlivekit"
	"github.com/koding/websocketproxy"
)

// TODO: make auto gen

func Router(sl *ServiceLocatorSite) {
	remote, err := url.Parse("http://localhost:3000")
	if err != nil {
		panic(err)
	}
	remoteWS, err := url.Parse("ws://localhost:3000")
	if err != nil {
		panic(err)
	}

	handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/app" {
				r.URL.Path = "/"
			}

			if r.URL.Query().Get("hlive") != "" {
				websocketproxy.NewProxy(remoteWS).ServeHTTP(w, r)
			} else {
				waitForPort()
				r.Host = remote.Host
				p.ServeHTTP(w, r)
			}
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)

	// hhr := hhot.NewRouter(sl)
	// phm := hhot.NewPageHistoryManager(sl.Config())

	newPageWrap := func(pFn func() *hlive.Page) func() *hlive.Page {
		return func() *hlive.Page {
			page := pFn()

			// Page level AppPubSub not Site level
			pubSub := hlivekit.NewPubSub()

			page.DOM().HTML().Add(hlivekit.InstallPubSub(pubSub))

			return page
		}
	}

	http.Handle("/", hlive.NewPageServerWithSessionStore(newPageWrap(pages.Index(sl)), sl.PageSessionStore()))
	http.HandleFunc("/app/", handler(proxy))
}

// TODO: add timeout
func waitForPort() {
	for {
		timeout := time.Second
		conn, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", "3000"), timeout)
		if err != nil {
			// fmt.Println("Connecting error:", err)
		}
		if conn != nil {
			conn.Close()
			return
		}

		time.Sleep(time.Millisecond * 100)
	}
}
