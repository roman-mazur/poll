package umdd

// k8s-deployment OMIT
import "rmazur.io/poll-defs/model" // HL

k8s: deployment: {
	kind: "Deployment"
	spec: {
		replicas: 1
		template: spec: {
			containers: [{
				name:  "pollsvc"
				image: "some-registry.io/pollsvc:0.1.0"

				resources: limits: memory: "\(model.summary.memory)Mi" // HL
			}]
		}
	}
}
// k8s-deployment-end OMIT

// terraform-waf OMIT
terraform: aws_waf_rate_based_rule: poll_svc_rate_limits: {
	name:       "poll_svc_rate_limits"
	rate_key:   "IP"
	rate_limit: >=(model.summary.CPS * 300) // Per 5 minutes. // HL
}
// terraform-waf-end OMIT

// instance-filter OMIT
instanceFilter: {
	CurrentGeneration: true
	MemoryInfo: SizeInMiB: >model.summary.memory & <=(model.summary.memory * 2) // HL
}

selectedInstance: { // HL
	candidates: [ for c in InstanceTypes if (c & instanceFilter) != _|_ {c}]

	info:     candidates[0]
	typeName: info.InstanceType
}
// instance-filter-end OMIT

// free-tier OMIT
instanceFilter: {
	CurrentGeneration: true
	FreeTierEligible:  true // HL
	MemoryInfo: SizeInMiB: >model.summary.memory & <=(model.summary.memory * 2)
}
// free-tier-end OMIT

// final OMIT
resource: aws_instance: poll_server: {
	ami:           "${data.aws_ami.poll_server_ami.id}"
	instance_type: selectedInstanceType.name // HL
}

data: aws_ami: poll_server_ami: {
	most_recent: true

	filter: [
		{name: "name", values: ["al2023-ami-2023*"]},
		{
			name: "virtualization-type"
			values: selectedInstanceType.info.SupportedVirtualizationTypes
		},
	]
	owners: ["amazon"]
}
// final-end OMIT

memoryMetricCode: """
memory-stats, OMIT
aws cloudwatch get-metric-statistics --metric-name=mem_used_percent \ // HL
		--namespace=CWAgent \
		--statistics=Maximum --dimensions Name=host,Value=$hostname \
		--start-time "2023-11-09T08:00:00" --end-time "2023-11-10T20:00:00"
memory-stats-end, OMIT
"""

// memory-output OMIT
outputs: memory: {
	Label: "mem_used_percent"
	Datapoints: [{
		Timestamp: "2023-11-09T18:24:00+00:00"
		Maximum:   19.663855173832545 // HL
		Unit:      "Percent"
	}, {
		Timestamp: "2023-11-09T20:24:00+00:00"
		Maximum:   19.637523143386133
		Unit:      "Percent"
	}]
}
// memory-output-end OMIT

// memory-check OMIT
import (
	"rmazur.io/poll-defs/infra/deployment"
	"rmazur.io/poll-defs/infra/model"
)

// Actual memory usage should be less than predicted by the model.
outputs: memory: Datapoints: [...{

	#instanceMem: deployment.selectedInstanceType.info.MemoryInfo.SizeInMiB
	#modelMax: model.summary.memory / #instanceMem  * 100

	Maximum: <= #modelMax // HL
}]
// memory-check-end OMIT
