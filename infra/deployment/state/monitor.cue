package state

import (
	"rmazur.io/poll-defs/infra/deployment"
	"rmazur.io/poll-defs/infra/model"
	"rmazur.io/poll-defs/infra/monitoring"
)

checks: {
	live: monitoring.#ServerLivenessCheck & {
		#addr: deployData.full_address.value
		#output: version: deployment.pollSvc.version
	}

	memory: monitoring.#InstanceMemoryCheck & {
		#region:   deployment.awsRegion
		#hostname: deployData.poll_server_host_name.value
	}
}

// Ensure we don't keep sensitive data from Terraform.
deployData: [string]: sensitive: false

outputs: [name=string]: checks[name].#output

// Actual memory usage should be less than predicted by the model.
outputs: memory: Datapoints: [...{
	Maximum: <=(model.summary.memory / deployment.selectedInstanceType.info.MemoryInfo.SizeInMiB * 100)
}]
