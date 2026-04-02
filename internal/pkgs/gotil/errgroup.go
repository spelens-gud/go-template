package gotil

import (
	"context"

	"{{.ProjectName}}/internal/pkgs/thirdparty/errgroup"

	"git.bestfulfill.tech/devops/go-core/kits/kcontext"
)

// GoErrGroup 并发执行 fs，任一个返回错误则取消并返回该错误；全部成功才返回 nil。
func GoErrGroup(ctx context.Context, fs ...func(context.Context) (err error)) error {
	if len(fs) == 0 {
		return nil
	}
	if len(fs) == 1 {
		return fs[0](ctx)
	}
	eg := errgroup.WithCancel(ctx)
	for _, f := range fs {
		eg.Go(f)
	}
	return eg.Wait()
}

// GoEGWithLimit 与 GoErrGroup 相同，但最多同时运行 limit 个任务（借鉴官方 errgroup.SetLimit）。
func GoEGWithLimit(ctx context.Context, limit int, fs ...func(context.Context) (err error)) error {
	if len(fs) == 0 {
		return nil
	}
	if len(fs) == 1 {
		return fs[0](ctx)
	}
	eg := errgroup.WithCancel(ctx)
	eg.SetLimit(limit)
	for _, f := range fs {
		eg.Go(f)
	}
	return eg.Wait()
}

// GetAsync 获取异步结果
func GetAsync(ctx context.Context, f func(context.Context) (any, error)) (ret any, err error) {
	waitingDone := ctx.Done()
	retChan := make(chan any)

	// 异步获取
	go func(ctx context.Context) {
		defer close(retChan)
		ret, err = f(ctx)
		select {
		case <-waitingDone:
		default:
			if err != nil {
				retChan <- err
			} else {
				retChan <- ret
			}
		}
	}(kcontext.Detach(ctx))

	select {
	case <-waitingDone:
		err = ctx.Err()
		return
	case ret = <-retChan:
		if r, ok := ret.(error); ok {
			ret = nil
			err = r
		}
	}
	return ret, err
}
