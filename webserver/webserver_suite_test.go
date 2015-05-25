package webserver_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWebserver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Webserver Suite")
}
