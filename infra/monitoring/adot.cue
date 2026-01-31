@extern(embed)
package monitoring

// AWS Distro for OpelTelemetry.
adot: {
	otlpPort:    4318
	setupScript: string @embed(file=adot-build.sh,type=text)
	installDir:  "/opt/aws/aws-otel-collector"
	ctlCmd:      "\(installDir)/bin/aws-otel-collector-ctl"
	configPath:  "\(installDir)/bin/poll-config.json"

	config: {
		receivers: otlp: protocols: http: endpoint: "localhost:\(otlpPort)"
		processors: batch: timeout:       "60s"
		exporters: awsemf: log_retention: 7
		service: pipelines: {
			metrics: {
				receivers: ["otlp"]
				processors: ["batch"]
				exporters: ["awsemf"]
			}
		}
	}
}
