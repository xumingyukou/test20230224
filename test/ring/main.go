package main

import (
	"container/ring"
	"fmt"
)

func main() {
	r := ring.New(10)
	fmt.Println(r, r.Len())
	//for i, now := 0, r.Next(); i < r.Len(); now, i = now.Next(), i+1 {
	//	now.Value = i + 10
	//}
	x := r.Next()
	for i := 0; i < 10; i++ {
		x.Value = i
		x = x.Prev()
	}
	fmt.Println(x.Value, 11111)
	// 遍历链表的值
	for i, now := 0, r.Next(); i < r.Len(); now, i = now.Next(), i+1 {
		fmt.Printf("%#v => %#v, ", i, now.Value)
	}
	fmt.Println(r.Value)
	fmt.Println(r.Prev().Value)
	fmt.Println(r.Next().Value)
}
