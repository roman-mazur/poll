package monitoring

import (
	"encoding/json"
	"time"
)

// A template for a command that checks poll service memory usage.
#ServiceMemoryCheck: {
	#region: string

	_sqlQuery: (#maxSqlQuery & {metricName: "go.memory.used"}).query

	cmd: [
		"aws", "--region", #region,
		"cloudwatch", "get-metric-data",
		"--start-time", "$(date -v -6H +'%Y-%m-%dT%H:%M:%S%z')",
		"--end-time", "$(date +'%Y-%m-%dT%H:%M:%S%z')",
		"--metric-data-queries", "'\(json.Marshal([#metricQuery & {#sql: _sqlQuery}]))'",
	]

	#output: {
		MetricDataResults: [{
			Id: string
			Label: string
			Timestamps: [...time.Time & string]
			Values: [...float64]
			StatusCode: "Complete"
		}]
		Messages: []
	}
}

#maxSqlQuery: {
	metricName: string
	query: #"SELECT MAX("\#(metricName)") FROM SCHEMA(pollsvc, OTelLib)"#
}

#metricQuery: {
	#sql: string

	Id: string | *"q1"
	Expression: #sql
	ReturnData: true
	Period: 60
}

#ServerLivenessCheck: {
	#addr: string

	cmd: ["curl", "https://\(#addr)/ping"]

	#output: version: string
}
