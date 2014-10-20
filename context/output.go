package context

import "strconv"

// Output represents the response written to a client
type Output struct {
	Status      int
	ContentType string
	Event       *Event
}

// NewOutput returns a new Output
func NewOutput(event *Event) *Output {
	return &Output{
		Event: event,
	}
}

// Header writes a response header to the client
func (output *Output) Header(key string, value string) {
	output.Event.ResponseWriter.Header().Set(key, value)
}

// Body writes the response to the client
func (output *Output) Body(content []byte) {
	output.Header("Content-Length", strconv.Itoa(len(content)))
	output.Event.ResponseWriter.Write(content)
}
