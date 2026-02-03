package umdd

import "time"

// vote-use-case-json OMIT
{
	"useCase": {
		"vote": {
			"submitRate": "5s"
			"concurentUsers": 1000
		}
	}
}
// vote-use-case-json-end OMIT

// vote-use-case-cue OMIT
useCase: vote: {
	submitRate: "5s"
	concurentUsers: 1000
}
// vote-use-case-cue-end OMIT

// constraints-1 OMIT

// In the same file we can add:

useCase: [name=string]: {
	submitPeriod:   string
	concurentUsers: >0 & <=5000 // HL
}
// constraints-1-end OMIT

// constraints-2 OMIT
//import "time"

useCase: [name=string]: {
	submitPeriod:   string & time.Duration // HL
	concurentUsers: >0 & <=5000
}
// constraints-2-end OMIT

// constraints-3 OMIT
useCase: [name=string]: close({ // HL
	submitPeriod:   string & time.Duration
	concurentUsers: >0 & <=5000
})
// constraints-3-end OMIT

// all-use-cases OMIT
useCase: {
	vote: {
		submitPeriod:   "5s"
		concurentUsers: 1000
	}
	annotate: {
		submitPeriod:   "2m"
		concurentUsers: 10
	}
	fetch: {
		submitPeriod:   "5s"
		concurentUsers: annotate.concurentUsers
	}
}
useCase: [name=string]: close({
	submitPeriod:   string & time.Duration
	concurentUsers: >0 & <=5000
})
// all-use-cases-end OMIT

// call-rate OMIT
useCase: [name=string]: close({
	submitPeriod:   string & time.Duration
	concurentUsers: >0 & <=5000

	CPS: concurentUsers / (time.ParseDuration(submitPeriod) / time.ParseDuration("1s")) // HL
})
// call-rate-end OMIT

// label-sample OMIT
useCase: annotate: recordSample: structure.#Label & {

	talk_name: strings.Join(list.Repeat(["x"], 80), "")

	timestamp: "2022-02-24T03:00:00Z"

	name:      strings.Join(list.Repeat(["x"], 50), "")
}
// label-sample-end OMIT

// memory-requirement OMIT
eventDuration: "8h" // HL

useCase: [name=string]: close({
	// ...
	submitPeriod: string & time.Duration

	recordSample: _
	recordSize:   len(json.Marshal(recordSample)) // HL

	_recordCount: time.ParseDuration(eventDuration) / time.ParseDuration(submitPeriod) * concurentUsers

	memory: _recordCount * recordSize // HL
})
// memory-requirement-end OMIT

// summary OMIT
summary: {
	CPS: math.Round(list.Sum([for uc in useCase { uc.CPS }]))

	memoryMB: math.Round(list.Sum([for uc in useCase { uc.memory }]) / 1024 / 1024)
}
// summary-end OMIT

// summary-limited-cps OMIT
summary: {
	CPS: math.Round(list.Sum([for uc in useCase { uc.CPS }]))
	CPS: <=60 // HL
}
// summary-limited-cps-end OMIT

// vote-use-case-modified OMIT
summary: CPS: <=60 // HL

useCase: vote: {
	submitRate: "20s" // HL
	usersCount: 1000
}
// vote-use-case-modified-end OMIT
