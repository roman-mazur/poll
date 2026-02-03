package state

import (
	"rmazur.io/poll-defs/infra/deployment"
	"rmazur.io/poll-defs/infra/model"
	"rmazur.io/poll-defs/infra/monitoring"
)

checks: {
	// Is /ping responding?
	live: monitoring.#ServerLivenessCheck & {
		#addr: deployData.full_address.value
	}

	// App memory usage.
	memory: monitoring.#ServiceMemoryCheck & {
		#region: deployment.dcRegion
	}

	// Operation counters.
	for name, case in model.useCase {
		"operation_\(name)": monitoring.#OperationRateCheck & {
			#region: deployment.dcRegion
			#name:   name
		}
	}
}

// Expect the correct version to be deployed.
outputs: live: version: deployment.pollSvc.version

// Actual memory usage should be less than predicted by the model.
outputs: memory: MetricDataResults: [{
	#v: <=(model.summary.memoryMB * 1024 * 1024)
	Values: [#v, ...#v]
}]

// Validate if our usage model matches actual usage.
for name, case in model.useCase {
	outputs: "operation_\(name)": {
		#v: <=case.CPS
		MetricDataResults: [{Values: [#v, ...#v]}]
	}
}
