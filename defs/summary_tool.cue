package defs

import (
	"encoding/yaml"

	"tool/cli"
)

command: {
	for name, value in summary {
		"\(name)": print: (#prettyPrint & {input: value}).output
	}
}

#prettyPrint: {
	input:  _
	output: cli.Print & {
		text: yaml.Marshal(input)
	}
}
