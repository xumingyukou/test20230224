package ratelimit

//https://github.com/sunanzhi/ratelimit-go/blob/main/ratelimit.go

import (
	"clients/transform"
	"container/ring"
	"math"
	"sync"
	"time"
)

// bucketNode 窗口节点
type bucketNode struct {
	startTime int64 // 节点开始时间ms
	endTime   int64 // 节点结束时间ms
	count     int64 // 节点统计访问数量
}

// SlideWindow 滑动窗口
type SlideWindow struct {
	limitTimeMs      int64      // 限流区间 单位ms
	limitCount       int64      // 限流数量
	count            int64      // 统计总使用数量
	bucketCount      int64      // 滑动窗口个数
	bucketIntervalMs int64      // 滑动窗口区间 单位ms
	bucketChain      *ring.Ring // 滑动窗口环形链

	mu sync.Mutex // 锁
}

//NewSlideWindow 初始化: 限流区间、滑动窗口个数、限流数量
func NewSlideWindow(limitTimeSec, bucketCount, limitCount int64) *SlideWindow {
	if bucketCount <= 0 {
		bucketCount = 10
	}
	s := &SlideWindow{
		limitTimeMs:      limitTimeSec,
		limitCount:       limitCount,
		count:            0,
		bucketCount:      bucketCount,
		bucketIntervalMs: int64(math.Floor(float64(limitTimeSec * 1000 / bucketCount))),
		bucketChain:      ring.New(int(bucketCount)), // 初始化滑动窗口
	}
	s.Clear()
	return s
}

func (s *SlideWindow) Clear() {
	var (
		oriStartTime = time.Now().UnixMilli()
		cur          = s.bucketChain
	)
	for i := 0; i < int(s.bucketCount); i++ {
		cur.Value = bucketNode{
			startTime: oriStartTime,
			endTime:   oriStartTime + s.bucketIntervalMs,
			count:     0,
		}
		oriStartTime -= s.bucketIntervalMs
		cur = cur.Prev()
	}
	//最后移动到时间最小的位置
	s.bucketChain = s.bucketChain.Next()
}

func (s *SlideWindow) Slide() {
	s.mu.Lock()
	defer s.mu.Unlock()
	var (
		latest           = s.bucketChain.Prev().Value.(bucketNode)
		moveCount  int64 = (time.Now().UnixMilli() - latest.endTime) / s.bucketIntervalMs
		curEndTime       = latest.endTime
	)
	if moveCount <= 0 {
		return //不滑动
	}
	for i := 0; int64(i) < transform.Min(moveCount, s.bucketCount); i++ {
		// 减去即将淘汰的节点数据统计
		s.count = transform.Max(s.count-s.bucketChain.Value.(bucketNode).count, 0)

		s.bucketChain.Value = bucketNode{
			startTime: curEndTime + (moveCount-int64(i)-1)*s.bucketIntervalMs,
			endTime:   curEndTime + (moveCount-int64(i))*s.bucketIntervalMs,
			count:     0,
		}
		s.bucketChain = s.bucketChain.Next()
	}
}

// Allow 限流
func (s *SlideWindow) Allow() bool {
	return s.AllowN(1)
}

func (s *SlideWindow) AllowN(count int64) bool {
	s.Slide()
	s.mu.Lock()
	defer s.mu.Unlock()
	// 全部节点超过限制
	if s.count+count > s.limitCount {
		return false
	}
	s.count += count
	cur := s.bucketChain.Prev().Value.(bucketNode)
	s.bucketChain.Prev().Value = bucketNode{
		startTime: cur.startTime,
		endTime:   cur.endTime,
		count:     cur.count + count,
	}
	return true
}

func (s *SlideWindow) Remain() (int64, int64) {
	return s.limitCount, s.limitCount - s.count
}

func (s *SlideWindow) Set(used int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.AllowN(transform.Max(used-s.count, 0))
}
