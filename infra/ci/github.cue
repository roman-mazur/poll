package ci

import (
	"strings"

	"cue.dev/x/githubactions"
)

github: workflows: [n=string]: githubactions.#Workflow & {name: n}

github: workflows: main: {

	on: {
		push: branches: ["main"]
		pull_request: branches: ["main"]
	}

	jobs: main: {
		"runs-on": "ubuntu-latest"
		steps: [
			{name: "Checkout", uses: "actions/checkout@v4"},
			{name: "Set up Go", uses: "actions/setup-go@v4", with: "go-version": "1.25.3"},
			#multilineRun & {
				name: "Test"
				#lines: [
					"go install cuelang.org/go/cmd/cue",
					"go test ./...",
				]
			},
			#multilineRun & {
				name: "Evaluate infra code"
				#lines: [
					"go generate ./infra/...",
					"cue export -e terraform ./infra/deployment --out cue",
				]
			},
		]
	}
}

#multilineRun: {
	name: string
	#lines: [...string]
	run: strings.Join(#lines, "\n")
}
