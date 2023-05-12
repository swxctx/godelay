package xpool

import (
	"math"
)

var (
	// 默认协程池配置
	defaultPool Pool
)

func init() {
	defaultPool = newPool(math.MaxInt32)
}

// Go 执行func
func Go(f func()) {
	defaultPool.Go(f)
}

// SetCap 设置pool容量, 默认math.MaxInt32
func SetCap(cap int32) {
	defaultPool.SetCap(cap)
}

// WorkerCount 返回协程池正在执行的协程数量
func WorkerCount() int32 {
	return defaultPool.WorkerCount()
}
