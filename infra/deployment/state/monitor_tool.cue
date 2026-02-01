package state

import (
	"strings"

	"tool/cli"
	"tool/exec"
)

command: check: {
	for name, checkCmd in checks {
		(name): {
			log: cli.Print & {text: "Importing data for [\(name)]"}

			// Get the data exposed by monitoring.
			get: exec.Run & {
				$after: log
				cmd: ["bash", "-c", strings.Join(checkCmd.cmd, " ")]
				stdout: string
			}

			// Import this data into CUE.
			"import": exec.Run & {
				cmd: ["cue", "import", "-p", "state", "-l", "outputs: \(name):", "-f", "-o", "\(name)_out.cue", "json:", "-"]
				stdin: get.stdout
			}
		}

		// Validate the expectations using cue vet.
		validate: exec.Run & {
			$after: [for name, _ in checks {check[name]["import"]}]
			cmd: ["cue", "vet", "-c"]
		}
	}
}
