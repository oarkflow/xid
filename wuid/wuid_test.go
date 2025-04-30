package wuid

import (
	"testing"

	"github.com/oarkflow/xid"
)

func BenchmarkNew(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = New()
	}
}

func BenchmarkNewInt64(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = New().Int64()
	}
}

func BenchmarkGenerateIDString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = GenerateIDString()
	}
}

func BenchmarkGenerateWithTimestamp(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = GenerateWithTimestamp()
	}
}

func BenchmarkWuidGenerateID(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = xid.New().String()
	}
}
