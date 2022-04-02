package goretry

import "time"

type BackoffStrategy int

const (
	StrategyConstant BackoffStrategy = iota
	StrategyLinear
	StrategyFibonacci
)

type Strategy func(times int) time.Duration

func Constant(d time.Duration) Strategy {
	return func(times int) time.Duration {
		return d
	}
}

func Linear(d time.Duration) Strategy {
	return func(times int) time.Duration {
		return (d * time.Duration(times))
	}
}

func Fibonacci(d time.Duration) Strategy {
	return func(times int) time.Duration {
		return (d * time.Duration(fibonacciNumber(times)))
	}
}

func fibonacciNumber(n int) int {
	if n == 0 || n == 1 {
		return n
	}
	return fibonacciNumber(n-1) + fibonacciNumber(n-2)
}
