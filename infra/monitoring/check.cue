package monitoring

import (
	"encoding/json"
	"time"
)

// A template for a check that verifies that the /ping endpoint works and returns expected data.
#ServerLivenessCheck: {
	#addr: string

	cmd: ["curl", "https://\(#addr)/ping"]

	#output: version: string
}

// A template for a command that checks poll service memory usage.
#ServiceMemoryCheck: _awsMetricsCheck & {
	#since: "6H"
	#sqlQuery: #"SELECT MAX("go.memory.used") FROM SCHEMA(pollsvc, OTelLib)"#
}

// A template for a command that checks the actual call rate for the defined operation or scenario.
#OperationRateCheck: _awsMetricsCheck & {
	#since: "6H"
	#name: string
	#sqlQuery: #"SELECT MAX("operation.\#(#name)_total") FROM SCHEMA(pollsvc, OTelLib)"#
}

_awsMetricsCheck: {
	#region: string
	#since: string
	#sqlQuery: string

	cmd: [
		"aws", "--region", #region,
		"cloudwatch", "get-metric-data",
		"--start-time", "$(date -v -\(#since) +'%Y-%m-%dT%H:%M:%S%z')",
		"--end-time", "$(date +'%Y-%m-%dT%H:%M:%S%z')",
		"--metric-data-queries", "'\(json.Marshal([#metricQuery & {#sql: #sqlQuery}]))'",
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
