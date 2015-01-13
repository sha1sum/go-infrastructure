package context

import (
	"net/http"

	"github.com/aarongreenlee/webserver/render"

	"github.com/davecgh/go-spew/spew"
	"github.com/julienschmidt/httprouter"
)

// Context is created on every request and models the event, or request.
// The webserver will populate a RequestContext with any data provided by the
// client from a form, URL, or recognized data type sent in the request body.
type Context struct {
	// RequestContentLength contains a count of incoming bytes.
	RequestContentLength int
	// ResponseContentLength contains a count of outgoing bytes.
	ResponseContentLength int

	renderer render.Renderer

	Input          *Input
	Output         *Output
	Request        *http.Request
	ResponseWriter http.ResponseWriter
    Params   httprouter.Params
}

// New produces a new request context event.
func New(w http.ResponseWriter, req *http.Request, params httprouter.Params) *Context {
	var c = new(Context)

	c.Input = NewInput(req)
	c.Output = NewOutput(c)
	c.Request = req
	c.ResponseWriter = w
    c.Params = params

	if c.Input.Is("POST") || c.Input.Is("PUT") {
		c.Input.Body()
	}

	return c
}

// *****************************************************************************
// Handling Conveinence
// *****************************************************************************

// BadRequest issues a bad request
func (c *Context) BadRequest(output interface{}) {
	c.Output.Status = http.StatusBadRequest

	switch output.(type) {
	case error:
		c.Output.Body([]byte(output.(error).Error()))
	case []byte:
		c.Output.Body(output.([]byte))
	case string:
		c.Output.Body([]byte(output.(string)))
	default:
		c.Output.Body([]byte("Bad Request"))
	}
}

// InternalError issues a 500 Internal Server Errror
func (c *Context) InternalError(output interface{}) {
	c.Output.Status = http.StatusInternalServerError

	switch output.(type) {
	case error:
		c.Output.Body([]byte(output.(error).Error()))
	case []byte:
		c.Output.Body(output.([]byte))
	case string:
		c.Output.Body([]byte(output.(string)))
	default:
		c.Output.Body([]byte("Internal Server Error"))
	}

}

// *****************************************************************************
// Rendering
// *****************************************************************************

// HTML outputs the provided HTML string
func (c *Context) HTML(output string) error {
    
    c.Output.Header("Content-Type", "text/html")
	c.Output.Body([]byte(output))

	return nil
}

// HTMLTemplate renders the HTML view specified by it's filename omitting the file extension.
func (c *Context) HTMLTemplate(name string, args interface{}) error {
    
	content, err := render.HTML.Render(name, args)
	if err != nil {
		return err
	}

    c.Output.Header("Content-Type", "text/html")
	c.Output.Body(content)

	return nil
}

// Dump spews the provided value to the stdout and is useful for debugging.
func (c *Context) Dump(v interface{}) {
	spew.Dump(v)
}
