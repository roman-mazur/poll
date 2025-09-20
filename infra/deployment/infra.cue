package deployment

import (
	"encoding/json"

	"github.com/roman-mazur/cuetf/aws"
	"github.com/roman-mazur/cuetf/cloudflare"
	"rmazur.io/poll-defs/infra/monitoring"
)

cwa: monitoring.cwa

awsRegion: "eu-central-1"

terraform: {
	terraform: required_providers: {
		aws: version: "= 6.14.0"
		cloudflare: {
			source:  "cloudflare/cloudflare"
			version: "= 5.10.1"
		}
	}

	aws.#Terraform
	provider: aws: region: awsRegion

	#EC2Permissions & {#serverName: "poll_server"}

	resource: aws_instance: poll_server: {
		ami:                  "${data.aws_ami.poll_server_ami.id}"
		iam_instance_profile: "${aws_iam_instance_profile.poll_server.name}"
		instance_type:        selectedInstanceType.name
		tags: Name: "pollsvc server"

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

		echo "poll svc started" >> /opt/init.log

		yum install -y amazon-cloudwatch-agent
		echo '\(json.Marshal(cwa.config))' > \(cwa.installDir)/bin/config.json
		\(cwa.installDir)/bin/amazon-cloudwatch-agent-ctl -a fetch-config -m ec2 -s -c file:\(cwa.installDir)/bin/config.json

		echo "CWA started" >> /opt/init.log
		"""
	}

	resource: aws_eip: poll_server_ip: {
		instance: "${aws_instance.poll_server.id}"
		domain:   "vpc"
	}

	#amazonLinuxVersion: "2023"
	data: aws_ami: poll_server_ami: {
		most_recent: true

		filter: [
			{name: "name", values: ["al\(#amazonLinuxVersion)-ami-\(#amazonLinuxVersion)*"]},
			{name: "virtualization-type", values: selectedInstanceType.info.SupportedVirtualizationTypes},
			{name: "architecture", values: selectedInstanceType.info.ProcessorInfo.SupportedArchitectures},
		]
		owners: ["amazon"]
	}

	cloudflare.#Terraform & {
		resource: cloudflare_dns_record: poll_server: {
			zone_id: "d383a7704b48586d1bc8c2f949712e28"
			name:    "poll"
			content: "${aws_eip.poll_server_ip.public_ip}"
			type:    "A"
			ttl:     1
			proxied: true
		}
	}

	output: poll_server_host_name: value: "${aws_instance.poll_server.private_dns}"
}

#EC2Permissions: aws.#Terraform & {
	#serverName: string

	#assumePolicy: {
		Version: "2012-10-17"
		Statement: [{
			Action: "sts:AssumeRole"
			Principal: Service: "ec2.amazonaws.com"
			Effect: "Allow"
		}]
	}

	resource: aws_iam_role: {
		(#serverName): {
			name:               #serverName
			path:               "/system/"
			assume_role_policy: json.Marshal(#assumePolicy)
			description:        "Role for \(#serverName)"
		}
	}
	resource: aws_iam_role_policy_attachment: {
		"\(#serverName)-cwa": {
			role:       "${aws_iam_role.\(#serverName).name}"
			policy_arn: "arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy"
		}
	}
	resource: aws_iam_instance_profile: {
		(#serverName): {
			name: #serverName
			role: "${aws_iam_role.\(#serverName).name}"
			path: "/system/"
		}
	}
}
