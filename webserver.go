// Package webserver is responsible for web server operations including creating
// new web servers, registering handlers, and rendering and caching templates.
// One might think of this package as our own web framework that uses conventions
// to consistently work across products and projects.
package webserver

import (
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"git.wreckerlabs.com/in/webserver/context"
	"git.wreckerlabs.com/in/webserver/render"
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

const (
	// MIMEJSON represents the standard classification for data encoded as JSON.
	MIMEJSON = "application/json"
	// MIMEHTML represents the standard classification for HTML web pages.
	MIMEHTML = "text/html"
	// MIMEXML represents the standard classification for data encoded as XML.
	MIMEXML = "application/xml"
	// MIMEXMLTEXT represents the standard classification for a XML text document.
	MIMEXMLTEXT = "text/xml"
	// MIMEPLAIN represents the standard classification for plain text data without
	// any encoded format and is generally human readable text data.
	MIMEPLAIN = "text/plain"
	// MIMEFORM represents form data encoded by a Web browser posted to a server.
	MIMEFORM = "application/x-www-form-urlencoded"
	// MIMECSS represents the standard classificaton for Cascading Style Sheets.
	MIMECSS = "text/css"
	// MIMEJS represents the standard classification for JavaScript.
	MIMEJS = "application/javascript"
	// MIMEPNG represents the standard classificaton for PNG images.
	MIMEPNG = "image/png"
	// MIMEJPEG represents the standard classificaton for JPEG/JPG images.
	MIMEJPEG = "image/jpeg"
	// MIMEGIF represents the standard classificaton for GIF images.
	MIMEGIF = "image/gif"
	// MIMEXICON represents the proposed classification for icons such as favicon images
	MIMEXICON = "image/x-icon"
)

const (
	packagename = "webserver:"
	// defaultResponse404 is returned if the server is unable to render the response
	// using the configured SystemTemplate. This can happen if a template file does not
	// exist at the configured path.
	defaultResponse404 = `<html><head><title>404 Not Found</title><style>body{background-color:black;color:white;margin:20%;}</style></head><body><center><h1>404 Not Found</h1><hr><p>WreckerLabs Webserver</p></center></body></html>`
)

type (
	// HandlerFunc is a request event handler and accepts a RequestContext
	HandlerFunc func(*context.Context)

	// Server represents an instance of the webserver.
	Server struct {
		*RouteNamespace
		contextPool    sync.Pool          // Manage our RequestContext event pool
		router         *httprouter.Router // Delegate to the httprouter package for performant route matching
		MissingHandler []HandlerFunc
	}

	// Conventions define configuration and are set to our conventional, default
	// values.
	Conventions struct {
		// Reference to the conventions of the webserver's rendering engine
		Render *render.Conventions
		// LogDebugMessages if true, enables debug observation logging
		LogDebugMessages bool
		// LogWarningMessages if true, enables warning logging behavior.
		LogWarningMessages bool
		// LogErrorMessages if true, enables error logging behavior
		LogErrorMessages bool
		// EnableStaticFileServer if true, enables the serving of static assets such as CSS, JS, or other files.
		EnableStaticFileServer bool
		// StaticFilePath defines the releative root directory static files can be served from. Example "public" or "web-src/cdn/"
		StaticFilePath string
		// SystemTemplates is a map of keys to template paths. Default templates
		// such as `onMissingHandler` (404) are configurable here allowing developers
		// to customize exception pages for each implementation.
		SystemTemplates map[string]string
		// A map of directory paths the webserver should serve static files from
		staticDir map[string]string
	}
)

var (
	// Settings allows a developer to override the conventional settings of the
	// webserver.
	Settings = Conventions{
		Render:                 &render.Settings,
		LogDebugMessages:       false,
		LogWarningMessages:     false,
		LogErrorMessages:       true,
		EnableStaticFileServer: false,
		SystemTemplates: map[string]string{
			"onMissingHandler": "errors/onMissingHandler",
		},
		staticDir: make(map[string]string),
	}
	// If we fail to find a configured onMissingHandler once we will stop looking
	seekOnMissingHandler = true
)

// New returns a new WebServer.
func New() *Server {
	log.SetLevel(log.DebugLevel)
	s := &Server{}
	// Setup an initial route namespace
	s.RouteNamespace = &RouteNamespace{
		prefix: "/",
		parent: nil,
		server: s}

	s.router = httprouter.New()
	s.router.NotFound = s.onMissingHandler

	return s
}

// Start launches the webserver so that it begins listening and serving requests
// on the desired address.
func (s *Server) Start(address string) {
	log.WithFields(log.Fields{"event": packagename + "Start", "topic": "ListenAndServe", "key": address}).Info("Starting")

	if err := http.ListenAndServe(address, s); err != nil {
		log.WithFields(log.Fields{"event": "Start", "topic": "ListenAndServe", "key": address}).Fatal("Unable to start webserver")
		panic(err)
	}
}

// ServeHTTP handles all requests of our web server
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	starttick := time.Now()

	requestPath := req.URL.Path

	if Settings.EnableStaticFileServer {
		for prefix, staticDir := range Settings.staticDir {
			if Settings.LogDebugMessages {
				log.WithFields(log.Fields{"event": packagename + "ServeHTTP", "topic": "Serve Static", "key": requestPath}).Debug("Evaluating static route")
			}
			if strings.HasPrefix(requestPath, prefix) {
				filePath := staticDir + requestPath[len(prefix):]
				fileInfo, err := os.Stat(filePath)
				if err != nil {
					if Settings.LogWarningMessages {
						log.WithFields(log.Fields{"event": packagename + "ServeHTTP", "topic": "Serve Static", "key": requestPath}).Warn("Unable to load file")
					}
					s.onMissingHandler(w, req)
					return
				}
				// TODO: Serve Directory Listing? Throw a 403 Forbidden Error? Defaulting to 404 is probably not robust enough for our web server
				if fileInfo.IsDir() {
					s.onMissingHandler(w, req)
				}

				if Settings.LogDebugMessages {
					log.WithFields(log.Fields{"event": packagename + "ServeHTTP", "topic": "Serve Static", "key": requestPath}).Debug("Serving file")
				}

				// TODO: Enable gZIP support if allowed for css, js, etc.
				http.ServeFile(w, req, filePath)
				return
			}
		}
	}

	s.router.ServeHTTP(w, req)

	log.WithFields(log.Fields{"event": packagename + "ServeHTTP", "ms": time.Since(starttick) * time.Millisecond, "path": requestPath}).Info("Request complete")
}

// captureRequest builds a new Event to model a request/response handled
// by our Webserver.
func (s *Server) captureRequest(
	w http.ResponseWriter,
	req *http.Request,
	params httprouter.Params,
	handlers []HandlerFunc) *context.Context {

	event := context.New(w, req, params)

	return event
}

// onMissingHandler replies to the request with an HTTP 404 not found error.
// This function is triggered when we are unable to match a route.
func (s *Server) onMissingHandler(w http.ResponseWriter, req *http.Request) {
	event := s.captureRequest(w, req, nil, s.MissingHandler)
	event.StatusCode = http.StatusNotFound

	if seekOnMissingHandler {
		template := Settings.SystemTemplates["onMissingHandler"]
		err := event.HTML(template, nil)
		if err != nil {
			log.WithFields(log.Fields{"event": packagename + "onMissingHandler", "topic": "Configured Template", "key": template}).Error("Failed to render template")
			seekOnMissingHandler = false
		}
	}

	if !seekOnMissingHandler {
		event.Output.Body([]byte(defaultResponse404))
	}
}
