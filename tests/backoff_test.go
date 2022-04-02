package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	goretry "github.com/ag9920/go-retry"
)

func Test_BackoffFibonacci(t *testing.T) {
	ctx := context.Background()
	definedErr := errors.New("err from outputErr")
	tsArray := []int64{}
	targetTimes := 6
	var retryTimes int
	err := goretry.Do(ctx, func() error {
		retryTimes++
		tsArray = append(tsArray, time.Now().Unix())
		return definedErr
	}, goretry.WithBackOffStrategy(goretry.StrategyFibonacci, time.Second), goretry.WithMaxRetryTimes(targetTimes))
	if err == nil {
		t.Errorf("should return err")
	}
	if !errors.Is(err, definedErr) {
		t.Errorf("returned err is not the same as defined, err=%v", err)
	}
	if retryTimes != targetTimes || len(tsArray) != targetTimes {
		t.Errorf("retry times not expected")
	}
	// fibonacci wait time series: 1,1,2,3,5,8
	for idx := range tsArray {
		if idx == 0 {
			continue
		}
		if tsArray[idx]-tsArray[idx-1] != int64(fibonacciNumber(idx)) {
			t.Errorf("not follow fibonacci backoff, gap=%v, tsArray=%v", tsArray[idx]-tsArray[idx-1], tsArray)
		}
	}
}

func fibonacciNumber(n int) int {
	if n == 0 || n == 1 {
		return n
	}
	return fibonacciNumber(n-1) + fibonacciNumber(n-2)
}
