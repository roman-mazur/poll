package model

import (
	"encoding/json"
	"list"
	"math"
	"strings"
	"time"

	"rmazur.io/poll/votes/structure"
)

eventDuration: "8h"

#dataModel: {
	sample:     _
	recordSize: len(json.Marshal(sample))

	submitPeriod:       string
	#submitPeriodValid: time.Duration(submitPeriod) & true

	recordsCountPerUser: uint & math.Round(time.ParseDuration(eventDuration)/time.ParseDuration(submitPeriod))
	usersCount:          uint

	recordsCount: usersCount * recordsCountPerUser
	memorySize:   recordsCount * recordSize
}

#talkNameSample: strings.Join(list.Repeat(["x"], 80), "")

vote: #dataModel & {
	sample: structure.#Vote & {
		talk_name: #talkNameSample
		voter_id:  "b70fc2dc"
		timestamp: time.Parse(time.RFC3339, "2022-02-24T03:00:00Z")
		value:     10
	}

	submitPeriod: "5s"
	usersCount:   1000
}

label: #dataModel & {
	sample: structure.#Label & {
		talk_name: #talkNameSample
		timestamp: time.Parse(time.RFC3339, "2022-02-24T03:00:00Z")
		name:      strings.Join(list.Repeat(["x"], 50), "")
	}

	submitPeriod: "2m"
	usersCount:   50
}
