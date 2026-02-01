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

	memory: monitoring.#ServiceMemoryCheck & {
		#region: deployment.dcRegion
	}
}

// Expect the correct version to be deployed.
outputs: live: version: deployment.pollSvc.version

// Actual memory usage should be less than predicted by the model.
outputs: memory: MetricDataResults: [{
	#v: <=(model.summary.memoryMB * 1024 * 1024)
	Values: [#v, ...#v]
}]
