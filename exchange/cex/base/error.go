package base

import (
	"fmt"
	"time"
)

type ApiError struct {
	Code          int32  // http状态码
	UnknownStatus bool   // 执行状态未知，调用者应该据此采取重发或者查询确保状态确定
	BizCode       int32  // 业务错误码
	ErrMsg        string // 错误信息
}

func (e ApiError) Error() string {
	return fmt.Sprintf("error: code=%d bizCode=%d message=%s", e.Code, e.BizCode, e.ErrMsg)
}

// RateLimitErr 限流错误
type RateLimitErr struct {
	Name             string
	Capacity, Remain int64
}

func (r RateLimitErr) Error() string {
	return fmt.Sprint(r.Name, " limit rate err: capacity:", r.Capacity, " remain:", r.Remain, " ts:", time.Now().UnixMicro())
}

func NewRateLimitErr(name string, capacity, remain int64) RateLimitErr {
	return RateLimitErr{name, capacity, remain}
}
