package deployment

import (
	cloudRegion "github.com/roman-mazur/cuetf/aws/regions/eucentral1"
	"rmazur.io/poll-defs/infra/model"
	"rmazur.io/poll-defs/infra/monitoring"
)

dcRegion: cloudRegion.#Name

pollSvc: {
	version:      "v0.0.11"
	arch:         selectedInstanceType.info.ProcessorInfo.SupportedArchitectures[0]
	downloadLink: "https://github.com/roman-mazur/poll/releases/download/\(version)/pollsvc-\(arch)-linux"
	memReq:       model.summary.memoryMB

	installPath: "/usr/bin/pollsvc"
	certPath:    "/opt/pollsvc"

	systemd: """
	[Unit]
	Description=Poll Service
	After=network.target
	StartLimitIntervalSec=0

	[Service]
	Type=simple
	Restart=always
	RestartSec=1
	User=root
	Environment="OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:\(monitoring.adot.otlpPort)"
	ExecStart=\(installPath) --addr=:443 --tls=\(certPath) --admin-secret="\(inputs.admin.secret)"

	[Install]
	WantedBy=multi-user.target
	"""
}

instanceFilter: {
	CurrentGeneration: true
	FreeTierEligible:  true
	MemoryInfo: SizeInMiB: >pollSvc.memReq & <=(pollSvc.memReq * 5)
}

selectedInstanceType: {
	candidates: [for c in cloudRegion.InstanceTypes if (c & instanceFilter) != _|_ {c}]
	candidateNames: [for c in candidates {c.InstanceType}]

	info: candidates[len(candidates)-1]
	name: info.InstanceType
}
