# go-retry

A Stateless Golang Library For Retry Mechanism.

You can easily integrate go-retry into your project. All features tested.

Check the default retry config in `config.go`. And feel free to customize config with options.

## Feature Supported
- Max retry times
- Abort execution
- Timeout
- Backoff strategy
- Choose whether to throw panic
- Hooks

## Getting Started

`go get github.com/ag9920/go-retry`


## Usage

```
import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	goretry "github.com/ag9920/go-retry"
)

func main() {
	ctx := context.Background()
	goretry.Do(ctx, randFn,
		goretry.WithMaxRetryTimes(10),
		goretry.WithBeforeHook(printTime),
		goretry.WithRecoverPanic(true))
}

func randFn() error {
	randNum := rand.Int31()
	if randNum < 15 {
		return errors.New("less")
	}
	if randNum < 20 {
		panic("less than 20")
	}
	return nil
}

func printTime() {
	fmt.Println(time.Now().Unix())
}

```


## Contribution

If you want to contribute resources to go-retry, Pull Requests are welcomed!