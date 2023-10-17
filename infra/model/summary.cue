package model

import "math"

#MB: {
	#bytes:  uint
	math.Round(#bytes / 1024 / 1024)
}

summary: {
	samples: {
		v: vote.sample
		l: label.sample
	}

	req: {
		memMB: {#MB, #bytes: vote.memorySize + label.memorySize}
	}
}
