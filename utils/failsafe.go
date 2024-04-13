package utils

import "context"

func FailSafeNoBreak(ctx context.Context, runFunc func(context.Context), fallbackFunc func(context.Context, error)) {
	defer func() {
		if rval := recover(); rval != nil {
			var err error
			err, _ = rval.(error)
			fallbackFunc(ctx, err)
		}
	}()
	runFunc(ctx)
}

func FailSafeNoBreakWithError(ctx context.Context, runFunc func(context.Context) error, fallbackFunc func(context.Context, error) error, tags map[string]string) {
	defer func() {
		if rval := recover(); rval != nil {
			var err error
			err, _ = rval.(error)
			_ = fallbackFunc(ctx, err)
		}
	}()
	err := runFunc(ctx)
	if err != nil {
		_ = fallbackFunc(ctx, err)
	}
}
