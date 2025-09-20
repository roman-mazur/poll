package ci

import (
	"strings"
)

github: workflows: [n=string]: {name: n}

github: workflows: main: {

	on: {
		push: branches: ["main"]
		pull_request: branches: ["main"]
	}

	jobs: main: {
		"runs-on": "ubuntu-latest"
		steps: [
			{name: "Checkout", uses: "actions/checkout@v4"},
			{name: "Set up Go", uses: "actions/setup-go@v4", with: "go-version": "1.25.1"},
			{
				name: "Test"
				run: strings.Join([
					"go install cuelang.org/go/cmd/cue",
					"go test ./...",
				], "\n")
			},
			{
				name: "Evaluate infra code"
				run:  "cue export -e terraform ./infra/deployment --out cue"
			},
		]
	}
}
