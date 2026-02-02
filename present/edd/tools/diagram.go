package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func diagram(dir, include string) {
	for _, entry := range must(os.ReadDir(dir)) {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".d2") {
			continue
		}
		if include != "" && !strings.Contains(name, include) {
			continue
		}
		fmt.Println(name)
		in := filepath.Join(dir, name)
		out := filepath.Join(dir, strings.TrimSuffix(name, ".d2")+".png")
		d2(in, out)
	}
}

func d2(in, out string) {
	cmd := exec.Command("d2", "-t", "7", "--sketch", in, out)
	res, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		log.Fatal(string(res))
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
