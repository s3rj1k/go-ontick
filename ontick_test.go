package ontick_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/s3rj1k/go-ontick"
)

type TickKey int

const (
	duration = 100 * time.Millisecond
	sleep    = duration * 15
	key      = TickKey(42)
)

func TestOnTick(t *testing.T) {
	ctx := context.Background()

	ticker := ontick.New(ctx, duration, key)
	defer ticker.Stop()

	var wg sync.WaitGroup
	wg.Add(1)

	ticker.Do(func(ctx context.Context) {
		wg.Done()
	})

	wg.Wait()
}

func TestGetTickTimeFromContext(t *testing.T) {
	testTime := time.Now()

	tests := []struct {
		name     string
		ctx      context.Context
		wantTime time.Time
		wantZero bool
	}{
		{
			name:     "With tick time",
			ctx:      context.WithValue(context.Background(), key, testTime),
			wantTime: testTime,
			wantZero: false,
		},
		{
			name:     "Without tick time",
			ctx:      context.Background(),
			wantTime: time.Time{},
			wantZero: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			ticker := ontick.New(ctx, duration, key)
			defer ticker.Stop()

			time.Sleep(sleep)

			gotTime := ticker.GetTickTimeFromContext(tt.ctx)
			if !gotTime.Equal(tt.wantTime) {
				t.Errorf("GetTickTimeFromContext() gotTime = %v, want %v", gotTime, tt.wantTime)
			}

			if gotTime.IsZero() != tt.wantZero {
				t.Errorf("GetTickTimeFromContext() gotTime.IsZero() = %v, want %v", gotTime.IsZero(), tt.wantZero)
			}
		})
	}
}

func TestStopWithContext(t *testing.T) {
	ctx := context.Background()

	onTick := ontick.New(ctx, duration, key)
	tickCount := new(atomic.Int64)

	onTick.Do(func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				tickCount.Add(1)
			}
		}
	})

	time.Sleep(sleep)
	onTick.Stop()

	tickCountBefore := tickCount.Load()
	onTick.Wait()
	tickCountAfter := tickCount.Load()

	if tickCountAfter != tickCountBefore {
		t.Errorf("Expected no more ticks after Stop, but tick count changed from %d to %d", tickCountBefore, tickCountAfter)
	}
}

func TestStopWithoutContext(t *testing.T) {
	ctx := context.Background()

	onTick := ontick.New(ctx, duration, key)
	tickCount := new(atomic.Int64)

	onTick.Do(func(ctx context.Context) {
		tickCount.Add(1)
	})

	time.Sleep(sleep)
	onTick.Stop()

	tickCountBefore := tickCount.Load()
	onTick.Wait()
	tickCountAfter := tickCount.Load()

	if tickCountAfter-1 != tickCountBefore {
		t.Errorf("Expected no more ticks after Stop, but tick count changed from %d to %d", tickCountBefore, tickCountAfter)
	}
}
