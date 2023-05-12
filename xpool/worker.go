package xpool

import (
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

var (
	workerPool sync.Pool
)

// init
func init() {
	workerPool.New = newWorker
}

type worker struct {
	pool *pool
}

// newWorker
func newWorker() interface{} {
	return &worker{}
}

// run
func (w *worker) run() {
	go func() {
		for {
			var (
				t *task
			)

			w.pool.taskLock.Lock()
			if w.pool.taskHead != nil {
				t = w.pool.taskHead
				w.pool.taskHead = w.pool.taskHead.next
				atomic.AddInt32(&w.pool.taskCount, -1)
			}

			if t == nil {
				// 没有任务需要执行
				w.close()
				w.pool.taskLock.Unlock()
				w.recycle()
				return
			}
			w.pool.taskLock.Unlock()

			func() {
				defer func() {
					if r := recover(); r != nil {
						msg := fmt.Sprintf("POOL: panic in pool: %v: %s", r, debug.Stack())
						fmt.Println(msg)
					}
				}()

				// 执行
				t.f()
			}()

			// 回收使用
			t.recycle()
		}
	}()
}

func (w *worker) close() {
	w.pool.decWorkerCount()
}

func (w *worker) clear() {
	w.pool = nil
}

func (w *worker) recycle() {
	w.clear()
	workerPool.Put(w)
}
