package main

import (
	"fmt"

	"github.com/go-gia/go-infrastructure/logger"
	"github.com/go-gia/go-infrastructure/webserver"
	"github.com/go-gia/go-infrastructure/webserver/context"
)

// AppWebInterface represents a application specific structure. You may want to add
// fields for loggers, repos, interactors, or other things to enhance your
// handler.
type AppWebInterface struct {
	Endpoints []webserver.HandlerDef
	logger    logger.Logger
}

// *****************************************************************************
// HANDLER SAMPLE
// *****************************************************************************

// We'll demo the webserver using an API for kicks and grins.
var api = &AppWebInterface{
	Endpoints: []webserver.HandlerDef{},
	logger:    logger.New(false, true),
}

func init() {
	// A simple example of how to register one handler using our HandlerDef system
	// to auto-document this API endpoint
	api.Endpoints = append(api.Endpoints, webserver.HandlerDef{
		Alias:               "ExampleGet",
		Method:              "GET",
		Path:                "/api/sample",
		Documentation:       "/some/path/to/documentation.html",
		DurationExpectation: "1ms",
		Handler: func(ctx *context.Context) {
			ctx.HTML("Listing stuff")
		},
	})
}

// *****************************************************************************
// COMPLEX HANDLER SAMPLE
// *****************************************************************************

// ComplexSampleParams demonstrates how we can define params that a handler can
// use that also support documentation.
type ComplexSampleParams struct {
	// Firstname is just an example
	Firstname string `json:"firstName"`
	// Lastname is another example
	Lastname string `json:"lastName"`
}

// ComplexSample is a handler with metadat declared
func (h AppWebInterface) ComplexSample(ctx *context.Context) {
	h.logger.Debug("Primary handler is executing")
	ctx.HTML("Everything is awesome when you are part of a team.")
}

// RequireSession will be registered to execute before our
// primary handler.
func RequireSession(ctx *context.Context) {
	fmt.Printf("Pretend we're validating a session")
}

// ThrottleCheck will be registered to execute before our
// primary handler.
func ThrottleCheck(ctx *context.Context) {
	fmt.Printf("Pretend we're checking for rate limits")
}

// ComplexSamplePostHandler will be registered to execute after our
// primary handler. This could be some form of analysis or something else.
func ComplexSamplePostHandler(ctx *context.Context) {
	fmt.Printf("PostHandler is executing")
}

func init() {
	// AuthHandlerDef could be used by many handlers
	AuthHandlerDef := webserver.HandlerDef{
		Alias:         "AuthCheck",
		Documentation: "/some/path/to/documentation/about/auth.html",
		Handler:       RequireSession,
	}
	// ThrottleHandlerDef could be used by many handlers
	ThrottleHandlerDef := webserver.HandlerDef{
		Alias:         "ThrottleCheck",
		Documentation: "/some/path/to/documentation/about/throttle.html",
		Handler:       ThrottleCheck,
	}

	// Bring the complex stuff home now...
	api.Endpoints = append(api.Endpoints, webserver.HandlerDef{
		Alias:               "ExampleComplexGet",
		Method:              "GET",
		Path:                "/api/sample/two",
		Documentation:       "/some/path/to/documentation2.html",
		DurationExpectation: "1ms",
		Handler: func(ctx *context.Context) {
			ctx.HTML("Listing stuff 2")
		},
		Params:      ComplexSampleParams{},
		PreHandlers: []webserver.HandlerDef{AuthHandlerDef, ThrottleHandlerDef},
	})
}

// *****************************************************************************
// BRING IT ALL TOGTHER
// *****************************************************************************

func main() {
	log := logger.New(false, true)
	ws := webserver.New(log)

	// You can serve static files. It is quite easy.
	ws.FILES("/public", "static")

	// We can pass the HandlerDefs we've created in this application specific
	// Handler to the Webserver and it will now read the HandlerDef(s) and
	// register those HandlerFunc(s) with the webserver
	ws.RegisterHandlerDefs(api.Endpoints)

	// We can also register a HandlerFunc manualy without all that HandlerDef fuss.
	ws.GET("/", func(ctx *context.Context) {
		ctx.HTML("Winter is coming--but, it's not here yet.")
	})

	// Once the webserver is configured you will want to listen for clients
	ws.Start(":8888")
}
