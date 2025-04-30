package main

import (
	"testing"

	"github.com/oarkflow/xid"
)

func BenchmarkCustomGenerateID(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = New().String()
	}
}

// BenchmarkWuidGenerateID measures the performance of the wuid library's generator.
// Ensure that the wuid generator is properly initialized.
// Here we assume that wuid provides a Next() function that returns a unique ID.
func BenchmarkWuidGenerateID(b *testing.B) {
	// The wuid library might need some initialization if you haven't done it already.
	// For example: wuid.Next() is used to generate the next unique ID.
	for n := 0; n < b.N; n++ {
		_ = xid.New().String()
	}
}
