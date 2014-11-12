package webserver

import (
	"net/http"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

type (
	// RouteNamespace groups routes according to a specific URL entry point or prefix.
	RouteNamespace struct {
		prefix string
		parent *RouteNamespace
		server *Server
	}
)

func (rns *RouteNamespace) buildPath(p string) string {
	return path.Join(rns.prefix, p)
}

// Handle registers handlers!
func (rns *RouteNamespace) Handle(method string, path string, handlers []HandlerFunc) {
	p := rns.buildPath(path)

	if Settings.LogDebugMessages {
		log.WithFields(log.Fields{"event": packagename + "Handle", "method": method, "path": p}).Debug("Registering handler")
	}

	// Serve the request
	rns.server.router.Handle(method, path, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		if Settings.LogDebugMessages {
			log.WithFields(log.Fields{"event": packagename + "Handle", "method": method, "path": path}).Info("Capturing request")
		}
		event := rns.server.captureRequest(w, req, nil, handlers)

		// TODO Execute all handlers passed in their order stopping if one
		// chooses to write to the body unless we can/should simply append
		// the event body.
		handlers[0](event)

		// Write the response to the client
		if event.StatusCode == 0 {
			//w.WriteHeader(event.StatusCode)
		}
	})
}

// GET is a conveinence method for registering handlers
func (rns *RouteNamespace) GET(path string, handlers ...HandlerFunc) {
	if Settings.LogDebugMessages {
		log.Printf("Registering GET: %s", path)
	}

	rns.Handle("GET", path, handlers)
}

// FILES creates a
func (rns *RouteNamespace) FILES(url string, path string) {

	if !Settings.EnableStaticFileServer {
		Settings.EnableStaticFileServer = true
	}

	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}

	if Settings.LogDebugMessages {
		log.Printf("Registering static route %s -> %s", url, path)
	}

	Settings.staticDir[url] = path
}

// POST is a conveinence method for registering handlers
func (rns *RouteNamespace) POST(path string, handlers ...HandlerFunc) {
	if Settings.LogDebugMessages {
		log.Printf("Registering POST: %s", path)
	}

	rns.Handle("POST", path, handlers)
}
