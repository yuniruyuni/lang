package main

import (
	"testing"
)

const code = `
if 1 == 1 {
	if 2 == 2 {
		if 3 == 3 {
			1
		} else {
			0
		}
	} else {
		0
	}
} else {
	0
}`

func BenchmarkCompile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		//nolint:errcheck // for benchmark, err checking should be skipped.
		Compile(code)
		b.StopTimer()
	}
}
