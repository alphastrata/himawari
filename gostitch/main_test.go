package main

import "testing"

func BenchmarkConcatFullDisc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		processFullDiscs("../tiles/20220121/214000", "benchmark_out.jpg", false)

	}

}
func BenchmarkConcatFullDisc2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		processFullDiscs("../tiles/20220121/233000", "benchmark2_out.jpg", false)

	}

}
