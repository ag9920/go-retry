package tests

import (
	"context"
	"errors"
	"testing"

	goretry "github.com/ag9920/go-retry"
)

func Test_Checker(t *testing.T) {
	ctx := context.Background()

	retryErr := errors.New("need retry err")
	notRetryErr := errors.New("no need to retry")

	var retryTimes int
	err := goretry.Do(ctx, func() error {
		retryTimes++
		if retryTimes == 1 {
			return retryErr
		}
		return notRetryErr
	}, goretry.WithRetryChecker(
		func(err error) (needRetry bool) {
			return err != notRetryErr
		}))
	if err == nil {
		t.Errorf("should return err")
	}
	if err != goretry.ErrorAbort {
		t.Errorf("returned err is not the same as abort, err=%v", err)
	}
	if retryTimes != 2 {
		t.Errorf("retry times not expected")
	}
}
