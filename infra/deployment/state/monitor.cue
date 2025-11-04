package state

import (
	"rmazur.io/poll-defs/infra/deployment"
	"rmazur.io/poll-defs/infra/model"
	"rmazur.io/poll-defs/infra/monitoring"
)

checks: {
	live: monitoring.#ServerLivenessCheck & {
		#addr: deployData.full_address.value
	}

	memory: monitoring.#InstanceMemoryCheck & {
		#region:   deployment.awsRegion
		#hostname: deployData.poll_server_host_name.value
	}
}

// Expect the correct version to be deployed.
outputs: live: version: deployment.pollSvc.version

// Actual memory usage should be less than predicted by the model.
outputs: memory: Datapoints: [...{
	Maximum: <=(model.summary.memory / deployment.selectedInstanceType.info.MemoryInfo.SizeInMiB * 100)
}]
