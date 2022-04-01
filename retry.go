package goretry

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"time"
)

type RetryFunc func() error
type HookFunc func()
type RetryChecker func(err error) (needRetry bool)

func Do(ctx context.Context, fn RetryFunc, opts ...Option) error {
	if fn == nil {
		return ErrorEmptyRetryFunc
	}
	var (
		abort         = make(chan struct{}, 1) // caller choose to abort retry
		overload      = make(chan struct{}, 1) // downstream return overload signal
		run           = make(chan error, 1)
		panicInfoChan = make(chan string, 1)

		timer  *time.Timer
		runErr error
	)
	config := newDefaultConfig()
	for _, o := range opts {
		o(config)
	}

	if config.Timeout > 0 {
		timer = time.NewTimer(config.Timeout)
	}

	go func() {
		var err error
		defer func() {
			if e := recover(); e == nil {
				return
			} else {
				panicInfoChan <- fmt.Sprintf("retry panic detected, err=%v, stack:%s", e, debug.Stack())
			}
		}()
		for i := 0; i < config.MaxRetryTimes; i++ {
			config.BeforeTry()
			err = fn()
			config.AfterTry()
			if err == nil {
				run <- nil
				return
			}
			// stop retry when overload
			if errors.Is(err, ErrorOverload) {
				overload <- struct{}{}
				return
			}
			// check whether to retry
			if config.RetryChecker != nil {
				needRetry := config.RetryChecker(err)
				if !needRetry {
					abort <- struct{}{}
					return
				}
			}
			if config.Strategy != nil {
				interval := config.Strategy(i + 1)
				<-time.After(interval)
			}
		}
		run <- err
	}()
	select {
	case <-ctx.Done():
		// context deadline exceed
		return ErrorContextDeadlineExceed
	case <-timer.C:
		// timeout
		return ErrorTimeout
	case <-abort:
		// caller abort
		return ErrorAbort
	case <-overload:
		// downstream overload
		return ErrorOverload
	case msg := <-panicInfoChan:
		// panic occurred
		if !config.RecoverPanic {
			panic(msg)
		}
		runErr = fmt.Errorf("panic occurred=%s", msg)
	case e := <-run:
		// normal run
		if e != nil {
			runErr = fmt.Errorf("retry failed, err=%w", e)
		}
	}
	return runErr
}
