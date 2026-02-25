package metron

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/maypok86/otter/v2"
)

type slidingWindow struct {
	cache  *otter.Cache[string, []int64]
	key    string
	limit  int
	window time.Duration
}

func (sw *slidingWindow) wait(ctx context.Context) error {
	var waitDur time.Duration

	sw.cache.Compute(sw.key, func(state []int64, found bool) ([]int64, otter.ComputeOp) {
		now := time.Now()
		cutoff := now.Add(-sw.window).UnixNano()

		start := 0
		for start < len(state) && state[start] <= cutoff {
			start++
		}

		active := state[start:]

		if len(active) < sw.limit {
			newState := make([]int64, len(active)+1)
			copy(newState, active)
			newState[len(active)] = now.UnixNano()

			return newState, otter.WriteOp
		}

		expiry := time.Unix(0, active[0]).Add(sw.window)
		waitDur = expiry.Sub(now)

		newState := make([]int64, len(active))
		copy(newState, active[1:])
		newState[len(active)-1] = expiry.UnixNano()

		return newState, otter.WriteOp
	})

	if waitDur <= 0 {
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(waitDur):
		return nil
	}
}

type rateLimiter struct {
	cache     *otter.Cache[string, []int64]
	cachePath string
	perMin    *slidingWindow
	perDay    *slidingWindow
}

func newRateLimiter(cachePath string) (*rateLimiter, error) {
	c, err := otter.New[string, []int64](&otter.Options[string, []int64]{
		MaximumSize: 2,
	})
	if err != nil {
		return nil, err
	}

	if cachePath != "" {
		if err = otter.LoadCacheFromFile(c, cachePath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}

	return &rateLimiter{
		cache:     c,
		cachePath: cachePath,
		perMin:    &slidingWindow{cache: c, key: "per_min", limit: 30, window: time.Minute},
		perDay:    &slidingWindow{cache: c, key: "per_day", limit: 10_000, window: 24 * time.Hour},
	}, nil
}

func (rl *rateLimiter) Wait(ctx context.Context) error {
	if err := rl.perMin.wait(ctx); err != nil {
		return err
	}

	return rl.perDay.wait(ctx)
}

func (rl *rateLimiter) Close() error {
	if rl.cachePath == "" {
		return nil
	}

	return otter.SaveCacheToFile(rl.cache, rl.cachePath)
}
