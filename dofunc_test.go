package ontick_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/s3rj1k/go-ontick"
)

func TestDoFunc(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	var (
		wg    sync.WaitGroup
		count atomic.Int64
	)

	ontick.DoFunc(ctx, &wg, 50*time.Millisecond, 42, func(ctx context.Context) {
		if _, ok := ctx.Value(42).(time.Time); ok {
			count.Add(1)
		}
	})

	wg.Wait()

	if count.Load() < 1 {
		t.Error("Expected at least one executions")
	}
}
