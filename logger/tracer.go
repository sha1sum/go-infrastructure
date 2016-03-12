package logger

import "fmt"

const (
  escape = "\x1b"
  red    = "1;31"
  green  = "1;32"
  blue   = "1;34"
)

type (
// Tracer accepts low-level debugging information and its implementations
// decide what to do with the messages it receives
  Tracer interface {
    Trace(title string, messages []interface{})
  }

// NilTracer implements the Tracer interface and is used for discarding traces.
// This implementation is usually used in production to discard debug info.
  NilTracer struct{}

// StdOutTracer implements the Tracer interface and is used for outputting
// trace information to stdout (standard output). This is useful when needing
// verbose debug information.
  StdOutTracer struct{}
)

// Trace for NilTracer will simply do nothing and return.
func (NilTracer) Trace(string, []interface{}) {
  return
}

// Trace for StdOutTracer will output a number of messages to stdout, attempting
// to use color for titles and indentation for messages.
func (t StdOutTracer) Trace(title string, messages []interface{}) {
  fmt.Printf("%s[%smTRCE%s[0m\n\t%s[%sm%s:%s[0m\n\n", escape, red, escape, escape, blue, title, escape)
  for i, m := range messages {
    fmt.Printf("\t\t%s[%sm* %s[0m%+v\n", escape, green, escape, m)
    if i+1 == len(messages) {
      fmt.Printf("\n")
    }
  }
}
