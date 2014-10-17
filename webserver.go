// Package webserver is responsible for web server operations including creating
// new web servers, registering handlers, and rendering and caching templates.
// One might think of this package as our own web framework that uses conventions
// to consistently work across products and projects.
package webserver

import "sync"

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

type (
	// Handler is a request event handler
	Handler func(*RequestContext)

	// RouteNamespace groups routes according to a specific URL entry point or prefix.
	RouteNamespace struct {
		Handlers []Handler
		prefix   string
		parent   *RouteNamespace
		server   *Server
	}

	// Server represents an instance of the webserver.
	Server struct {
		*RouteNamespace
		eventPool sync.Pool // Manage our RequestContext event pool
	}
)

// New returns a new WebServer
func New() *Server {
	ws := &Server{}
	// Setup an initial route namespace
	ws.RouteNamespace = &RouteNamespace{
		Handlers: nil,
		prefix:   "/",
		parent:   nil,
		server:   ws}

	return ws
}

/*
type (
	// Webserver is a functional webserver that a main program can use with the
	// standard http package to serve traffic
	webserver struct {
		Mux              *mux.Router //Gorrilla MUX to resolve runtime routing of requests
		HTTPHandlerCount int         // Count the handlers registered for our server
	}
)

func (ws *webserver) RegisterHandler(h *handler.Handler) {
	if h.SupportsHTTP() {
		log.Println("Registering handler: ", h.Name)
		ws.Mux.HandleFunc(h.HTTPPath, h.HTTPRunner).Methods(h.HTTPMethod)
		ws.HTTPHandlerCount++
	}
}

// NewWebserver is a factory method to construct a new webserver given the
// provided handlers. The new webserver will automatically be configured to
// server static assets from the /cdn/css; /cdn/js; and /cdn/m directories
func NewWebserver(handlers []*handler.Handler) *webserver {
	r := mux.NewRouter()
	ws := &webserver{Mux: r}

	for _, h := range handlers {
		ws.RegisterHandler(h)
	}

	// Register some static asset routes to serve files from the applications
	ws.Mux.PathPrefix("/cdn/").Handler(http.FileServer(http.Dir("./")))

	r.NotFoundHandler = http.HandlerFunc(notFound)

	return ws
}

// ServeHTTP sends our 404 response to the client
func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("FourOhFour! What have you done?!?!?"))
}

// ServeHTTP handles all requests of our web server and delegates to the
// gorilla mux package for routing and actual handler execution.
func (ws *webserver) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	re := NewRequestEvent(r)

	// When using the recorder there is something lost in translation for static
	// files with a "If Modified Since" header with the request--which browsers
	// tend to do for second+ requests of the static assets. As a quick fix we
	// will keep things simple and use a real response writer and avoid our
	// lifecycle for static assets.
	isCDN, pathError := path.Match("/cdn/ * / *", r.URL.Path)
	if pathError != nil {
		log.Fatal(pathError)
		isCDN = false
	}

	if isCDN {
		log.Println("[webserver] CDN simulation: ", r.URL.Path)
		ws.Mux.ServeHTTP(w, r)
		return
	}

	// Record the way the handlers treat ResponseWriter using the test package's
	// Recorder. This technique will prevent any headers from being written
	// before we can write them in this handler. This is important because as
	// soon as a header is written bytes are sent to the client and we loose our
	// chance to write headers in subsequent processing.
	rec := httptest.NewRecorder()

	// Provide the recorder to the Gorilla mux package to router and execute
	// the registered handler (if any)
	ws.Mux.ServeHTTP(rec, r)

	// Determine how many bytes we've written from the handler system
	re.ResponseContentLength = len(rec.Body.Bytes())

	// Actually write headers to the ResponseWriter by copying any handlers
	// set by our executed handlers into the ResponseWriter
	for k, v := range rec.Header() {
		w.Header()[k] = v
	}

	// Write the content length as measured by our
	w.Header().Set("Content-Length", strconv.Itoa(re.ResponseContentLength))
	w.Header().Set("VND.wreckerlabs.com.runtime", strconv.FormatFloat(re.GetCurrentRuntime().Seconds(), 'f', 6, 64))

	// Set the status code--which also starts sending bytes back to the client
	// and prevnts us from sending any more headers
	w.WriteHeader(200)

	// Write the original body provided by the handler to the ResponseWriter
	w.Write(rec.Body.Bytes())
}
*/
