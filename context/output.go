package context

import "strconv"

// Output represents the response written to a client
type Output struct {
	Status      int
	ContentType string
	Context     *Context
}

// NewOutput returns a new Output
func NewOutput(c *Context) *Output {
	return &Output{
		Context: c,
	}
}

// Header writes a response header to the client
func (output *Output) Header(key string, value string) {
	output.Context.ResponseWriter.Header().Set(key, value)
}

// Body writes the response to the client
func (output *Output) Body(content []byte) {
	output.Header("Content-Length", strconv.Itoa(len(content)))
	output.Context.ResponseWriter.Write(content)
}
