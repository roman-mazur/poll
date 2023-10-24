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
		user_data: """
		#!/bin/bash
		echo "Hello!"

		echo "start pollsvc init" >> /root/init-log
		curl -f -L \(pollSvc.downloadLink) > \(pollSvc.installPath) && chmod 755 \(pollSvc.installPath)

		echo "write pollsvc systemd config" >> /root/init-log
		cat > /etc/systemd/system/pollsvc.service <<EOF
		\(pollSvc.systemd)
		EOF

		systemctl enable pollsvc
		systemctl start pollsvc
		"""
	}

	data: aws_ami: poll_server_ami: {
		most_recent: true

		filter: [
			{name: "name", values: ["al2023-ami-2023*"]},
			{name: "virtualization-type", values: awsInstanceType.info.SupportedVirtualizationTypes},
			{name: "architecture", values: ["x86_64" & or(awsInstanceType.info.ProcessorInfo.SupportedArchitectures)]},
		]
		owners: ["amazon"]
	}
}
