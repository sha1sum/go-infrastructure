package webserver_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

const suite = "Webserver"

func TestWebserver(t *testing.T) {
	RegisterFailHandler(Fail)
	//junitReporter := reporters.NewJUnitReporter("./reports/junit/" + suite + "_junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, suite+" Suite", []Reporter{})
}
