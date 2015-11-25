// Package webserver is responsible for web server operations including creating
// new web servers, registering handlers, and rendering and caching templates.
// One might think of this package as our own web framework that uses conventions
// to consistently work across products and projects.
package webserver

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-gia/go-infrastructure/logger"
	"github.com/go-gia/go-infrastructure/webserver/context"
	"github.com/go-gia/go-infrastructure/webserver/render"
	"github.com/gorilla/mux"
)

const (
	// GET http method reference
	GET = "GET"
	// POST http method reference
	POST = "POST"
	// OPTIONS http method reference
	OPTIONS = "OPTIONS"
	// DELETE http method reference
	DELETE = "DELETE"
	// PUT http method reference
	PUT = "PUT"
	// PATCH http method reference
	PATCH = "PATCH"
	// HEAD http method reference
	HEAD = "HEAD"
)

// defaultResponse404 is returned if the server is unable to render the response
// using the configured SystemTemplate. This can happen if a template file does not
// exist at the configured path.
const defaultResponse404 = `<html><head><title>404 Not Found</title><style>body{background-color:black;color:white;margin:20%;}</style></head><body><center><h1>404 Not Found</h1></center></body></html>`

type (
	// Server represents an instance of the webserver.
	Server struct {
		contextPool   sync.Pool
		methodRouters map[string]*mux.Router

		MissingHandler []HandlerFunc

		// HandlerDef maintains a map of all registered handler definitions
		HandlerDef      map[string]HandlerDef
		handlerDefMutex sync.Mutex

		logger logger.Logger
	}

	// Conventions defines our configuration.
	Conventions struct {
		// Reference to the conventions of the webserver's rendering engine
		Render *render.Conventions
		// EnableStaticFileServer if true, enables the serving of static assets such as CSS, JS, or other files.
		EnableStaticFileServer bool
		// StaticFilePath defines the relative root directory static files can be served from. Example "public" or "web-src/cdn/"
		StaticFilePath string
		// SystemTemplates is a map of keys to template paths. Default templates
		// such as `onMissingHandler` (404) are configurable here allowing developers
		// to customize exception pages for each implementation.
		SystemTemplates map[string]string
		// A map of directory paths the webserver should serve static files from
		staticDir map[string]string
		// Flag requests that take longer than N milliseconds. Default is 250ms (1/4th a second)
		RequestDurationWarning time.Duration
	}

	// HandlerFunc is a request event handler and accepts a RequestContext
	HandlerFunc func(*context.Context)
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

	// ErrWebserverDuplicateMethod is thrown when there's a route that has duplicate methods (read: Two PUT requests on the same route)
	ErrWebserverDuplicateMethod = errors.New("Duplicate Method on a route.")
	// ErrWebserverRequestHeaderCountWrong is thrown when the request header counts are wrong between similar routes.
	ErrWebserverRequestHeaderCountWrong = errors.New("The current route doesn't have the same number of RequestHeaders as a previous route header.")
	// ErrWebserverRequestHeaderMismatch is thrown when the request header of a route doesn't match a request header of another route.
	ErrWebserverRequestHeaderMismatch = errors.New("The routes have RequestHeaders mismatch. For similiar routes, all the verbs should use the same RequestHeaders.")
)

// New returns a new WebServer.
func New(
	logger logger.Logger) *Server {

	s := &Server{
		logger:        logger,
		HandlerDef:    make(map[string]HandlerDef),
		methodRouters: make(map[string]*mux.Router),
	}

	// Be sure to setup at least one router. Additional method routers
	// can be defined when HandlerFuncs are registered.
	s.methodRouters[GET] = mux.NewRouter()

	// TODO We need to reset a default missing handler
	// s.router.NotFound = s.onMissingHandler

	return s
}

// Start launches the webserver so that it begins listening and serving requests
// on the desired address.
func (s *Server) Start(address string) {
	if err := http.ListenAndServe(address, s); err != nil {
		panic(err)
	}
}

// ServeHTTP handles all requests of our web server
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	starttick := time.Now() // TODO Look at ticker? Inside time package.

	requestPath := req.URL.Path

	s.logger.Context(logger.Fields{
		"requestPath": requestPath,
		"method":      req.Method,
	}).Debug("GO-GIA Webserver is receiving a request")

	if Settings.EnableStaticFileServer {
		for prefix, staticDir := range Settings.staticDir {
			s.logger.Context(logger.Fields{"method": req.Method, "requestPath": requestPath}).Debug("Evaluating static route")

			if strings.HasPrefix(requestPath, prefix) {
				filePath := staticDir + requestPath[len(prefix):]
				fileInfo, err := os.Stat(filePath)
				if err != nil {
					s.logger.Context(logger.Fields{"filepath": filePath, "requestPath": requestPath}).Warn("Static file not found")
					s.onMissingHandler(w, req)
					return
				}
				// TODO: Serve Directory Listing? Throw a 403 Forbidden Error? Defaulting to 404 is probably not robust enough for our web server
				if fileInfo.IsDir() {
					s.onMissingHandler(w, req)
				}

				s.logger.Context(logger.Fields{"filepath": filePath, "requestPath": requestPath}).Debug("Serving static file")

				// TODO: Enable gZIP support if allowed for css, js, etc.
				http.ServeFile(w, req, filePath)
				return
			}
		}
	}

	router, ok := s.methodRouters[req.Method]
	if !ok {
		router = s.methodRouters[GET]
	}

	router.ServeHTTP(w, req)

	duration := time.Since(starttick)
	if duration >= Settings.RequestDurationWarning {
		//s.WarningLogger.Printf(logprefix+"Request complete; Path: %s; Duration: %fs", requestPath, duration.Seconds())
	} else {
		//s.DebugLogger.Printf(logprefix+"Request complete; Path: %s; Duration: %fs", requestPath, duration.Seconds())
	}
}

// captureRequest builds a new Event to model a request/response handled
// by our Webserver.
func (s *Server) captureRequest(
	w http.ResponseWriter,
	req *http.Request,
	handlers []HandlerFunc) *context.Context {

	event := context.New(w, req)

	return event
}

// onMissingHandler replies to the request with an HTTP 404 not found error.
// This function is triggered when we are unable to match a route.
func (s *Server) onMissingHandler(w http.ResponseWriter, req *http.Request) {
	context := s.captureRequest(w, req, s.MissingHandler)

	context.Output.Status = http.StatusNotFound

	s.logger.Context(logger.Fields{"method": req.Method, "requestPath": req.URL.Path, "statusCode": 404}).Debug("Handler not found")

	if seekOnMissingHandler {
		template := Settings.SystemTemplates["onMissingHandler"]
		err := context.HTMLTemplate(template, nil)
		if err != nil {
			s.logger.Context(logger.Fields{"template": template}).Warn("Failed single attempt to load configured onMissingHandler template--serving default response")
			seekOnMissingHandler = false
		}
	}

	if !seekOnMissingHandler {
		context.Output.Body([]byte(defaultResponse404))
	}
}

// Handle registers HandlerFuncs with the webserver.
func (s *Server) Handle(method string, path string, handlers []HandlerFunc, postHandlers []HandlerFunc) {
	router, ok := s.methodRouters[method]
	if !ok {
		router = mux.NewRouter()
		s.methodRouters[method] = router
	}
	s.logger.Context(logger.Fields{"method": method, "path": path}).Debug("Registering Route")

	router.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		event := s.captureRequest(w, req, handlers)
		// Run through our handler chain
		for _, h := range handlers {
			if event.BreakHandlerChain {
				break
			}
			h(event)
		}

		// Run through any post handlers. These are not allowed to write
		// to the client.
		if postHandlers != nil {
			for _, h := range postHandlers {
				h(event)
			}
		}
	}).Methods(method)
}

// FILES registers a url and directory path to serve static files. The webserver
// will serve all static files in any directories under these paths. Executing
// this method enables the static file server flag.
func (s *Server) FILES(url string, path string) {
	if !Settings.EnableStaticFileServer {
		Settings.EnableStaticFileServer = true
	}

	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}

	Settings.staticDir[url] = path
}

// GET is a convenience method for registering handlers
func (s *Server) GET(path string, handlers ...HandlerFunc) {
	s.Handle("GET", path, handlers, nil)
}

// POST is a convenience method for registering handlers
func (s *Server) POST(path string, handlers ...HandlerFunc) {
	s.Handle("POST", path, handlers, nil)
}

// PUT is a convenience method for registering handlers
func (s *Server) PUT(path string, handlers ...HandlerFunc) {
	s.Handle("PUT", path, handlers, nil)
}
