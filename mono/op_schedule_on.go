package mono

import (
	"context"

	"github.com/jjeffcaii/reactor-go"
	"github.com/jjeffcaii/reactor-go/internal"
	"github.com/jjeffcaii/reactor-go/scheduler"
)

type monoScheduleOn struct {
	source rs.RawPublisher
	sc     scheduler.Scheduler
}

func (m *monoScheduleOn) SubscribeWith(ctx context.Context, s rs.Subscriber) {
	actual := internal.NewCoreSubscriber(ctx, s)
	m.sc.Worker().Do(func() {
		m.source.SubscribeWith(ctx, actual)
	})
}

func newMonoScheduleOn(s rs.RawPublisher, sc scheduler.Scheduler) *monoScheduleOn {
	return &monoScheduleOn{
		source: s,
		sc:     sc,
	}
}
