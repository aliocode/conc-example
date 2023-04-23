package concexample

import "testing"

func BenchmarkWithConcWgNoPanics(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WithConcWgNoPanics()
	}
}

func BenchmarkWithBuiltinWgNoPanics(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WithBuiltinWgNoPanics()
	}
}

func BenchmarkWithConcWgRecovered(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WithConcWgRecovered()
	}
}

func BenchmarkWithBuiltinWgRecovered(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WithBuiltinWgRecovered()
	}
}
