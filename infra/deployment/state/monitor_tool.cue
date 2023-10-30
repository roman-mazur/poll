package state

import (
	"strings"

	"tool/exec"
)

command: check: {
	for name, command in checks {
		(name): {
			get: exec.Run & {
				cmd: ["bash", "-c", strings.Join(command.cmd, " ")]
				stdout: string
			}
			"import": exec.Run & {
				cmd: ["cue", "import", "-p", "state", "-l", "outputs: \(name):", "-f", "-o", "\(name)_out.cue", "json:", "-"]
				stdin: get.stdout
			}
			validate: exec.Run & {
				$after: check[name]["import"]
				cmd: ["cue", "vet"]
			}
		}
	}
}
