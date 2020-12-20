package tablemgr

import (
	"fmt"
	"testing"
	"time"
)

func Test_Tick(t *testing.T) {
	ticker := time.Tick(50 * time.Millisecond)
	now := time.Now()
	for range ticker {
		fmt.Println("time diff:", time.Since(now).Milliseconds())
	}
}
