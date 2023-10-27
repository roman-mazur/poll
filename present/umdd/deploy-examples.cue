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

selectedInstance: {
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
		{name: "virtualization-type", values: selectedInstanceType.info.SupportedVirtualizationTypes},
	]
	owners: ["amazon"]
}
// final-end OMIT
