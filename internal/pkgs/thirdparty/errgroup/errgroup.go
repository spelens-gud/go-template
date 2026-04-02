package errgroup

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

// 基于 Kratos errgroup
// 借鉴 golang.org/x/sync/errgroup进行增强：panic 恢复、任务接收 context、SetLimit/TryGo

// Group 是一组执行子任务的 goroutine 的集合，任一个返回非 nil 错误会取消 context 并作为 Wait 的返回值。
// 零值可用，且不会在错误时取消。
type Group struct {
	err        error                                // 错误
	wg         sync.WaitGroup                       // 等待所有任务完成
	errOnce    sync.Once                            // 错误只返回一次
	workerOnce sync.Once                            // 任务接收 context 只执行一次
	ch         chan func(ctx context.Context) error // 任务接收 context
	chs        []func(ctx context.Context) error    // 任务接收 context
	chsMu      sync.Mutex                           // 锁
	sem        chan struct{}                        // 可选：SetLimit 时用于限制并发数(信号量)
	ctx        context.Context                      // 可选：任务接收 context
	cancel     func()                               // 可选：任务取消
}

// WithCancel 创建一个 Group 及从 ctx 派生的可取消 context.
// 任一个 Go 传入的函数返回非 nil 错误或 Wait 返回时，该派生 context 会被取消.
func WithCancel(ctx context.Context) *Group {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{ctx: ctx, cancel: cancel}
}

// do 并发运行一个任务
func (g *Group) do(f func(ctx context.Context) error) {
	defer func() {
		if g.sem != nil {
			<-g.sem
		}
		g.wg.Done()
	}()

	ctx := g.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	var err error
	// 自动捕获
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 64<<10)
			buf = buf[:runtime.Stack(buf, false)]
			err = fmt.Errorf("errgroup: panic recovered: %s\n%s", r, buf)
		}
		if err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}()
	err = f(ctx)
}

// SetLimit 将同时运行的 goroutine 数限制为 n.
// Go 调用在达到上限时会阻塞.
// 应在任何 Go 调用之前调用；与 GOMAXPROCS 二选一使用.
func (g *Group) SetLimit(n int) {
	if n <= 0 {
		panic("errgroup: SetLimit n must be positive")
	}
	if g.sem != nil {
		return
	}
	g.sem = make(chan struct{}, n)
}

// GOMAXPROCS 使用固定 n 个 worker 的池子执行任务(任务入队，worker 取任务执行).
// 与 SetLimit 二选一使用；SetLimit 为“最多 n 个并发”，GOMAXPROCS 为“n 个 worker 排队”.
func (g *Group) GOMAXPROCS(n int) {
	if n <= 0 {
		panic("errgroup: GOMAXPROCS must be positive")
	}
	g.workerOnce.Do(func() {
		g.ch = make(chan func(context.Context) error, n)
		for i := 0; i < n; i++ {
			go func() {
				for fn := range g.ch {
					g.do(fn)
				}
			}()
		}
	})
}

// Go 在新 goroutine 中执行 f。第一个返回非 nil 错误的会取消 group 的 context，该错误由 Wait 返回.
// 若调用了 SetLimit，在达到上限时 Go 会阻塞直到有空位.
func (g *Group) Go(f func(ctx context.Context) error) {
	if g.sem != nil {
		g.sem <- struct{}{}
	}
	g.wg.Add(1)
	if g.ch != nil {
		select {
		case g.ch <- f:
			return
		default:
			g.chsMu.Lock()
			g.chs = append(g.chs, f)
			g.chsMu.Unlock()
			return
		}
	}
	go g.do(f)
}

// TryGo 尝试在新 goroutine 中执行 f。若已设置 SetLimit 且当前已达上限则不启动并返回 false，否则返回 true.
// 未设置 SetLimit 时行为同 Go，始终返回 true.
func (g *Group) TryGo(f func(ctx context.Context) error) bool {
	if g.sem != nil {
		select {
		case g.sem <- struct{}{}:
		default:
			return false
		}
	}
	g.wg.Add(1)
	if g.ch != nil {
		select {
		case g.ch <- f:
			return true
		default:
			g.chsMu.Lock()
			g.chs = append(g.chs, f)
			g.chsMu.Unlock()
			return true
		}
	}
	go g.do(f)
	return true
}

// Wait 阻塞直到所有通过 Go 启动的函数返回，并返回其中第一个非 nil 错误(若有).
func (g *Group) Wait() error {
	if g.ch != nil {
		g.chsMu.Lock()
		pending := g.chs
		g.chs = nil
		g.chsMu.Unlock()
		for _, f := range pending {
			g.ch <- f
		}
		g.wg.Wait()
		close(g.ch)
		g.ch = nil
	} else {
		g.wg.Wait()
	}
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}
