package goretry

import (
	"errors"
	"time"
)

type Config struct {
	MaxRetryTimes int
	Timeout       time.Duration
	RetryChecker  func(err error) (needRetry bool)
	Strategy      Strategy
	RecoverPanic  bool
	BeforeTry     HookFunc
	AfterTry      HookFunc
}

var (
	DefaultMaxRetryTimes = 3
	DefaultTimeout       = time.Minute
	DefaultRetryChecker  = func(err error) bool {
		return !errors.Is(err, ErrorAbort) // not abort error, should continue retry
	}
)

func newDefaultConfig() *Config {
	return &Config{
		MaxRetryTimes: DefaultMaxRetryTimes,
		RetryChecker:  DefaultRetryChecker,
		Timeout:       DefaultTimeout,
		BeforeTry:     func() {},
		AfterTry:      func() {},
	}
}
