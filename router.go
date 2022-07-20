package hhot

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	l "github.com/SamHennessy/hlive"
	"github.com/SamHennessy/hlive/hlivekit"
	"github.com/rs/zerolog"
)

const routePartWild = "<wild>"

const (
	TopicRedirectInternal        = "hhot_redirect_internal"
	TopicRedirectInternalHistory = "hhot_redirect_internal_history"
)

type route struct {
	path      string
	parts     []string
	methods   []string
	allowPost bool
	allowGet  bool
	pageFn    func() *l.Page
}

func (r route) IsZero() bool {
	return r.pageFn == nil
}

func (r route) MethodAllowed(method string) bool {
	if len(r.methods) == 0 {
		return true
	}

	method = strings.ToUpper(method)
	if method == "GET" {
		return r.allowGet
	}

	if method == "POST" {
		return r.allowPost
	}

	for i := 0; i < len(r.methods); i++ {
		if method == r.methods[i] {
			return true
		}
	}

	return false
}

type routeTreeNode struct {
	kids  map[string]*routeTreeNode
	route route
}

func NewRouter(sl ServiceLocator) *Router {
	return &Router{
		routes: map[string]route{},
		tree: &routeTreeNode{
			kids: map[string]*routeTreeNode{},
		},
		store:  sl.PageSessionStore(),
		config: sl.Config(),
		log:    sl.Logger(),
	}
}

type Router struct {
	routes map[string]route
	tree   *routeTreeNode
	store  *l.PageSessionStore
	config Config
	log    *zerolog.Logger
}

func (r *Router) normalisePath(path string) string {
	return strings.Trim(path, "/")
}

func (r *Router) pathToParts(path string) (string, []string) {
	return path, strings.Split(r.normalisePath(path), "/")
}

func (r *Router) Add(path string, pageFn func() *l.Page, methods ...string) {
	path, parts := r.pathToParts(path)

	newRoute := route{
		path:   path,
		parts:  parts,
		pageFn: pageFn,
	}

	for i := 0; i < len(methods); i++ {
		m := strings.ToUpper(methods[i])
		newRoute.methods = append(newRoute.methods, m)

		if m == "GET" {
			newRoute.allowGet = true
		}

		if m == "POST" {
			newRoute.allowPost = true
		}
	}

	r.routes[path] = newRoute

	r.addRouteToTree(newRoute)

}

func (r *Router) addRouteToTree(rout route) {
	parent := r.tree

	for i := 0; i < len(rout.parts); i++ {
		part := rout.parts[i]
		if strings.HasPrefix(part, "{") {
			part = routePartWild
		}

		kid := parent.kids[part]
		if kid == nil {
			kid = &routeTreeNode{
				kids: map[string]*routeTreeNode{},
			}

			parent.kids[part] = kid
		}

		// End?
		if i == len(rout.parts)-1 {
			kid.route = rout
			parent.kids[part] = kid
		}

		// next level down
		parent = kid
	}
}

func (r *Router) matchInternal(path string) (normalisePath string, parts []string, rout route) {
	normalisePath = r.normalisePath(path)
	// Remove base path if needed
	bp := strings.Trim(r.config.BasePath(), "/")
	if bp != "" {
		// TODO: maybe we should 404 if it doesn't begin with the base path
		if strings.HasPrefix(normalisePath, bp) {
			normalisePath = normalisePath[len(bp):]
		}
	}

	_, parts = r.pathToParts(normalisePath)

	parent := r.tree

	for i := 0; i < len(parts); i++ {
		part := parts[i]

		kid, exists := parent.kids[part]
		if !exists {
			// Variable?
			kid, exists = parent.kids[routePartWild]
			if !exists {
				return
			}
		}

		// End?
		if i == len(parts)-1 {
			rout = kid.route
			return
		}

		// Next level down
		parent = kid
	}

	return
}

type ctxKey string

const ctxKeyRouteInfo ctxKey = "r"

type RouteData struct {
	Path       string
	PathParts  []string
	PathParams map[string]string
}

func RouteDataFromContext(ctx context.Context) *RouteData {
	data, _ := ctx.Value(ctxKeyRouteInfo).(*RouteData)

	return data
}

func (r *Router) routeDataFromRequest(path string, parts []string, rout route) *RouteData {
	// Should not happen but just in case
	if len(parts) != len(rout.parts) {
		return nil
	}

	data := &RouteData{
		Path:       path,
		PathParts:  parts,
		PathParams: map[string]string{},
	}

	for i := 0; i < len(rout.parts); i++ {
		if strings.HasPrefix(rout.parts[i], "{") {
			data.PathParams[strings.Trim(rout.parts[i], "{}")] = parts[i]
		}
	}

	return data
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path, parts, rout := r.matchInternal(req.URL.Path)

	if rout.IsZero() {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Page Not Found (404)"))

		return
	}

	ctx := context.WithValue(req.Context(), ctxKeyRouteInfo, r.routeDataFromRequest(path, parts, rout))

	l.NewPageServerWithSessionStore(rout.pageFn, r.store).ServeHTTP(w, req.WithContext(ctx))
}

func (r *Router) ReplacePage(newPath string, oldPage *l.Page, isHistory bool) {
	path, parts, rout := r.matchInternal(newPath)

	// Found?
	if rout.IsZero() {
		r.log.Error().Str("route", newPath).Msg("route not found")

		return
	}

	sess := r.store.Get(oldPage.GetSessionID())
	if sess == nil {
		r.log.Error().Str("id", oldPage.GetSessionID()).Msg("page session not found")

		return
	}

	if oldPage.IsConnected() {
		oldPage.Close(sess.InitialContext)
	}

	newPage := rout.pageFn()
	newPage.DOMBrowser = oldPage.DOMBrowser
	sess.Page = newPage

	// Browser History
	if !isHistory {
		HistoryPush(path, newPage.DOM.HTML)
	}

	ctx := context.WithValue(sess.InitialContext, ctxKeyRouteInfo, r.routeDataFromRequest(path, parts, rout))

	err := newPage.ServeWS(ctx, sess.ID, sess.Send, sess.Receive)
	if err != nil {
		fmt.Println("ERROR: Replace Page: ServerWS: ", err)
		return
	}
}

// TODO: Google says these are not crawlable
func Link(path string, elements ...any) *InternalRoute {
	a := &InternalRoute{
		Component: l.C("a",
			l.PreventDefault(),
			l.Attrs{"href": path},
			elements,
		),
	}

	a.Add(l.On("click", func(ctx context.Context, e l.Event) {
		InternalRedirect(a.pubSub, path)
	}))

	return a
}

type InternalRoute struct {
	*l.Component

	pubSub *hlivekit.PubSub
}

func (a *InternalRoute) PubSubMount(_ context.Context, pubSub *hlivekit.PubSub) {
	a.pubSub = pubSub
}

func InternalRedirect(pubSub *hlivekit.PubSub, path string) {
	pubSub.Publish(TopicRedirectInternal, path)
}

func InternalRedirectListener(pubSub *hlivekit.PubSub, hhr *Router, page *l.Page) (unsubscribe func()) {
	subFn := pubSub.SubscribeFunc(func(message hlivekit.QueueMessage) {
		path, ok := message.Value.(string)
		if !ok {
			return
		}

		go hhr.ReplacePage(path, page, message.Topic == TopicRedirectInternalHistory)
	}, TopicRedirectInternal, TopicRedirectInternalHistory)

	return func() {
		pubSub.Unsubscribe(subFn, TopicRedirectInternal, TopicRedirectInternalHistory)
	}
}
