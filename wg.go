package concexample

import (
	"context"
	"sync"
	"time"

	"github.com/sourcegraph/conc"
)

const (
	workerLimit       = 1000
	workerMinDuration = time.Nanosecond * 10
	contextTimeout    = workerMinDuration * workerLimit
	cancelTimeout     = contextTimeout - workerMinDuration
)

func WithConcWgNoPanics() {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	wg := conc.NewWaitGroup()
	for i := 0; i < workerLimit; i++ {
		m := i
		wg.Go(func() {
			concJobNoPanic(ctx, workerMinDuration*time.Duration(m))
		})
	}
	_ = wg.WaitAndRecover()
	cancel()
}

func WithBuiltinWgNoPanics() {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	wg := &sync.WaitGroup{}
	for i := 0; i < workerLimit; i++ {
		m := i
		wg.Add(1)
		go builtinJobNoPanic(ctx, wg, workerMinDuration*time.Duration(m))
	}
	wg.Wait()
	cancel()
}

func WithConcWgRecovered() {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	wg := conc.NewWaitGroup()
	for i := 0; i < workerLimit; i++ {
		m := i
		wg.Go(func() {
			concJobWithPanic(ctx, workerMinDuration*time.Duration(m))
		})
	}
	<-time.NewTimer(cancelTimeout).C
	cancel()
	wg.WaitAndRecover()
}

func WithBuiltinWgRecovered() {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	wg := &sync.WaitGroup{}
	for i := 0; i < workerLimit; i++ {
		m := i
		wg.Add(1)
		go builtinJobWithPanic(ctx, wg, workerMinDuration*time.Duration(m))
	}
	<-time.NewTimer(cancelTimeout).C
	cancel()
	wg.Wait()
}

func concJobNoPanic(ctx context.Context, timeout time.Duration) {
	timer := time.NewTimer(timeout)
	select {
	case <-timer.C:
	case <-ctx.Done():
	}
}
func concJobWithPanic(ctx context.Context, timeout time.Duration) {
	timer := time.NewTimer(timeout)
	select {
	case <-timer.C:
	case <-ctx.Done():
		panic("context canceled")
	}
}

func builtinJobNoPanic(ctx context.Context, wg *sync.WaitGroup, timeout time.Duration) {
	timer := time.NewTimer(timeout)
	select {
	case <-timer.C:
		wg.Done()
	case <-ctx.Done():
		wg.Done()
	}
}
func builtinJobWithPanic(ctx context.Context, wg *sync.WaitGroup, timeout time.Duration) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	timer := time.NewTimer(timeout)
	select {
	case <-timer.C:
		wg.Done()
	case <-ctx.Done():
		wg.Done()
		panic("context canceled")
	}
}
