package deployment

import (
	// "rmazur.io/cuetf/aws"
	"rmazur.io/cuetf/cloudflare"
)

terraform: {
	// aws.#Terraform
	provider: aws: region: "eu-central-1"

	resource: aws_instance: poll_server: {
		ami:           "${data.aws_ami.poll_server_ami.id}"
		instance_type: awsInstanceType.name
		tags: Name: "pollsvc server"
		associate_public_ip_address: false

		user_data: """
		#!/bin/bash

		curl -f -L \(pollSvc.downloadLink) > \(pollSvc.installPath) && chmod 755 \(pollSvc.installPath)
		setcap CAP_NET_BIND_SERVICE=+eip \(pollSvc.installPath)

		cat > /etc/systemd/system/pollsvc.service <<EOF
		\(pollSvc.systemd)
		EOF

		mkdir -p \(pollSvc.certPath)
		cat > \(pollSvc.certPath)/cert.pem <<EOF
		\(inputs.tls.certificate)
		EOF
		cat > \(pollSvc.certPath)/pkey.pem <<EOF
		\(inputs.tls.private_key)
		EOF

		systemctl enable pollsvc
		systemctl start pollsvc

		echo "poll svc started"
		"""
	}

	resource: aws_eip: poll_server_ip: {
		instance: "${aws_instance.poll_server.id}"
		domain:   "vpc"
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

	cloudflare.#Terraform & {
		resource: cloudflare_record: poll_server: {
			zone_id: "d383a7704b48586d1bc8c2f949712e28"
			name:    "poll"
			value:   "${aws_eip.poll_server_ip.public_ip}"
			type:    "A"
			ttl:     1
			proxied: true
		}
	}
}
