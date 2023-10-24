package deployment

import (
	"rmazur.io/cuetf/aws"
	"rmazur.io/cuetf/aws/regions/eucentral1"
)

terraform: {
	aws.#Terraform
	provider: aws: region: eucentral1.#Name

	resource: aws_instance: poll_server: {
		ami: "todo"
		instance_type: awsInstanceType.name
		availability_zones: eucentral1.AvailabilityZones[1].ZoneName
	}

	resource: null_resource: something: {}
}
