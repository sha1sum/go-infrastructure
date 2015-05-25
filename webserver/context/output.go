package context

import (
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"io"
	"strconv"
)

// Output represents the response written to a client
type Output struct {
	Status      int
	ContentType string
	Context     *Context
}

// NewOutput returns a new Output
func NewOutput(c *Context) *Output {
	output := &Output{
		Context: c,
	}
	// 200 OK by default
	output.Status = 200

	return output
}

// Body sets response body content and writes it to the client.
func (output *Output) Body(content []byte) {
	writer := output.Context.ResponseWriter.(io.Writer)

	output.Context.ResponseWriter.WriteHeader(output.Status)
	output.Header("Content-Length", strconv.Itoa(len(content)))

	writer.Write(content)

	switch writer.(type) {
	case *gzip.Writer:
		writer.(*gzip.Writer).Close()
	case *flate.Writer:
		writer.(*flate.Writer).Close()
	}
}

// Header writes a response header to the client
func (output *Output) Header(key string, value string) {
	output.Context.ResponseWriter.Header().Set(key, value)
}

// JSON is a conveinence method for writing JSON to the response body
func (output *Output) JSON(data interface{}, indent bool) error {
	output.Header("Content-Type", "application/json;charset=UTF-8")

	var content []byte
	var err error

	if indent {
		content, err = json.MarshalIndent(data, "", "  ")
	} else {
		content, err = json.Marshal(data)
	}
	if err != nil {
		panic(err)
	}

	output.Body(content)
	return nil
}
