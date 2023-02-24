package ratelimit

import (
	"clients/logger"
	"clients/transform"
	"sync"
	"time"
)

type TokenBucket struct {
	rate        float64 //固定的token放入速率, r/s
	capacity    int64   //桶的容量
	tokens      int64   //桶中当前token数量
	lastTokenMs int64   //桶上次放token的时间戳 ms

	lock sync.Mutex
}

func NewTokenBucket(rate float64, capacity int64) *TokenBucket {
	l := &TokenBucket{
		rate:        rate,
		tokens:      capacity, // 热启动
		capacity:    capacity,
		lastTokenMs: time.Now().UnixMilli(),
	}
	return l
}

func (l *TokenBucket) Allow() bool {
	return l.AllowN(1)
}

func (l *TokenBucket) AllowN(count int64) bool {
	l.lock.Lock()
	defer l.lock.Unlock()
	now := time.Now().UnixMilli()
	addTokens := int64(float64(now-l.lastTokenMs) * l.rate / 1000)

	l.tokens = transform.Min(l.tokens+addTokens, l.capacity) // 先添加令牌
	if addTokens > 0 {
		l.lastTokenMs = now
	}
	if l.tokens >= count {
		// 还有令牌，领取令牌
		l.tokens -= count
		logger.Logger.Trace("remaining token ", l.tokens, "/", l.capacity)
		return true
	} else {
		// 没有令牌,则拒绝
		return false
	}
}

func (l *TokenBucket) Set(used int64) {
	l.tokens = transform.Max(l.capacity-used, 0)
	l.lastTokenMs = time.Now().UnixMilli()
}

func (l *TokenBucket) Clear() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.tokens = 0
	l.lastTokenMs = time.Now().UnixMilli()
}

func (l *TokenBucket) Remain() (int64, int64) {
	return l.capacity, l.tokens
}
