package test

import (
	"testing"

	"github.com/go-gia/go-infrastructure/logger"
)

func TestLogging(t *testing.T) {
	log := logger.NewLogMock(false, false)

	log.Context(logger.Fields{
		"Stuff": false,
	}).Debug("Recording Stuff")

}
