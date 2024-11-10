package goso

import (
	"testing"
)

func BenchmarkTest(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
	}
}
