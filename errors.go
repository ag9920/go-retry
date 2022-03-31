package goretry

import "errors"

var (
	ErrorOverload              = errors.New("overload err")
	ErrorAbort                 = errors.New("stop retry")
	ErrorTimeout               = errors.New("retry timeout")
	ErrorContextDeadlineExceed = errors.New("context deadline exceeded")
	ErrorEmptyRetryFunc        = errors.New("empty retry function")
)
