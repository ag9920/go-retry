package goretry

import "time"

type Option func(c *Config)

func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

func WithMaxRetryTimes(times int) Option {
	return func(c *Config) {
		c.MaxRetryTimes = times
	}
}

func WithRecoverPanic() Option {
	return func(c *Config) {
		c.RecoverPanic = true
	}
}

func WithBeforeHook(hook HookFunc) Option {
	return func(c *Config) {
		c.BeforeTry = hook
	}
}

func WithAfterHook(hook HookFunc) Option {
	return func(c *Config) {
		c.AfterTry = hook
	}
}

func WithRetryChecker(checker RetryChecker) Option {
	return func(c *Config) {
		c.RetryChecker = checker
	}
}

func WithBackOffStrategy(s BackoffStrategy, duration time.Duration) Option {
	return func(c *Config) {
		switch s {
		case StrategyConstant:
			c.Strategy = Constant(duration)
		case StrategyLinear:
			c.Strategy = Linear(duration)
		case StrategyFibonacci:
			c.Strategy = Fibonacci(duration)
		}
	}
}

func WithCustomBackOffStrategy(s Strategy) Option {
	return func(c *Config) {
		c.Strategy = s
	}
}
