package base

import (
	"clients/conn/ratelimit"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	R_BASIC_IP RateLimitType = "ip"
	//Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
)

// RateLimitType 限流类型
// 类型 + 时间 + 其他信息
// ip/account + duration + url/instId/symbol
type RateLimitType string

// RateLimitInfo 限流相关信息
type RateLimitInfo struct {
	BucketType   ratelimit.BUCKET_TYPE //令牌桶或者时间滑动窗口
	Type         RateLimitType
	LimitTimeSec float64 // 限流区间 单位s
	BucketCount  int64   // 时间窗口数量
	LimitCount   int64   // 限流数量
	Instance     ratelimit.LimitI
}

/*
RateLimitMgr
- 创建limitMap，用于binance全局的限制，
*/
type RateLimitMgr struct {
	bucketType ratelimit.BUCKET_TYPE
	limitMap   map[RateLimitType]RateLimitInfo

	lock sync.Mutex
}

func NewRateLimitMgr() *RateLimitMgr {
	r := &RateLimitMgr{
		bucketType: ratelimit.TOKENBUCKET, //使用令牌桶
		limitMap:   make(map[RateLimitType]RateLimitInfo),
	}
	return r
}

func (r *RateLimitMgr) CreateLimitItem(bucketType ratelimit.BUCKET_TYPE, item RateLimitType, count int64) error {
	if count < 1 {
		return errors.New(fmt.Sprint(bucketType, item, "rate limit capacity err:", count))
	}
	limitTimeSec, err := ParseRateLimitTime(item)
	if err != nil {
		return err
	}
	info := RateLimitInfo{
		BucketType:   bucketType,
		Type:         item,
		LimitTimeSec: limitTimeSec,
		BucketCount:  10,
		LimitCount:   count,
	}
	info.Instance = CreateInstance(&info)
	r.limitMap[item] = info
	return nil
}

func (r *RateLimitMgr) Consume(c RateLimitConsume) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	// 消费的时候创建
	if _, ok := r.limitMap[c.LimitTypeName]; !ok {
		if err := r.CreateLimitItem(r.bucketType, c.LimitTypeName, c.Limit); err != nil {
			return err
		}
	}
	if !r.limitMap[c.LimitTypeName].Instance.AllowN(c.Count) {
		total, remain := r.limitMap[c.LimitTypeName].Instance.Remain()
		rle := NewRateLimitErr(string(c.LimitTypeName), total, remain)
		return ApiError{Code: 429, UnknownStatus: false, BizCode: 0, ErrMsg: rle.Error()}
	}
	return nil
}

func (r *RateLimitMgr) GetInstance(item RateLimitType) *RateLimitInfo {
	if info, ok := r.limitMap[item]; !ok {
		return nil
	} else {
		return &info
	}
}

func CreateInstance(info *RateLimitInfo) ratelimit.LimitI {
	if info.BucketType == ratelimit.SLIDEWINDOW {
		return ratelimit.NewSlideWindow(int64(info.LimitTimeSec), info.BucketCount, info.LimitCount)
	} else { //if info.BucketType == ratelimit.TOKENBUCKET
		return ratelimit.NewTokenBucket(float64(info.LimitCount)/float64(info.LimitTimeSec), info.LimitCount)
	}
}

type RateLimitConsume struct {
	Count, Limit  int64 // 消耗的次数。限制总量：用于创建RateLimit
	LimitTypeName RateLimitType
}

func NewRLConsume(type_ RateLimitType, count, limit int64) *RateLimitConsume {
	return &RateLimitConsume{
		LimitTypeName: type_,
		Count:         count,
		Limit:         limit,
	}
}

type ReqUrlInfo struct {
	Url                 string
	RateLimitConsumeMap []*RateLimitConsume
	OptionKeys          []string
}

func NewReqUrlInfo(url string, items ...*RateLimitConsume) ReqUrlInfo {
	if items == nil {
		items = []*RateLimitConsume{} // 防止传入nil，导致遍历错误
	}
	res := ReqUrlInfo{
		Url:                 url,
		RateLimitConsumeMap: items,
	}
	return res
}

func ParseRateLimitTime(type_ RateLimitType) (limitTimeSec float64, err error) {
	var (
		content   = strings.Split(string(type_), "_")
		limitTime time.Duration
	)
	if len(content) < 2 {
		return 0, errors.New("invalid format type:" + string(type_))
	}
	if limitTime, err = time.ParseDuration(content[1]); err != nil {
		return 0, err
	}
	limitTimeSec = float64(limitTime.Milliseconds()) / 1000
	return
}
