// Package webserver is responsible for web server operations including creating
// new web servers, registering handlers, and rendering and caching templates.
// One might think of this package as our own web framework that uses conventions
// to consistently work across products and projects.
package webserver

import (
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/aarongreenlee/webserver/context"
	"github.com/aarongreenlee/webserver/render"
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
	logprefix = "Package: webserver; Message: "
	// defaultResponse404 is returned if the server is unable to render the response
	// using the configured SystemTemplate. This can happen if a template file does not
	// exist at the configured path.
	defaultResponse404 = `<html><head><title>404 Not Found</title><style>body{background-color:black;color:white;margin:20%;}</style></head><body><center><h1>404 Not Found</h1></center></body></html>`
)

type (
	// Server represents an instance of the webserver.
	Server struct {
		*RouteNamespace
		contextPool    sync.Pool          // Manage our RequestContext event pool
		router         *httprouter.Router // Delegate to the httprouter package for performant route matching
		MissingHandler []HandlerFunc

		// HandlerDef maintains a map of all registered handler definitions
		HandlerDef      map[string]HandlerDef
		handlerDefMutex sync.Mutex

		DebugLogger   *log.Logger
		InfoLogger    *log.Logger
		WarningLogger *log.Logger
		ErrorLogger   *log.Logger
	}

	// Conventions defines our configuration.
	Conventions struct {
		// Reference to the conventions of the webserver's rendering engine
		Render *render.Conventions
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
		// Flag requests that take longer than N miliseconds. Default is 250ms (1/4th a second)
		RequestDurationWarning time.Duration
	}

	// HandlerFunc is a request event handler and accepts a RequestContext
	HandlerFunc func(*context.Context)

	// HandlerDef provides for a system to organize HandlerFunc metadata. Use of a
	// HandlerDef to describe a HandlerFunc is not required but provides a way
	// to eaisly configure advanced behavior and document that behavior.
	//
	// This abstraction binds documentation to implementation. This tight
	// coupeling between the two helps reduce documentaton buridens and
	// ensures documentation is kept current.
	HandlerDef struct {
		// A string to name the handler
		Alias string
		// The method to interact with the handler (i.e. GET or POST)
		Method string
		// The URL Path to access the handler
		Path string
		// The location of a HTML file describing the HandlerDef behavior in detail.
		Documentation string
		// The maximum time this HandlerFunc should take to process. This information is useful for performance testing.
		DurationExpectation string
		// An optional reference to a structure containing input paramaters for the HandlerFunc.
		Params interface{}
		// The handler to register
		Handler HandlerFunc
		// A chain of handlers to process before executing the primary HandlerFunc
		PreHandlers []HandlerDef
		// A chain of handlers to process after executing the primary HandlerFunc
		PostHandlers []HandlerDef
	}
)

var (
	// Settings allows a developer to override the conventional settings of the
	// webserver.
	Settings = Conventions{
		Render:                 &render.Settings,
		EnableStaticFileServer: false,
		SystemTemplates: map[string]string{
			"onMissingHandler": "errors/onMissingHandler",
		},
		staticDir:              make(map[string]string),
		RequestDurationWarning: time.Second / 4,
	}
	// If we fail to find a configured onMissingHandler once we will stop looking
	seekOnMissingHandler = true
)

// New returns a new WebServer.
func New(
	debugLogger *log.Logger,
	infoLogger *log.Logger,
	warningLogger *log.Logger,
	errorLogger *log.Logger) *Server {

	s := &Server{
		DebugLogger:   debugLogger,
		InfoLogger:    infoLogger,
		WarningLogger: warningLogger,
		ErrorLogger:   errorLogger,

		HandlerDef: make(map[string]HandlerDef),
	}
	// Setup an initial route namespace
	s.RouteNamespace = &RouteNamespace{
		prefix:        "/",
		server:        s,
		DebugLogger:   debugLogger,
		InfoLogger:    infoLogger,
		WarningLogger: warningLogger,
		ErrorLogger:   errorLogger}

	s.router = httprouter.New()
	s.router.NotFound = s.onMissingHandler

	return s
}

// RegisterHandlerDefs accepts a slice of HandlerDefs and registers each
func (s *Server) RegisterHandlerDefs(h []HandlerDef) {
	for _, hd := range h {
		s.RegisterHandlerDef(hd)
	}
}

// RegisterHandlerDef accepts a HandlerDef and registers it's behavior with the
// webserver.
func (s *Server) RegisterHandlerDef(h HandlerDef) {
	chain := []HandlerFunc{}

	// Pre
	for _, a := range h.PreHandlers {
		chain = append(chain, a.Handler)
	}
	// Target
	chain = append(chain, h.Handler)
	// Post
	for _, a := range h.PostHandlers {
		chain = append(chain, a.Handler)
	}

	// Register
	switch h.Method {
	case "GET":
		fallthrough
	case "PUT":
		fallthrough
	case "DELETE":
		fallthrough
	case "POST":
		s.Handle(h.Method, h.Path, chain)

	case "":
		// do nothing--middleware only
	default:
		panic("Unable to register handler due to unknown method: " + h.Method)
	}

	// Register the handler def for auto-documentation
	s.handlerDefMutex.Lock()
	defer s.handlerDefMutex.Unlock()

	s.HandlerDef[h.Method+":"+h.Path] = h
}

// Start launches the webserver so that it begins listening and serving requests
// on the desired address.
func (s *Server) Start(address string) {
	s.InfoLogger.Printf(logprefix+"Starting webserver; Address: %s;", address)

	if err := http.ListenAndServe(address, s); err != nil {
		s.ErrorLogger.Printf(logprefix+"Unable to start; Address: %s; Error: %s", address, err.Error())
		panic(err)
	}
}

// ServeHTTP handles all requests of our web server
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	starttick := time.Now() // TODO Look at ticker? Inside time package.

	requestPath := req.URL.Path

	if Settings.EnableStaticFileServer {
		for prefix, staticDir := range Settings.staticDir {
			s.DebugLogger.Printf(logprefix+"Evaluating static route; Path: %s:%s;", req.Method, requestPath)
			if strings.HasPrefix(requestPath, prefix) {
				filePath := staticDir + requestPath[len(prefix):]
				fileInfo, err := os.Stat(filePath)
				if err != nil {
					s.WarningLogger.Printf(logprefix+"Static asset not found; Path: %s;", requestPath)
					s.onMissingHandler(w, req)
					return
				}
				// TODO: Serve Directory Listing? Throw a 403 Forbidden Error? Defaulting to 404 is probably not robust enough for our web server
				if fileInfo.IsDir() {
					s.onMissingHandler(w, req)
				}

				s.DebugLogger.Printf(logprefix+"Serving static asset; Path: %s;", requestPath)

				// TODO: Enable gZIP support if allowed for css, js, etc.
				http.ServeFile(w, req, filePath)
				return
			}
		}
	}

	s.router.ServeHTTP(w, req)

	duration := time.Since(starttick)
	if duration >= Settings.RequestDurationWarning {
		s.WarningLogger.Printf(logprefix+"Request complete; Path: %s; Duration: %fs", requestPath, duration.Seconds())
	} else {
		s.DebugLogger.Printf(logprefix+"Request complete; Path: %s; Duration: %fs", requestPath, duration.Seconds())
	}
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
	w.WriteHeader(event.StatusCode)

	s.DebugLogger.Printf(logprefix+"Executing onMissingHandler; Address: %s;", req.URL.Path)

	if seekOnMissingHandler {
		template := Settings.SystemTemplates["onMissingHandler"]
		err := event.HTMLTemplate(template, nil)
		if err != nil {
			s.ErrorLogger.Printf(logprefix+"Failed single attempt to load configured onMissingHandler template--serving default response; Path: %s;", template)
			seekOnMissingHandler = false
		}
	}

	if !seekOnMissingHandler {
		event.Output.Body([]byte(defaultResponse404))
	}
}
