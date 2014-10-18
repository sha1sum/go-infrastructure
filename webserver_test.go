package webserver_test

import (
	. "git.wreckerlabs.com/in/webserver"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Webserver", func() {

	ws := New()

	Describe("A webserver", func() {
		Context("without any configuration", func() {
			It("should have no handlers defined", func() {
				Expect(len(ws.Handlers)).To(Equal(0))
			})
			It("should have a root RouteNamespace", func() {
				Expect(ws.Prefix).To(Equal("/"))
			})
		})
	})
})
