package transform

import (
	"fmt"
	"testing"
	"time"
)

func TestSetDate(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().In(loc)
	fmt.Println(time.Date(now.Year(), now.Month()+7, 0, 0, 0, 0, 0, now.Location()))
}

func TestGetDate(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2022, 9, 30, 5, 1, 0, 0, loc)
	res := GetThisWeek(now, 5, 16)
	fmt.Println("this week:", res)
	res = GetNextWeek(now, 5, 16)
	fmt.Println("next week:", res)
	res = GetThisMonth(now, 5, 16)
	fmt.Println("this month:", res)
	res = GetNextMonth(now, 5, 16)
	fmt.Println("next month:", res)
	res = GetThisQuarter(now, 5, 16)
	fmt.Println("this quarter:", res)
	res = GetNextQuarter(now, 5, 16)
	fmt.Println("next quarter:", res)

}

func TestGetLocationTime(t *testing.T) {
	fmt.Println(time.Now().In(BeiJingLoc).Format("20060102150405.000000-0700"))
	fmt.Println(time.Now().In(LondonLoc).Format("20060102150405.000000-0700"))
	fmt.Println(time.Now().UTC().Format("20060102150405.000000-0700"))
}
