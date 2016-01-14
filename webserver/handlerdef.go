package webserver

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-gia/go-infrastructure/webserver/context"
	openapi "github.com/sha1sum/golang-openapi"
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
		Alias string `json:"alias,omitempty"`
		// The method to interact with the handler (i.e. GET or POST)
		Method string `json:"method,omitempty"`
		// The URL Path to access the handler
		Path string `json:"path,omitempty"`
		// The location of a HTML file describing the HandlerDef behavior in detail.
		Documentation string `json:"documentationURL,omitempty"`
		// DocumentationMarkdown defines a description for the handler and takes markdown format
		DocumentationMarkdown string `json:"documentationMarkdown,omitempty"`
		// The maximum time this HandlerFunc should take to process. This information is useful for performance testing.
		DurationExpectation string `json:"duration,omitempty"`
		// OpenAPIParams is a list of parameters which the handler can accept
		OpenAPIParams []openapi.Parameter `json:"openAPIParams,omitempty"`
		// OpenAPIResponses is a map of http error codes as strings with their Response
		OpenAPIResponses map[string]openapi.Response `json:"openAPIResponses,omitempty"`
		// An optional structure containing search/query parameters for the
		// HandlerFunc.
		Params        interface{} `json:"params,omitempty"`
		ParamsExample interface{} `json:"paramsExample,omitempty"`
		// An optional structure documenting the request body supported.
		RequestBody        interface{} `json:"request,omitempty"`
		RequestBodyExample interface{} `json:"requestExample,omitempty"`
		// An optional structure documenting the response body.
		ResponseBody        interface{} `json:"-"`
		ResponseBodyExample interface{} `json:"responseExample,omitempty"`
		// An optional reference to a map describing response headers expected from the HandlerFunc.
		ResponseHeaders map[string]string `json:"responseHeaders,omitempty"`
		// An optional reference to a map describing required request headers of the HandlerFunc.
		RequestHeaders map[string]string `json:"requestHeaders,omitempty"`
		// The handler to register
		Handler HandlerFunc `json:"-"`
		// A chain of handlers to process before executing the primary HandlerFunc
		PreHandlers []HandlerDef `json:"prehandlers,omitempty"`
		// A chain of handlers to process after executing the primary HandlerFunc
		PostHandlers []HandlerDef `json:"posthandlers,omitempty"`
		// Summary is a short title for the handler
		Summary string `json:"summary,omitempty"`
		// Tags is a list of tags used for API documentation
		Tags []string `json:"tags,omitempty"`
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
	postChain := []HandlerFunc{}

	// Pre
	for _, a := range h.PreHandlers {
		chain = append(chain, a.Handler)
	}
	// Target
	chain = append(chain, h.Handler)

	for _, a := range h.PostHandlers {
		postChain = append(postChain, a.Handler)
	}

	// Register
	switch h.Method {
	case HEAD:
		fallthrough
	case GET:
		fallthrough
	case PUT:
		fallthrough
	case DELETE:
		fallthrough
	case OPTIONS:
		fallthrough
	case PATCH:
		fallthrough
	case POST:
		s.Handle(h.Method, h.Path, chain, postChain)

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
	get            bool
	put            bool
	post           bool
	delete         bool
	head           bool
	patch          bool
	RequestHeaders map[string]string
}

// RegisterHandlerDefsAndOptions accepts a slice of HandlerDefs and registers
// each unique route. After all the routes have been determined it then creates
// new HandlerDefs to create OPTIONS methods for each unique route. The created
// OPTIONS handlers read the `HandlerDef.RequestHeaders` and tell the client
// that they
func (s *Server) RegisterHandlerDefsAndOptions(h []HandlerDef) error {

	optionsMap := map[string]optionsMetadata{}
	headers := make(map[string]map[string]string)

	// Loop through all the HandlerDefs and collect methods / paths.
	for _, hd := range h {
		if _, pathExists := optionsMap[hd.Path]; !pathExists {
			optionsMap[hd.Path] = optionsMetadata{}
			headers[hd.Path] = hd.RequestHeaders
		}

		// Open up the current route
		o := optionsMap[hd.Path]
		switch strings.ToUpper(hd.Method) {
		case GET:
			o.get = true
		case POST:
			o.post = true
		case PUT:
			o.put = true
		case DELETE:
			o.delete = true
		case HEAD:
			o.head = true
		case PATCH:
			o.patch = true
		}

		if len(hd.RequestHeaders) != len(headers[hd.Path]) {
			return ErrWebserverRequestHeaderCountWrong
		}

		if len(hd.RequestHeaders) > 0 {
			for key, value := range headers[hd.Path] {
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

	j, _ := json.Marshal(optionsMap)
	fmt.Printf("The JSON is \n\n%s\n\n", j)

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
	if meta.get {
		methods = append(methods, GET)
	}
	if meta.put {
		methods = append(methods, PUT)
	}
	if meta.post {
		methods = append(methods, POST)
	}
	if meta.delete {
		methods = append(methods, DELETE)
	}
	if meta.head {
		methods = append(methods, HEAD)
	}
	if meta.patch {
		methods = append(methods, PATCH)
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

// ToOpenAPI returns an string with the path, HTTP verb, and OpenAPI request
func (h HandlerDef) ToOpenAPI() (path string, verb string, req openapi.Request) {
	return h.Path, strings.ToLower(h.Method), openapi.Request{
		Summary:     h.Summary,
		Description: h.DocumentationMarkdown,
		Parameters:  h.OpenAPIParams,
		Tags:        h.Tags,
		Responses:   h.OpenAPIResponses,
	}
}
