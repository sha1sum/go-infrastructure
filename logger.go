package webserver

import (
	"log"
	"os"
	"time"
)

var (
	green  = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white  = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red    = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	reset  = string([]byte{27, 91, 48, 109})
)

var (
	stdlogger = log.New(os.Stdout, "", 0)
	errlogger = log.New(os.Stderr, "", 0)
)

func consolelog(level string, message string) {

	var color string
	switch {
	case level == "info":
		color = white
	case level == "debug":
		color = yellow
	case level == "error":
		color = red
	case level == "success":
		color = green
	default:
		color = red
	}

	stdlogger.Printf("[webserver] %v | %s %s %4v | %s",
		time.Now().Format("2006/01/02 - 15:04:05"),
		color, level, reset, message,
	)
}
