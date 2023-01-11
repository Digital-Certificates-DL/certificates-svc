package service

import (
	"context"
	"fmt"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	log := logan.New()

	ctx, _ := context.WithCancel(context.Background())

	ch := make(chan int, 10)
	j := 0
	for i := 0; i < 10; i++ {
		fmt.Println("start   ", i)
		go running.UntilSuccess(ctx, log, "test", func(ctx context.Context) (bool, error) {

			test, err := func() (bool, error) {
				if j > 114 {
					ch <- j
					return true, nil
				}

				return false, err
			}()
			j++
			fmt.Println(j)
			return test, err
		}, time.Millisecond*100, time.Millisecond*150) //todo move config file
	}

	close(ch)
	fmt.Println("close")
}
