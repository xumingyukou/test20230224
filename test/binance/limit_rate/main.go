package main

import (
	"fmt"
	"time"
)

func main() {
	r := "24h1s500ms1us1ns"
	fmt.Println(r[1:4])
	a, err := time.ParseDuration(r)
	fmt.Println(a, a.Seconds(), int64(a.Seconds()), err)
	var d *int
	var b, c int
	d = &b
	d = &c
	fmt.Println(d)
}
