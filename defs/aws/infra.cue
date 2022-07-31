package aws

import (
	"rmazur.io/cuetf/aws"
	"rmazur.io/poll-cue/defs"
)

let memReq = defs.summary.req.memMB

itypes: [
	for t in aws.InstanceTypes if t.CurrentGeneration && t.MemoryInfo.SizeInMiB >= memReq && t.MemoryInfo.SizeInMiB < memReq*3 {
		{typ: t.InstanceType, m: t.MemoryInfo.SizeInMiB}
	},
]
