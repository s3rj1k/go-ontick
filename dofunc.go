package ontick

import (
	"context"
	"sync"
	"time"
)

func DoFunc[T comparable](ctx context.Context, wg *sync.WaitGroup, duration time.Duration, key T, f func(context.Context)) {
	ticker := time.NewTicker(duration)

	wg.Add(1)

	go func() {
		defer func() {
			ticker.Stop()
			wg.Done()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				f(context.WithValue(ctx, key, t))
			}
		}
	}()
}
