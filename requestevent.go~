package webserver

import (
	"log"
	"net/http"
	"time"

	"github.com/sdming/gosnow" // ID generation inspired by twitter's snowflake technique
)

var snowflake, _ = gosnow.Default()

// RequestEvent represents a single command to the system. For example, 'create user', or 'crete users', or, 'find
// users where firstname like Aaron`' may all be individual commands sent to the system. Each of these requests
// to do something creates a RequestEvent.
//
// RequestEvents...
// ...measure incoming bytes
// ...measure outgoing bytes
// ...measure processing time
// ...support dispatching measurements to O.C.E.AN
type RequestEvent struct {
	// An ID to represent this MachineEvnet
	ID uint64 `json:"id"`
	// Incoming Bytes
	RequestContentLength int `json:"requestContentLength"`
	// Outgoing Bytes
	ResponseContentLength int `json:"responseContentLength"`
	// Startime of our RequestEvent
	StartTime time.Time `json:"startTime"`
	// The endtime of our RequestEvent--How much time did it take?
	Runtime time.Time `json:"runtime"`
	// Address Requested
	Address string `json:"address"`
	// True if this RequestEvent is HTTP
	IsHTTP bool `json:"isHTTP"`
	// True if this request should lookup data
	IsReadRequest bool `json:"isReadRequest"`
	// True if this request should mutate data
	IsWriteRequest bool `json:"isWriteRequest"`
	// True if this request has been authenticated
	IsAuthenticated bool `json:"isAuthenticated"`
	// True if this request is to a private area of the system. For example, requests to the public homepage will get
	// a RequestEvent and are considered public; however, incrementing the private data of a user applicaton is
	// considered private. Set to true by default by factory.
	IsProtectedResource bool `json:"isProtectedResource"`
	// The public key used to attempt authentication, if any.
	PublicKey string `json:"publicKey"`
	// Messages attributed to this RequestEvent
	Messages []RequestEventMessage `json:"messages"`
}

// RequestEventMessage captures system messages about a RequestEvent for
// runtime diagnoses, debugging, and otherwise observing system behavior.
type RequestEventMessage struct {
	// An internal ID to track which part of the system generated the message
	System string `json:"system"`
	// The level of the message 'warn','fatal','info',debug'.
	Level string `json:"level"`
	// The message
	Message string `json:"message"`
	// Details from the method
	Detail string `json:"detail"`
	// If true this message is private and only visible by system operators
	Private string `json:"private"`
	// CaptureTime
	CaptureTime time.Time `json:"captureTime"`
}

// RequestEvent constructs a new instance complete with a stamped startime.
func NewRequestEvent(r *http.Request) *RequestEvent {
	var start = time.Now()
	var re = new(RequestEvent)

	id, err := snowflake.Next()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	re.ID = id
	re.IsHTTP = true
	re.IsProtectedResource = true
	re.StartTime = start

	path := r.URL.Host + r.URL.Path
	if len(r.URL.RawQuery) > 0 {
		path += "?" + r.URL.RawQuery
	}
	if len(r.URL.Fragment) > 0 {
		path += "#" + r.URL.Fragment
	}

	re.Address = path
	log.Println("[RequestEvent] - Created a new RequestEvent for request with an address of:", re.Address)

	return re
}

// GetCurrentRuntime reads the starttime for the RequestEvent and returns the
// elapsed time duration.
func (re RequestEvent) GetCurrentRuntime() time.Duration {
	return time.Since(re.StartTime)
}
