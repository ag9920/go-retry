package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	goretry "github.com/ag9920/go-retry"
)

func Test_Hook(t *testing.T) {
	ctx := context.Background()

	arr := []int{}

	hookFn := func() {
		arr = append(arr, int(time.Now().Unix()))
	}

	var retryTimes int
	err := goretry.Do(ctx, func() error {
		retryTimes++
		return errors.New("output err")
	}, goretry.WithAfterHook(hookFn))
	if err == nil {
		t.Errorf("should return err")
	}
	if retryTimes != goretry.DefaultMaxRetryTimes {
		t.Errorf("retry times not expected")
	}
	if len(arr) != goretry.DefaultMaxRetryTimes {
		t.Errorf("after hook not called as expected")
	}
}
