package goretry

import "time"

type Option func(c *Config)

func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

func WithRetryInterval(interval time.Duration) Option {
	return func(c *Config) {
		c.RetryInterval = interval
	}
}

func WithMaxRetryTimes(times int) Option {
	return func(c *Config) {
		c.MaxRetryTimes = times
	}
}

func WithThrowPanic(throwPanic bool) Option {
	return func(c *Config) {
		c.ThrowPanic = throwPanic
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
