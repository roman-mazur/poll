//go:build linux || darwin

package main

import (
	"os"
	"syscall"
)

var shutdownSignals = []os.Signal{syscall.SIGTERM, os.Interrupt}
