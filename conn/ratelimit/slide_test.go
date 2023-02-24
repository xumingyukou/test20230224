package ratelimit

import (
	"fmt"
	"testing"
	"time"
)

func TestSlideWindow(t *testing.T) {

	s := NewSlideWindow(10, 10, 50)
	for i := 0; i < 1000; i++ {
		fmt.Println(i, s.AllowN(10), s.count)
		time.Sleep(time.Second)
		cur := s.bucketChain.Prev()
		for j := 0; int64(j) < s.bucketCount; j++ {
			fmt.Printf("%#v\n", cur.Value.(bucketNode))
			cur = cur.Prev()
		}
		fmt.Println()
	}
}
