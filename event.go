package webserver

import (
	"time"

	"github.com/sdming/gosnow"
)

// Snowflake ID generator
var snowflake, _ = gosnow.Default()

// Event is created on every request and models the event, or request.
// The webserver will populate a RequestContext with any data provided by the
// client from a form, URL, or recognized data type sent in the request body.
type Event struct {
	// ID is generated for each new RequestContext.
	id uint64 `json:"requestID"`
	// RequestContentLength contains a count of incoming bytes.
	RequestContentLength int `json:"requestContentLength"`
	// ResponseContentLength contains a count of outgoing bytes.
	ResponseContentLength int `json:"requestContentLength"`
	// StartTime can be used for performance metrics.
	StartTime time.Time `json:"startTime"`
	// Duration is set when a request is concluded and is a measure of how
	// long a request has taken.
	Duration time.Time `json:"duration"`
	// StatusCode
	StatusCode int
}

func (e Event) getID() uint64 {
	return e.id
}

// SetView declares an HTML view to render.
func (e *Event) SetView(view string, noLayout bool, args interface{}) {

}
