package main

import (
	"math"
	"sync"
	"time"
)

// 滑动窗口算法
type SlidingWindow struct {
	msTime   time.Time
	interval time.Duration
	number   int
	limit    int
	bucket   []int
	ciSlot   int
	cnSlot   int
	lock     sync.Locker
}

func NewSlidingWindow(msTime time.Time, interval time.Duration, number, limit int) *SlidingWindow {
	sw := new(SlidingWindow)
	sw.msTime = msTime
	sw.interval = interval
	sw.number = number
	sw.limit = limit
	sw.bucket = make([]int, number)
	sw.lock = new(sync.Mutex)
	return sw
}

func (s *SlidingWindow) sumBucket() int {
	sum := 0
	for _, v := range s.bucket {
		sum += v
	}
	return sum
}

func (s *SlidingWindow) Grant() bool {
	d := time.Now().Sub(s.msTime)
	s.lock.Lock()
	defer s.lock.Unlock()

	iSlot := int(math.Floor(float64(d / s.interval)))
	if iSlot >= s.ciSlot+2 {
		for k, _ := range s.bucket {
			s.bucket[k] = 0
		}
		s.ciSlot = iSlot
	}
	p := int(d % s.interval)
	nSlot := int(math.Floor(float64(p*s.number) / float64(s.interval)))
	if nSlot != s.cnSlot {
		s.bucket[nSlot] = 0
		s.cnSlot = nSlot
	}
	if s.sumBucket() >= s.limit {
		return false
	}
	s.bucket[nSlot]++
	return true
}

func (s *SlidingWindow) GetBucket() []int {
	res := make([]int, s.number)
	copy(res, s.bucket)
	return res
}
