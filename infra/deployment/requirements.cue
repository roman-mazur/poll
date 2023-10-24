package deployment

import (
	"rmazur.io/cuetf/aws/regions/eucentral1"
	"rmazur.io/poll-defs/infra/model"
)

#memReq: model.summary.req.memMB

instanceFilter: {
	CurrentGeneration: true
	FreeTierEligible: true
	MemoryInfo: SizeInMiB: >#memReq & <=(#memReq*2)
}

awsInstanceType: {
	candidates: [for c in eucentral1.InstanceTypes if (c & instanceFilter) != _|_ {c.InstanceType}]
	name: candidates[0]
}
