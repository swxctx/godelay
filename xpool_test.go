package godelay

import (
	"fmt"
	"testing"
	"time"
)

func TestGo(t *testing.T) {
	for i := 0; i < 10; i++ {
		Go(func() {
			fmt.Printf("xpool %d\n", time.Now().UnixNano())
		})
	}
	select {}
}
