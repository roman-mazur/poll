package deployment

import (
	"github.com/roman-mazur/cuetf/aws/regions/eucentral1"
	"rmazur.io/poll-defs/infra/model"
	"rmazur.io/poll-defs/infra/monitoring"
)

pollSvc: {
	version:      "v0.0.10"
	arch:         selectedInstanceType.info.ProcessorInfo.SupportedArchitectures[0]
	downloadLink: "https://github.com/roman-mazur/poll/releases/download/\(version)/pollsvc-\(arch)-linux"
	memReq:       model.summary.memory

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
	candidates: [for c in eucentral1.InstanceTypes if (c & instanceFilter) != _|_ {c}]

	info: candidates[0]
	name: info.InstanceType
}
