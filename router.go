package webserver

import (
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type (
	// RouteNamespace groups routes according to a specific URL entry point or prefix.
	RouteNamespace struct {
		prefix string
		server *Server

		DebugLogger   *log.Logger
		InfoLogger    *log.Logger
		WarningLogger *log.Logger
		ErrorLogger   *log.Logger
	}
)

func (rns *RouteNamespace) buildPath(p string) string {
	return path.Join(rns.prefix, p)
}

// Handle registers handlers!
func (rns *RouteNamespace) Handle(method string, path string, handlers []HandlerFunc) {
	p := rns.buildPath(path)

	rns.InfoLogger.Printf(logprefix+"Registering handler; Route: %s:%s; Quantity: %b", method, p, len(handlers))

	// Serve the request
	rns.server.router.Handle(method, path, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		rns.DebugLogger.Printf(logprefix+"Capturing request; Route: %s:%s", method, path)
		event := rns.server.captureRequest(w, req, nil, handlers)

		// TODO Execute all handlers passed in their order stopping if one
		// chooses to write to the body unless we can/should simply append
		// the event body.
		// handlers[0](event)
		for _, h := range handlers {
			rns.DebugLogger.Printf("Running Handler %v", h)
			h(event)
		}
	})
}

// FILES registers a url and directory path to serve static files. The webserver
// will serve all static files in any directories under these paths. Executing
// this method enables the static file server flag.
func (rns *RouteNamespace) FILES(url string, path string) {
	if !Settings.EnableStaticFileServer {
		Settings.EnableStaticFileServer = true
	}

	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}

	rns.InfoLogger.Printf(logprefix+"Registering static file path; Path: `%s`; URL: `%s`;", path, url)

	Settings.staticDir[url] = path
}

// GET is a conveinence method for registering handlers
func (rns *RouteNamespace) GET(path string, handlers ...HandlerFunc) {
	rns.Handle("GET", path, handlers)
}

// POST is a conveinence method for registering handlers
func (rns *RouteNamespace) POST(path string, handlers ...HandlerFunc) {
	rns.Handle("POST", path, handlers)
}

// PUT is a conveinence method for registering handlers
func (rns *RouteNamespace) PUT(path string, handlers ...HandlerFunc) {
	rns.Handle("PUT", path, handlers)
}
