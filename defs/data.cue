package defs

import (
	"encoding/json"
	"list"
	"math"
	"strings"
	"time"

	"rmazur.io/poll/votes"
)

eventDuration: "8h"

vote: #dataModel & {
	sample: votes.#Vote & {
		talk_name:  strings.Join(list.Repeat(["x"], 50), "")
		voter_id:   "b70fc2dc-4562-43f4-809f-153783bcfd41"
		timestamp: time.Parse(time.RFC3339, "2022-02-24T03:00:00Z")
		value:     10
	}

	submitPeriod: "5s"
	usersCount:   1000
}

label: #dataModel & {
	sample: votes.#Label & {
		talk_name:  strings.Join(list.Repeat(["x"], 50), "")
		timestamp: time.Parse(time.RFC3339, "2022-02-24T03:00:00Z")
		name:      strings.Join(list.Repeat(["x"], 50), "")
	}

	submitPeriod: "2m"
	usersCount:   50
}

#dataModel: {
	sample:     _
	recordSize: len(json.Marshal(sample))

	submitPeriod:       string
	_submitPeriodValid: time.Duration(submitPeriod) & true

	recordsCountPerUser: uint & math.Round(time.ParseDuration(eventDuration)/time.ParseDuration(submitPeriod))
	usersCount:          uint

	recordsCount: usersCount * recordsCountPerUser
	memorySize:   recordsCount * recordSize
}
