package umdd

// k8s-deployment OMIT
import "rmazur.io/poll-defs/infra/model" // HL

k8s: deployment: {
	kind: "Deployment"
	spec: {
		replicas: 1
		template: spec: {
			containers: [{
				name:  "pollsvc"
				image: "some-registry.io/pollsvc:0.1.0"

				resources: limits: memory: "\(model.summary.memoryMB)Mi" // HL
			}]
		}
	}
}
// k8s-deployment-end OMIT

// opentofu-waf OMIT
opentofu: aws_waf_rate_based_rule: poll_svc_rate_limits: {
	name:       "poll_svc_rate_limits"
	rate_key:   "IP"
	rate_limit: >=(model.summary.CPS * 300) // Per 5 minutes. // HL
}
// opentofu-waf-end OMIT

// instance-filter OMIT
filter: {
	CurrentGeneration: true
	MemoryInfo: SizeInMiB: >model.summary.memoryMB & <=(model.summary.memoryMB * 2) // HL
}

selectedInstanceType: {
	candidates: [for c in cloudRegion.InstanceTypes if (c & filter) != _|_ {c}] // HL

	info: candidates[len(candidates)-1]
	name: info.InstanceType
}
// instance-filter-end OMIT

// free-tier OMIT
filter: {
	CurrentGeneration: true
	FreeTierEligible:  true // HL
	MemoryInfo: SizeInMiB: >model.summary.memoryMB & <=(model.summary.memoryMB * 2)
}
// free-tier-end OMIT

// final OMIT
resource: aws_instance: poll_server: {
	ami:           "${data.aws_ami.poll_server_ami.id}"
	instance_type: selectedInstanceType.name // HL
}

data: aws_ami: poll_server_ami: {
	filter: [
		{name: "name", values: ["al2023-ami-2023*"]},
		{
			name: "virtualization-type"
			values: selectedInstanceType.info.SupportedVirtualizationTypes
		},
		{
			name: "architecture",
			values: selectedInstanceType.info.ProcessorInfo.SupportedArchitectures
		},
	]
}
// final-end OMIT

// memory-check OMIT
// PromQL queries to perform to verify the system.
promQuery: [name=string]: string
promQuery: memory: "max(go.memory.used)"

// Prometheus query results.
output: [name=string]: _
output: memory: values: [...<=model.summary.memoryMB*1024*1024]
// memory-check-end OMIT

// operations-check OMIT
for name, case in model.useCase {
	promQuery: "operation_\(name)": "rate(operation.\(name)_total[5m])"

	output: "operation_\(name)": values: [...<=case.CPS] // HL
}
// operations-check-end OMIT
