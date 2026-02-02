package main

import (
	"flag"
	"log"
)

func main() {
	flag.Parse()
	cmd := flag.Arg(0)
	if cmd == "" {
		log.Fatal("missing command")
	}

	switch cmd {
	case "diagram":
		dir := flag.Arg(1)
		if dir == "" {
			dir = "."
		}
		diagram(dir, flag.Arg(2))

	default:
		log.Fatalf("unknown command: %s", cmd)
	}
}
