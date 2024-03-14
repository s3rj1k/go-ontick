package ontick

import (
	"context"
	"sync"
	"time"
)

type OnTick[T comparable] struct {
	ticker *time.Ticker

	ctx    context.Context
	cancel context.CancelFunc

	onceDo   sync.Once
	onceStop sync.Once
	wg       sync.WaitGroup

	key T // key value for context.WithValue() to get TickTime
}

func New[T comparable](ctx context.Context, d time.Duration, key T) *OnTick[T] {
	ctx, cancel := context.WithCancel(ctx)

	return &OnTick[T]{
		ticker: time.NewTicker(d),
		ctx:    ctx,
		cancel: cancel,
		key:    key,
	}
}

func (et *OnTick[T]) GetTickTimeFromContext(ctx context.Context) time.Time {
	key := et.key

	tickTime, ok := ctx.Value(key).(time.Time)
	if !ok {
		return time.Time{}
	}

	return tickTime
}

func (et *OnTick[T]) Stop() {
	et.onceStop.Do(func() {
		et.ticker.Stop()
		et.cancel()
	})
}

func (et *OnTick[T]) Wait() {
	et.wg.Wait()
}

func (et *OnTick[T]) Do(funcs ...func(context.Context)) {
	et.onceDo.Do(func() {
		for _, f := range funcs {
			et.wg.Add(1)

			go func(f func(context.Context)) {
				defer et.wg.Done()

				for {
					select {
					case t := <-et.ticker.C:
						f(context.WithValue(et.ctx, et.key, t))
					case <-et.ctx.Done():
						return
					}
				}
			}(f)
		}
	})
}
