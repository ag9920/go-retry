package tests

import (
	"context"
	"errors"
	"testing"

	goretry "github.com/ag9920/go-retry"
)

func Test_SimpleRetry(t *testing.T) {
	ctx := context.Background()
	definedErr := errors.New("err from outputErr")
	var retryTimes int
	err := goretry.Do(ctx, func() error {
		retryTimes++
		return definedErr
	})
	if err == nil {
		t.Errorf("should return err")
	}
	if !errors.Is(err, definedErr) {
		t.Errorf("returned err is not the same as defined, err=%v", err)
	}
	if retryTimes != goretry.DefaultMaxRetryTimes {
		t.Errorf("retry times not expected")
	}
}

func Test_SimpleRetryWithAbortErr(t *testing.T) {
	ctx := context.Background()
	definedErr := errors.New("err from outputErr")
	var retryTimes int
	err := goretry.Do(ctx, func() error {
		retryTimes++
		if retryTimes == 2 {
			// return abort err to stop retry, use default RetryChecker
			return goretry.ErrorAbort
		}
		return definedErr
	})
	if err == nil {
		t.Errorf("should return err")
	}
	if err != goretry.ErrorAbort {
		t.Errorf("returned err is not abort, err=%v", err)
	}
	if retryTimes != 2 {
		t.Errorf("retry times not expected")
	}
}
