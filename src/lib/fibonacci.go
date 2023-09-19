package lib

import "sync"

var cache = make(map[int]int)
var mu sync.Mutex

func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

func FibonacciCache(n int) int {
	mu.Lock()
	if val, found := cache[n]; found {
		mu.Unlock()
		return val
	}
	mu.Unlock()

	if n <= 1 {
		return n
	}

	res := FibonacciCache(n-1) + FibonacciCache(n-2)

	mu.Lock()
	cache[n] = res
	mu.Unlock()

	return res
}
