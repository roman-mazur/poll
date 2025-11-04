package state

#StateStructure

#StateStructure: {
	// Terraform outputs.
	deployData: [string]: {
		sensitive: false // Ensure we don't keep sensitive data from Terraform.
		type:      string
		value:     _
	}

	// All the checks we perform on our infra.
	checks: [string]: {
		cmd: [string, ...string]
		#output: _
	}

	// Monitoring outputs.
	outputs: [name=string]: checks[name].#output
}
