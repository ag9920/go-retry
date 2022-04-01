# go-retry

A Stateless Golang Library For Retry Mechanism.

`go get github.com/ag9920/go-retry`

## Feature Supported
- Max retry times
- Abort execution
- Total timeout
- Retry interval after failed call
- Choose whether to throw panic
- Hooks

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
		goretry.WithRetryInterval(time.Second),
		goretry.WithThrowPanic(true))
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