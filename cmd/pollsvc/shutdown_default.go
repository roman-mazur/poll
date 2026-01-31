//go:build !linux && !darwin

package main

import "os"

var shutdownSignals = []os.Signal{os.Interrupt}
