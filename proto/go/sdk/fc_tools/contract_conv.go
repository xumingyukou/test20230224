package fc_tools

import (
	"strings"
	"time"

	"github.com/warmplanet/proto/go/common"
)

type CONTRACTTYPE string

var (
	BeiJingLoc *time.Location
	LondonLoc  *time.Location
)

const (
	THISWEEK    CONTRACTTYPE = "this_week"
	NEXTWEEK    CONTRACTTYPE = "next_week"
	THISMONTH   CONTRACTTYPE = "this_month"
	NEXTMONTH   CONTRACTTYPE = "next_month"
	THISQUARTER CONTRACTTYPE = "this_quarter"
	NEXTQUARTER CONTRACTTYPE = "next_quarter"

	DEFAULT_WEEKDAY = 5
	DEFAULT_HOUR    = 16
)

type SettlingTime struct {
	Week int64
	Hour int64
}

type ContractDate string //合约日期

var (
	ExchangeSettlingTime = make(map[common.Exchange]*SettlingTime)
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

	ExchangeSettlingTime[common.Exchange_ALL] = new(SettlingTime)
	ExchangeSettlingTime[common.Exchange_ALL].Week = DEFAULT_WEEKDAY
	ExchangeSettlingTime[common.Exchange_ALL].Hour = DEFAULT_HOUR
}

// 根据合约日期获得合约类型
func Get_SymbolType_From_Date(date ContractDate, exch common.Exchange, is_coin bool) common.SymbolType {
	c_symbol_type := common.SymbolType_INVALID_TYPE

	st := GetExchangeSettlingTime(exch)

	//根据日期串，获取symbol type
	now := time.Now().In(BeiJingLoc)
	switch string(date) {
	case GetThisWeek(now, st.Week, st.Hour).Format("060102"):
		c_symbol_type = common.SymbolType_FUTURE_THIS_WEEK
		if is_coin {
			c_symbol_type = common.SymbolType_FUTURE_COIN_THIS_WEEK
		}
	case GetNextWeek(now, st.Week, st.Hour).Format("060102"):
		c_symbol_type = common.SymbolType_FUTURE_COIN_NEXT_WEEK
		if is_coin {
			c_symbol_type = common.SymbolType_FUTURE_COIN_NEXT_WEEK
		}
	case GetThisMonth(now, st.Week, st.Hour).Format("060102"):
		c_symbol_type = common.SymbolType_FUTURE_THIS_MONTH
		if is_coin {
			c_symbol_type = common.SymbolType_FUTURE_COIN_THIS_MONTH
		}
	case GetNextMonth(now, st.Week, st.Hour).Format("060102"):
		c_symbol_type = common.SymbolType_FUTURE_NEXT_MONTH
		if is_coin {
			c_symbol_type = common.SymbolType_FUTURE_COIN_NEXT_MONTH
		}
	case GetThisQuarter(now, st.Week, st.Hour).Format("060102"):
		c_symbol_type = common.SymbolType_FUTURE_THIS_QUARTER
		if is_coin {
			c_symbol_type = common.SymbolType_FUTURE_COIN_THIS_QUARTER
		}
	case GetNextMonth(now, st.Week, st.Hour).Format("060102"):
		c_symbol_type = common.SymbolType_FUTURE_NEXT_QUARTER
		if is_coin {
			c_symbol_type = common.SymbolType_FUTURE_COIN_NEXT_QUARTER
		}
	}
	return c_symbol_type
}

// 根据合约类型获取合约日期
func Contract_SymbolType_To_Date(st common.SymbolType) ContractDate {
	//交割合约返回合约结束日期，永续合约返回空串
	type_str := common.SymbolType_name[int32(st)]
	type_str = strings.Replace(type_str, "FUTURE_", "", 1)
	type_str = strings.Replace(type_str, "COIN_", "", 1)
	type_str = strings.ToLower(type_str)
	c_date := GetDate(CONTRACTTYPE(type_str))
	return ContractDate(c_date)
}

// GetDate 获取060102格式的string时间
func GetDate(contract CONTRACTTYPE) string {
	// 当周、次周、当月、次月、当季度、次季度
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

func GetExchangeSettlingTime(exch common.Exchange) *SettlingTime {
	st, ok := ExchangeSettlingTime[exch]

	if !ok {
		st = ExchangeSettlingTime[common.Exchange_ALL]
	}

	return st
}
