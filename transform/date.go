package transform

import (
	"time"
)

type CONTRACTTYPE string

var (
	BeiJingLoc *time.Location
	LondonLoc  *time.Location
)

func init() {
	var err error
	BeiJingLoc, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	LondonLoc, err = time.LoadLocation("Europe/London")
	if err != nil {
		panic(err)
	}
}

const (
	THISWEEK    CONTRACTTYPE = "this_week"
	NEXTWEEK    CONTRACTTYPE = "next_week"
	THISMONTH   CONTRACTTYPE = "this_month"
	NEXTMONTH   CONTRACTTYPE = "next_month"
	THISQUARTER CONTRACTTYPE = "this_quarter"
	NEXTQUARTER CONTRACTTYPE = "next_quarter"
)

// GetDate 获取060102格式的string时间
func GetDate(contract CONTRACTTYPE) string {
	// 当周、次周、当季度、次季度
	now := time.Now().In(BeiJingLoc)
	switch contract {
	case THISWEEK:
		return GetThisWeek(now, 5, 16).Format("060102")
	case NEXTWEEK:
		return GetNextWeek(now, 5, 16).Format("060102")
	case THISMONTH:
		return GetThisMonth(now, 5, 16).Format("060102")
	case NEXTMONTH:
		return GetNextMonth(now, 5, 16).Format("060102")
	case THISQUARTER:
		return GetThisQuarter(now, 5, 16).Format("060102")
	case NEXTQUARTER:
		return GetNextQuarter(now, 5, 16).Format("060102")
	default:
		return ""
	}
}

func GetThisWeek(now time.Time, week, hour int64) time.Time {
	/*
		当周
		now：现在的时间
		week：结算日的星期
		hour：结算日的小时
	*/
	var (
		day     = week - int64(now.Weekday())
		nowHour = hour - int64(now.Hour())
	)
	if day > 0 || (day == 0 && nowHour > 0) {
		now = now.AddDate(0, 0, int(day))
	} else {
		now = now.AddDate(0, 0, int(7+day))
	}
	return time.Date(now.Year(), now.Month(), now.Day(), int(hour), 0, 0, 0, now.Location())
}

func GetNextWeek(now time.Time, week, hour int64) time.Time {
	/*
		次周
		now：现在的时间
		week：结算日的星期
		hour：结算日的小时
	*/
	return GetThisWeek(now, week, hour).AddDate(0, 0, 7)
}

func GetThisMonth(now time.Time, week, hour int64) time.Time {
	/*
		当月
		now：现在的时间
		week：结算日的星期
		hour：结算日的小时
	*/
	lastDayOfMonth := time.Date(now.Year(), now.Month()+1, 0, int(hour), 0, 0, 0, now.Location())
	day := week - int64(lastDayOfMonth.Weekday())
	if day > 0 {
		day -= 7
	}
	date := time.Date(now.Year(), now.Month(), lastDayOfMonth.Day()+int(day), int(hour), 0, 0, 0, now.Location())
	if now.Sub(date) > 0 {
		return GetThisMonth(time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()), week, hour)
	}
	return date
}

func GetNextMonth(now time.Time, week, hour int64) time.Time {
	/*
		次月
		now：现在的时间
		week：结算日的星期
		hour：结算日的小时
	*/
	thisMonth := GetThisMonth(now, week, hour)
	if thisMonth.Month() == now.Month()+1 {
		return GetThisMonth(time.Date(now.Year(), now.Month()+2, 1, 0, 0, 0, 0, now.Location()), week, hour)
	}
	return GetThisMonth(time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()), week, hour)
}

func GetThisQuarter(now time.Time, week, hour int64) time.Time {
	/*
		当季度
		now：现在的时间
		week：结算日的星期
		hour：结算日的小时
	*/
	var (
		month = now.Month()
		date  time.Time
	)
	if month%3 == 0 {
		lastDayOfQuarter := time.Date(now.Year(), now.Month()+1, 0, int(hour), 0, 0, 0, now.Location())
		day := int(week) - int(lastDayOfQuarter.Weekday())
		if day > 0 {
			day -= 7
		}
		date = lastDayOfQuarter.AddDate(0, 0, day)
		if now.Sub(date) > 0 {
			return GetThisQuarter(now.AddDate(0, 1, 0), week, hour)
		}
	} else {
		lastDayOfQuarter := time.Date(now.Year(), now.Month()+4-now.Month()%3, 0, 0, 0, 0, 0, now.Location())
		day := int(week) - int(lastDayOfQuarter.Weekday())
		if day > 0 {
			day -= 7
		}
		date = time.Date(lastDayOfQuarter.Year(), lastDayOfQuarter.Month(), lastDayOfQuarter.Day()+day, int(hour), 0, 0, 0, now.Location())
	}
	return date
}

func GetNextQuarter(now time.Time, week, hour int64) time.Time {
	// 次季度
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
	thisQuarter := GetThisQuarter(nowDate, week, hour)
	if thisQuarter.Month() == now.Month()+3 {
		return GetThisQuarter(time.Date(now.Year(), now.Month(), 1, now.Hour(), 0, 0, 0, now.Location()).AddDate(0, 4, 0), week, hour)
	} else {
		return GetThisQuarter(time.Date(now.Year(), now.Month(), 1, now.Hour(), 0, 0, 0, now.Location()).AddDate(0, 3, 0), week, hour)
	}
}
