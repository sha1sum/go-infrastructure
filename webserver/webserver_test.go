package webserver_test

import (
	"io/ioutil"
	"log"

	"github.com/aarongreenlee/webserver/context"

	. "github.com/aarongreenlee/webserver"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var logger = log.New(ioutil.Discard, "Test ", log.Ldate|log.Ltime)

var apiHandlerDef = HandlerDef{
	Alias:               "ExampleGet",
	Method:              "GET",
	Path:                "/api/sample",
	Documentation:       "/some/path/to/documentation.html",
	DurationExpectation: "1ms",
	Handler: func(ctx *context.Context) {
		ctx.HTML("Listing stuff")
	}}
var pageHandlerDef = HandlerDef{
	Alias:               "ExampleWebpage",
	Method:              "GET",
	Path:                "/",
	Documentation:       "/some/path/to/documentation.html",
	DurationExpectation: "1ms",
	Handler: func(ctx *context.Context) {
		ctx.HTML("Some page")
	}}
var authHandlerDef = HandlerDef{
	Alias:               "AuthMiddleware",
	Method:              "",
	Path:                "",
	Documentation:       "/some/path/to/documentation.html",
	DurationExpectation: "1ms",
	Handler: func(ctx *context.Context) {
		// pretend we manage sessions or something
	}}

var _ = Describe("Webserver without any configuration", func() {
	It("should have EnableStaticFileServer disabled by default", func() {
		Expect(Settings.EnableStaticFileServer).To(BeFalse())
	})
	Context("by convention", func() {
		It("should have only one SystemTemplates (404)", func() {
			Expect(len(Settings.SystemTemplates)).Should(Equal(1))
		})
		It("should have an onMissingHandler template defined by default", func() {
			Expect(Settings.SystemTemplates["onMissingHandler"]).Should(Equal("errors/onMissingHandler"))
		})
	})
})

var _ = Describe("Webserver RegisterHandlerDef", func() {
	ws := New(logger, logger, logger, logger)
	ws.RegisterHandlerDef(apiHandlerDef)

	It("should create a map with a key for each HandlerDef provided", func() {
		Expect(len(ws.HandlerDef)).Should(Equal(1))
	})
	It("should store the HandlerDef by method:path", func() {
		_, ok := ws.HandlerDef["GET:/api/sample"]
		Expect(ok).Should(Equal(true))
	})
})

var _ = Describe("Webserver RegisterHandlerDefs", func() {
	ws := New(logger, logger, logger, logger)
	ws.RegisterHandlerDefs([]HandlerDef{apiHandlerDef, pageHandlerDef})

	It("should register all handlers provided", func() {
		Expect(len(ws.HandlerDef)).Should(Equal(2))
	})
	It("should store the HandlerDefs by method:path", func() {
		var ok bool
		_, ok = ws.HandlerDef["GET:/api/sample"]
		Expect(ok).Should(Equal(true))
		_, ok = ws.HandlerDef["GET:/"]
		Expect(ok).Should(Equal(true))
		_, ok = ws.HandlerDef["POST:/api/sample"]
		Expect(ok).Should(Equal(false))
	})
})
