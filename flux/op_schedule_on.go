package flux

import (
	"context"

	"github.com/jjeffcaii/reactor-go"
	"github.com/jjeffcaii/reactor-go/scheduler"
)

type fluxSubscribeOn struct {
	source rs.RawPublisher
	sc     scheduler.Scheduler
}

func (p *fluxSubscribeOn) SubscribeWith(ctx context.Context, s rs.Subscriber) {
	p.sc.Worker().Do(func() {
		p.source.SubscribeWith(ctx, s)
	})
}

func newFluxSubscribeOn(source rs.RawPublisher, sc scheduler.Scheduler) *fluxSubscribeOn {
	return &fluxSubscribeOn{
		source: source,
		sc:     sc,
	}
}
