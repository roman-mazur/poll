package model

import (
	"encoding/json"
	"list"
	"math"
	"strings"
	"time"

	"rmazur.io/poll/votes/structure"
)

_talkNameSample: strings.Join(list.Repeat(["x"], 80), "")

useCase: {
	vote: {
		submitPeriod:   "20s"
		concurentUsers: 1000

		recordSample: structure.#Vote & {
			talk_name: _talkNameSample
			voter_id:  "b70fc2dc"
			timestamp: "2022-02-24T03:00:00Z"
			value:     10
		}
	}

	annotate: {
		submitPeriod:   "2m"
		concurentUsers: 10

		recordSample: structure.#Label & {
			talk_name: _talkNameSample
			timestamp: "2022-02-24T03:00:00Z"
			name:      strings.Join(list.Repeat(["x"], 50), "")
		}
	}

	fetch: {
		submitPeriod:   "5s"
		concurentUsers: annotate.concurentUsers
		recordSample:   null
	}
}

eventDuration: "8h"

useCase: [name=string]: close({
	submitPeriod:   string & time.Duration
	concurentUsers: >0 & <=1_000_000

	CPS: concurentUsers / (time.ParseDuration(submitPeriod) / time.ParseDuration("1s"))

	recordSample: _ | *null
	recordSize:   uint | *0
	if recordSample != null {
		recordSize: len(json.Marshal(recordSample))
	}
	memory: time.ParseDuration(eventDuration) / time.ParseDuration(submitPeriod) * concurentUsers * recordSize
})

summary: {
	CPS: >0 & <=60
	CPS: math.Round(list.Sum([ for uc in useCase {uc.CPS}]))

	memory: math.Round(list.Sum([ for uc in useCase {uc.memory}]) / 1024 / 1024)
}
