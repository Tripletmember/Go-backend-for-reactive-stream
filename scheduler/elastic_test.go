package scheduler_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jjeffcaii/reactor-go/scheduler"
)

func TestElastic(t *testing.T) {
	fmt.Println(scheduler.Elastic().Name())
	done := make(chan int)
	_ = scheduler.Elastic().Worker().Do(func() {
		done <- 1
	})
	assert.Equal(t, 1, <-done)
}

func TestNewElastic(t *testing.T) {
	const size = 100
	sc := scheduler.NewElastic(size)

	wg := sync.WaitGroup{}
	wg.Add(size * 2)
	for i := 0; i < size*2; i++ {
		assert.NoError(t, sc.Worker().Do(func() {
			wg.Done()
		}), "should not return error")
	}
	wg.Wait()

	assert.NoError(t, sc.Close())

	assert.Error(t, sc.Worker().Do(func() {
		// noop
	}), "should return error after closing scheduler")
}

func TestElasticBounded(t *testing.T) {
	assert.NotPanics(t, func() {
		const total = 1000
		var wg sync.WaitGroup
		wg.Add(total)
		worker := scheduler.ElasticBounded().Worker()
		for range [total]struct{}{} {
			err := worker.Do(func() {
				time.Sleep(10 * time.Millisecond)
				wg.Done()
			})
			assert.NoError(t, err)
		}
		wg.Wait()
	})
}
