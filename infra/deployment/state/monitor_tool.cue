package state

import (
	"strings"

	"tool/cli"
	"tool/exec"
)

command: check: {
	for name, command in checks {
		(name): {
			log: cli.Print & {text: "Performing check [\(name)]"}

			// Get the data exposed by monitoring.
			get: exec.Run & {
				$after: log
				cmd: ["bash", "-c", strings.Join(command.cmd, " ")]
				stdout: string
			}

			// Import this data into CUE.
			"import": exec.Run & {
				cmd: ["cue", "import", "-p", "state", "-l", "outputs: \(name):", "-f", "-o", "\(name)_out.cue", "json:", "-"]
				stdin: get.stdout
			}

			// Validate the expectations using cue vet.
			validate: exec.Run & {
				$after: check[name]["import"]
				cmd: ["cue", "vet"]
			}
		}
	}
}
