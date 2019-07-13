package flux_test

import (
	"context"
	"log"
	"testing"

	rs "github.com/jjeffcaii/reactor-go"
	"github.com/jjeffcaii/reactor-go/flux"
	"github.com/jjeffcaii/reactor-go/scheduler"
	"github.com/stretchr/testify/assert"
)

func createSinker(totals int) func(context.Context, flux.Sink) {
	return func(ctx context.Context, sink flux.Sink) {
		for i := 0; i < totals; i++ {
			v := i
			sink.Next(v)
		}
		sink.Complete()
	}
}

func createSubscriber() rs.Subscriber {
	var su rs.Subscription
	return rs.NewSubscriber(
		rs.OnNext(func(i interface{}) {
			log.Println("next:", i)
			su.Request(1)
		}),
		rs.OnSubscribe(func(s rs.Subscription) {
			s.Request(1)
			su = s
		}),
	)
}

func TestFluxCreate_SubscribeOn(t *testing.T) {
	done := make(chan struct{})
	flux.Create(createSinker(18)).
		Map(func(i interface{}) interface{} {
			return i.(int) * 2
		}).
		SubscribeOn(scheduler.Elastic()).
		Subscribe(context.Background(),
			rs.OnNext(func(v interface{}) {
				log.Println("next:", v)
			}),
			rs.OnComplete(func() {
				close(done)
			}),
		)
	<-done
}

func TestCreate(t *testing.T) {
	const totals = 18
	sinker := createSinker(totals)
	s := createSubscriber()
	log.Println("---test1")
	flux.Create(sinker).SubscribeWith(context.Background(), s)
	log.Println("---test2")
	flux.Create(sinker).
		Map(func(i interface{}) interface{} {
			return i.(int) * 2
		}).
		SubscribeWith(context.Background(), s)

	log.Println("---test3")

	var su rs.Subscription

	flux.Create(sinker).
		Filter(func(i interface{}) bool {
			return i.(int) >= totals-1
		}).
		DoOnDiscard(func(i interface{}) {
			// TODO: auto
			su.Request(1)
		}).
		Map(func(i interface{}) interface{} {
			return i.(int) * 2
		}).
		Subscribe(context.Background(),
			rs.OnNext(func(i interface{}) {
				log.Println("next:", i)
				assert.Equal(t, (totals-1)*2, i.(int), "bad value")
				su.Request(1)
			}),
			rs.OnSubscribe(func(s rs.Subscription) {
				su = s
				su.Request(1)
			}),
		)
}
