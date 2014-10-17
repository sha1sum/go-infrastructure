package webserver

// RequestContext is created on every request and models the event, or request.
// The webserver will populate a RequestContext with any data provided by the
// client from a form, URL, or recognized data type sent in the request body.
type RequestContext struct {
	// ID is generated for each new RequestContext
	ID uint64 `json:"requestID"`
}
