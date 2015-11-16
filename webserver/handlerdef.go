package webserver

import (
	"strings"

	"github.com/go-gia/go-infrastructure/webserver/context"
)

type (
	// HandlerDef provides for a system to organize HandlerFunc metadata. Use of a
	// HandlerDef to describe a HandlerFunc is not required but provides a way
	// to easily configure advanced behavior and document that behavior.
	//
	// This abstraction binds documentation to implementation. This tight
	// coupling between the two helps reduce documentation burdens and
	// ensures documentation is kept current.
	HandlerDef struct {
		// A string to name the handler
		Alias string
		// The method to interact with the handler (i.e. GET or POST)
		Method string
		// The URL Path to access the handler
		Path string
		// The location of a HTML file describing the HandlerDef behavior in detail.
		Documentation         string
		DocumentationMarkdown string
		// The maximum time this HandlerFunc should take to process. This information is useful for performance testing.
		DurationExpectation string
		// An optional structure containing input parameters for the HandlerFunc.
		Params        interface{}
		ParamsExample interface{}
		// An optional reference to a structure containing output for successful HandlerFunc calls.
		Response        interface{}
		ResponseExample interface{}
		// An optional reference to a map describing response headers expected from the HandlerFunc.
		ResponseHeaders map[string]string
		// An optional reference to a map describing required request headers of the HandlerFunc.
		RequestHeaders map[string]string
		// The handler to register
		Handler HandlerFunc
		// A chain of handlers to process before executing the primary HandlerFunc
		PreHandlers []HandlerDef
		// A chain of handlers to process after executing the primary HandlerFunc
		PostHandlers []HandlerDef
	}
)

// SetHandlerFunc returns a copy of the provided HandlerDef, with the provided
// HandlerFunc set and is helpful when the HandlerFun would like to reference
// it's HandlerDef. Without setting the HandlerFunc into a copy applications
// are unable to compile due to a initialization loop.
func SetHandlerFunc(def HandlerDef, f HandlerFunc) HandlerDef {
	def.Handler = f
	return def
}

// RegisterHandlerDefs accepts a slice of HandlerDefs and registers each
func (s *Server) RegisterHandlerDefs(h []HandlerDef) error {
	for _, hd := range h {
		s.RegisterHandlerDef(hd)
	}

	return nil
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
	case "HEAD":
		fallthrough
	case "GET":
		fallthrough
	case "PUT":
		fallthrough
	case "DELETE":
		fallthrough
	case "OPTIONS":
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

type optionsMetadata struct {
	Get            bool
	Put            bool
	Post           bool
	Delete         bool
	Head           bool
	RequestHeaders map[string]string
}

// RegisterHandlerDefsAndOptions accepts a slice of HandlerDefs and registers
// each unique route and then after all the routes have been determined, creates
// new HandlerDefs with the OPTIONS method for each unique route.
func (s *Server) RegisterHandlerDefsAndOptions(h []HandlerDef) error {
	optionsMap := map[string]optionsMetadata{}
	defaultHeaders := make(map[string]map[string]string)
	// Let's loop through all the HandlerDefs and get collect methods / paths
	for _, hd := range h {
		if _, pathExists := optionsMap[hd.Path]; !pathExists {
			optionsMap[hd.Path] = optionsMetadata{}
			defaultHeaders[hd.Path] = hd.RequestHeaders
		}

		// Open up the current route
		o := optionsMap[hd.Path]
		o.Get = strings.ToUpper(hd.Method) == "GET"
		o.Post = strings.ToUpper(hd.Method) == "POST"
		o.Put = strings.ToUpper(hd.Method) == "PUT"
		o.Delete = strings.ToUpper(hd.Method) == "DELETE"
		o.Head = strings.ToUpper(hd.Method) == "HEAD"

		if len(hd.RequestHeaders) != len(defaultHeaders[hd.Path]) {
			return ErrWebserverRequestHeaderCountWrong
		}

		if len(hd.RequestHeaders) > 0 {
			for key, value := range defaultHeaders[hd.Path] {
				if _, ok := hd.RequestHeaders[key]; ok {
					if hd.RequestHeaders[key] != value {
						return ErrWebserverRequestHeaderMismatch
					}
				}
			}

			o.RequestHeaders = hd.RequestHeaders
		}

		optionsMap[hd.Path] = o
	}

	// // Now let's add to the end of the incoming HandlerDefs
	for route, meta := range optionsMap {
		h = append(h, createOption(route, meta))
	}
	// // Now, let's register everything.
	for _, hd := range h {
		s.RegisterHandlerDef(hd)
	}

	return nil
}

func createOption(path string, meta optionsMetadata) HandlerDef {
	methods := []string{}
	if meta.Get {
		methods = append(methods, "GET")
	}
	if meta.Put {
		methods = append(methods, "PUT")
	}
	if meta.Post {
		methods = append(methods, "POST")
	}
	if meta.Delete {
		methods = append(methods, "DELETE")
	}
	if meta.Head {
		methods = append(methods, "HEAD")
	}

	return HandlerDef{
		Method: "OPTIONS",
		Path:   path,
		Handler: func(c *context.Context) {
			c.Output.Header("Access-Control-Allow-Methods", strings.Join(methods, ","))
			for header, value := range meta.RequestHeaders {
				c.Output.Header(header, value)
			}
			c.Output.Body([]byte{})
		},
	}
}
