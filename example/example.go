package example

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
