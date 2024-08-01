package server

import (
	"qooked/internal/config"
	"qooked/internal/instrumentation"
)

// Server definition
type Server struct {
	config config.Config
	instrumentation instrumentation.Instrumentation
	// router class
}

// initialize method

// run server method