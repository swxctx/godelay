package godelay

import (
	"sync"
	"sync/atomic"
)

type Pool interface {
	// SetCap 设置goroutine容量
	SetCap(cap int32)
	// Go 执行方法
	Go(f func())
	// WorkerCount 获取执行中的任务
	WorkerCount() int32
}

var (
	taskPool sync.Pool
)

func init() {
	taskPool.New = newTask
}

type task struct {
	f func()

	next *task
}

func newTask() interface{} {
	return &task{}
}

func (t *task) clear() {
	t.f = nil
	t.next = nil
}

func (t *task) recycle() {
	t.clear()
	taskPool.Put(t)
}

type pool struct {
	// goroutine 最大支持数量
	cap int32

	// linked list of tasks
	taskHead  *task
	taskTail  *task
	taskLock  sync.Mutex
	taskCount int32

	// 执行中的任务数量
	workerCount int32
}

// newPool
func newPool(cap int32) Pool {
	p := &pool{
		cap: cap,
	}
	return p
}

// SetCap
func (p *pool) SetCap(cap int32) {
	atomic.StoreInt32(&p.cap, cap)
}

func (p *pool) Go(f func()) {
	// 从pool获取
	t := taskPool.Get().(*task)
	t.f = f

	p.taskLock.Lock()

	if p.taskHead == nil {
		p.taskHead = t
		p.taskTail = t
	} else {
		p.taskTail.next = t
		p.taskTail = t
	}
	p.taskLock.Unlock()

	atomic.AddInt32(&p.taskCount, 1)

	// 有1个以上的任务待执行 && 容器没有满 || 当前没有任务在执行
	if (atomic.LoadInt32(&p.taskCount) >= 1 && p.WorkerCount() < atomic.LoadInt32(&p.cap)) || p.WorkerCount() == 0 {
		p.incWorkerCount()
		w := workerPool.Get().(*worker)
		w.pool = p
		w.run()
	}
}

func (p *pool) WorkerCount() int32 {
	return atomic.LoadInt32(&p.workerCount)
}

func (p *pool) incWorkerCount() {
	atomic.AddInt32(&p.workerCount, 1)
}

func (p *pool) decWorkerCount() {
	atomic.AddInt32(&p.workerCount, -1)
}
