package webserver_test

import (
	. "github.com/wreckerlabs/webserver"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Webserver without any configuration", func() {

	_ = New()

	It("should have LogDebugMessages disabled by default", func() {
		Expect(Settings.LogDebugMessages).To(BeFalse())
	})
	It("should have LogWarningMessages disabled by default", func() {
		Expect(Settings.LogWarningMessages).To(BeFalse())
	})
	It("should have EnableStaticFileServer disabled by default", func() {
		Expect(Settings.EnableStaticFileServer).To(BeFalse())
	})
	Context("by convention", func() {
		It("should have the expect number of SystemTemplates", func() {
			Expect(len(Settings.SystemTemplates)).Should(Equal(1))
		})
		It("should have an onMissingHandler", func() {
			Expect(Settings.SystemTemplates["onMissingHandler"]).Should(Equal("errors/onMissingHandler"))
		})
	})
})
