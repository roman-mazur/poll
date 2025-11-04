package monitoring

import "time"

// A template for a command that checks on the specified instance.
#InstanceMemoryCheck: {
	#region:   string
	#hostname: string

	cmd: [
		"aws", "--region", #region,
		"cloudwatch", "get-metric-statistics", "--namespace=CWAgent", "--metric-name=mem_used_percent", "--period=3600",
		"--statistics=Maximum", "--dimensions", "Name=host,Value=\(#hostname)",
		"--start-time", "$(date -v -6H +'%Y-%m-%dT%H:%M:%S%z')",
		"--end-time", "$(date +'%Y-%m-%dT%H:%M:%S%z')",
	]

	#output: {
		Label: "mem_used_percent"
		Datapoints: [#datapoint, ...#datapoint]
		#datapoint: {
			Timestamp: time.Time & string
			Maximum:   float64
			Unit:      "Percent"
		}
	}
}

#ServerLivenessCheck: {
	#addr: string

	cmd: ["curl", "https://\(#addr)/ping"]

	#output: version: string
}
