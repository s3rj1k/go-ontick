package ontick

import (
	"context"
	"runtime"
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

	semaphore chan struct{}

	key T // key value for context.WithValue() to get TickTime
}

func New[T comparable](ctx context.Context, duration time.Duration, concurrency int, key T) *OnTick[T] {
	ctx, cancel := context.WithCancel(ctx)

	cfg := &OnTick[T]{
		ticker: time.NewTicker(duration),
		ctx:    ctx,
		cancel: cancel,
		key:    key,
	}

	if concurrency < 1 {
		cfg.semaphore = make(chan struct{}, runtime.NumCPU()+1)
	} else {
		cfg.semaphore = make(chan struct{}, concurrency)
	}

	return cfg
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
						et.semaphore <- struct{}{}
						f(context.WithValue(et.ctx, et.key, t))
						<-et.semaphore
					case <-et.ctx.Done():
						return
					}
				}
			}(f)
		}
	})
}
