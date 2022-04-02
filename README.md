# go-retry

A Stateless Golang Library For Retry Mechanism

- Written in vanilla Go, no dependencies
- Max retry times
- Abort execution
- Timeout
- Backoff strategy
- Choose whether to throw panic
- Hooks

You can easily integrate go-retry into your project. All features tested.

## Installation

When used with Go modules, use the following import path

`go get github.com/ag9920/go-retry`


## Usage

go-retry provides a simple function signature. You need to pass the RetryFunc, and customize retry mechanism with the help of `Option` if needed.

`type RetryFunc func() error`

`func Do(ctx context.Context, fn RetryFunc, opts ...Option) error`

Suppose the function you want to add retry support has the following signature:

```
func GetData() error   // call rpc to query data from downstream
```


**Example 1: Simple retry with default config**

A basic retry with default max retry times and timeout. Check the default retry config in `config.go`. 

```
import (
    goretry "github.com/ag9920/go-retry"
)

err := goretry.Do(ctx, GetData())

if err != nil {
    // meaning GetData() already succeed in the last retry.
} else {
    // retry failed, need to handle business logic accordingly.
}

```

goretry will keep executing `GetData` until the returned err of `GetData()` become nil or the maximum retry times has been reached, which by default is 3).


**Example 2: Change retry condition with RetryChecker**

`type RetryChecker func(err error) (needRetry bool)`

Sometimes we want to abort the execution if a certain error is detected. For example, the downstream was already completed unable to serve your request regardless of how many retries was conducted. We want to identify such conditions and abort retry. You could simply define a `RetryChecker` to meet the requirement.

```
var OverloadError = errors.New("already overload, should abort retry")

func isDownstreamOverload(err error) bool {
    return !errors.Is(err, OverloadError)
}

err := goretry.Do(ctx, GetData(), goretry.WithRetryChecker(isDownstreamOverload))

```

goretry will call `isDownstreamOverload` after each failed call to `GetData()` and pass the returned err. In the above scenario, `isDownstreamOverload` will check whether the `err` returned from `GetData()` is an overload error. If Yes, it will return a `false`, meaning there's no need to retry. The overall retry process will abort accordingly.


**Example 3: Change MaxRetryTimes && Timeout**

```
err := goretry.Do(ctx, GetData(), goretry.WithMaxRetryTimes(10), goretry.WithTimeout(5 * time.Second))

```


**Example 4: Add hooks before/after each retry**

```
func printTimestamp() {
    fmt.Println(time.Now().Unix())
}

err := goretry.Do(ctx, GetData(), goretry.WithAfterHook(printTimestamp))

```

**Example 5: Update backoff strategy**

`type Strategy func(times int) time.Duration`

By default, goretry will keep retry instantly after the last failed call. This may not be the most suitable way for your scenario. You could choose a different backoff strategy.

Currently, available strategies are the following

```

type BackoffStrategy int

const (
	StrategyConstant BackoffStrategy = iota
	StrategyLinear
	StrategyFibonacci
)
```

Add your strategy to Options. 

```
baseDuration := 10*time.Millisecond
err := goretry.Do(ctx, GetData(), goretry.WithBackOffStrategy(StrategyFibonacci, baseDuration))
```

If all provided strategies can't satisfy your needs, you could also specify your own implementation and use `WithCustomBackOffStrategy`

```

err := goretry.Do(ctx, GetData(), goretry.WithCustomBackOffStrategy(
    func(times int) time.Duration {
		if times > 1 {
            return 50*time.Millisecond
        }
        return time.Millisecond
	},
))

```

## Demo

Check the `example` foler to see a mock business demo.

```

import (
	"context"
	"log"
	"time"

	goretry "github.com/ag9920/go-retry"
)

type Item struct {
	ID    string
	Name  string
	Price string
}

type ItemRepo interface {
	Query(id string) (Item, error)
}

type ItemService struct {
	repo ItemRepo
}

func (srv ItemService) GetItem(ctx context.Context, id string) (*Item, error) {
	var err error
	var item *Item
	// wrapper function
	retryFn := func() error {
		// send request to downstream to get data
		result, err := srv.repo.Query(id)
		if err == nil {
			item = &result
			return nil
		}
		if err.Error() == "overload" {
			return goretry.ErrorAbort // use prefined abort err to stop execution
		}
		return err
	}
	err = goretry.Do(ctx, retryFn,
		goretry.WithMaxRetryTimes(3),
		goretry.WithTimeout(500*time.Microsecond),
		goretry.WithBackOffStrategy(goretry.StrategyFibonacci, 10*time.Microsecond),
		goretry.WithRecoverPanic())

	if err != nil {
		log.Default().Printf("retry failed, id=%s, err=%v", id, err)
		return nil, err
	}
	return item, nil
}
```


## Contribution

If you want to contribute resources to go-retry, Pull Requests are welcomed!