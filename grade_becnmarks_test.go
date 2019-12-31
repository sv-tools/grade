package main

import "testing"

func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

var result int

func BenchmarkFib(b *testing.B) {
	var r int
	for n := 0; n < b.N; n++ {
		r = fib(10)
	}
	result = r
}

func BenchmarkFibParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var r int
		for pb.Next() {
			r = fib(10)
		}
		result = r
	})
}
