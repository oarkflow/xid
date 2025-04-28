package xid

import (
	"testing"
)

func BenchmarkIDString(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = New().String()
	}
}

// ...existing code or additional tests...
