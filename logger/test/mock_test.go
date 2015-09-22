package test

import (
	"testing"

	"github.com/go-gia/go-infrastructure/logger"
)

func TestStdoutLogging(t *testing.T) {
	settings := logger.Settings{
		Output: logger.Stdout{
			Level: "panic",
		},
	}

	log, err := logger.New(settings)
	if err != nil {
		t.Fatal("Error", err)
	}

	log.Context(logger.Fields{"foo": "bar1"}).Debug("Debug")
}

func TestStderrLogging(t *testing.T) {
	settings := logger.Settings{
		Output: logger.Stderr{
			Level: "panic",
		},
	}

	log, err := logger.New(settings)
	if err != nil {
		t.Fatal("Error", err)
	}

	log.Context(logger.Fields{"foo": "bar2"}).Debug("Debug")
}

func TestDiskLogging(t *testing.T) {
	settings := logger.Settings{
		Output: logger.Disk{
			Path:  "output.log",
			Level: "panic",
		},
	}

	log, err := logger.New(settings)
	if err != nil {
		t.Fatal("Error", err)
	}

	log.Context(logger.Fields{"foo": "bar3"}).Debug("Debug")
}

func TestLogglyLogging(t *testing.T) {
	settings := logger.Settings{
		Output: logger.LogglySettings{
			Domain: "web-rat.com",
			Token:  "gibberish",
			Level:  "panic",
			Tags:   []string{"version-1", "stuff", "more-stuff"},
		},
	}

	log, err := logger.New(settings)
	if err != nil {
		t.Fatal("Error", err)
	}
	log.Context(logger.Fields{"foo": "bar4"}).Debug("Debug")
}
