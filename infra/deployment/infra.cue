package deployment

import (
	// "rmazur.io/cuetf/aws"
)

terraform: {
	// aws.#Terraform
	provider: aws: region: "eu-central-1"

	resource: aws_instance: poll_server: {
		ami: "${data.aws_ami.poll_server_ami.id}"
		instance_type: awsInstanceType.name
		tags: Name: "pollsvc server"
	}

	data: aws_ami: poll_server_ami: {
		most_recent: true

		filter: [
			{name: "name", values: ["al2023-*"]},
			{name: "virtualization-type", values: awsInstanceType.info.SupportedVirtualizationTypes},
			{name: "architecture", values: ["x86_64" & or(awsInstanceType.info.ProcessorInfo.SupportedArchitectures)]},
			{name: "boot-mode", values: ["uefi-preferred"]},
		]
		owners: ["amazon"]
	}
}
