// Command infra allows inspecting the infrastructure definitions.
package main

import (
	"log"
	"os"
	"os/exec"
)

//go:generate go run cuelang.org/go/cmd/cue get go rmazur.io/poll/votes

func main() {
	cc := exec.Command("cue", os.Args[1:]...)
	cc.Stdout = os.Stdout
	cc.Stderr = os.Stderr
	cc.Stdin = os.Stdin
	err := cc.Run()
	if err != nil {
		if eErr, ok := err.(*exec.ExitError); ok {
			os.Exit(eErr.ExitCode())
		} else {
			log.Fatalf("problems running %s: %s", cc.Path, err)
		}
	}
}
