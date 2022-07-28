package defs

#MB: {
	input: uint
	output: input / 1024 / 1024
}

summary: {
	samples: {
		v: vote.sample
		l: label.sample
	}

	req: {
		mem: (#MB & { input: vote.memorySize + label.memorySize }).output
	}
}
