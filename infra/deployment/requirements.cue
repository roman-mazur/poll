package deployment

import (
	"rmazur.io/cuetf/aws/regions/eucentral1"
	"rmazur.io/poll-defs/infra/model"
)

pollSvc: {
	version:      "v0.0.6"
	downloadLink: "https://github.com/roman-mazur/poll/releases/download/\(version)/pollsvc-amd64"
	memReq:       model.summary.req.memMB

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
	ExecStart=\(installPath) --addr=:443 --tls=\(certPath) --admin-secret="\(inputs.admin.secret)"

	[Install]
	WantedBy=multi-user.target
	"""
}

instanceFilter: {
	CurrentGeneration: true
	FreeTierEligible:  true
	MemoryInfo: SizeInMiB: >model.summary.memory & <=(model.summary.memory * 5)
}

selectedInstanceType: {
	candidates: [_]
	candidates: [ for c in eucentral1.InstanceTypes if (c & instanceFilter) != _|_ {c}]

	info: candidates[0]
	name: info.InstanceType
}
