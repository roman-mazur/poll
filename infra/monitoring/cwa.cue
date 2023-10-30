package monitoring

// AWS CloudWatch agent configuration.
cwa: {
	installDir: "/opt/aws/amazon-cloudwatch-agent"

	config: {
		"agent": {
			"metrics_collection_interval": 60
			"run_as_user":                 "root"
		}
		"metrics": {
			"aggregation_dimensions": [
				[
					"InstanceId",
				],
			]
			"metrics_collected": {
				"cpu": {
					"measurement": [
						"cpu_usage_idle",
						"cpu_usage_iowait",
						"cpu_usage_user",
						"cpu_usage_system",
					]
					"metrics_collection_interval": 60
					"resources": [
						"*",
					]
					"totalcpu": false
				}
				"disk": {
					"measurement": [
						"used_percent",
						"inodes_free",
					]
					"metrics_collection_interval": 60
					"resources": [
						"*",
					]
				}
				"diskio": {
					"measurement": [
						"io_time",
					]
					"metrics_collection_interval": 60
					"resources": [
						"*",
					]
				}
				"mem": {
					"measurement": [
						"mem_used_percent",
					]
					"metrics_collection_interval": 60
				}
				"swap": {
					"measurement": [
						"swap_used_percent",
					]
					"metrics_collection_interval": 60
				}
			}
		}
	}
}
