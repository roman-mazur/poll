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
	candidates: [_]
	candidates: [for c in eucentral1.InstanceTypes if (c & instanceFilter) != _|_ {c}]

	info: candidates[0]
	name: info.InstanceType
}
