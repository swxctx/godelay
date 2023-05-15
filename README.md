### Xpool
使用Golang开发的协程池

#### 思想
使用sync.pool实现任务添加、任务执行pool，使用Go协程进行处理

#### 使用方法

```go
func TestGo(t *testing.T) {
	for i := 0; i < 10; i++ {
		Go(func() {
			fmt.Printf("xpool %d\n", time.Now().UnixNano())
		})
	}
	select {}
}
```
