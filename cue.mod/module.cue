module: "rmazur.io/poll-defs"
language: {
	version: "v0.9.0"
}
source: {
	kind: "git"
}
deps: {
	"cue.dev/x/githubactions@v0": {
		v:       "v0.2.0"
		default: true
	}
	"github.com/roman-mazur/cuetf@v0": {
		v: "v0.5.2"
	}
}
