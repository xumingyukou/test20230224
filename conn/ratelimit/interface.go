package ratelimit

type LimitI interface {
	Allow() bool                      //判断当前请求
	AllowN(int64) bool                //判断当前请求
	Remain() (capacity, remain int64) //剩余
	Clear()                           //清空
	Set(used int64)                   //设置已使用
}

type BUCKET_TYPE int8

const (
	SLIDEWINDOW BUCKET_TYPE = 0
	TOKENBUCKET BUCKET_TYPE = 1
)
