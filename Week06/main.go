package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type SlidingWindow struct {
	msTime   time.Time
	interval time.Duration
	number   int
	limit    int
	bucket   []int
	prevSlot int
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

	p := int(d % s.interval)
	slot := int(math.Ceil(float64(p * s.number / int(s.interval))))
	if slot == 0 && s.prevSlot == (s.number-1) {
		s.bucket[slot] = 0
	}
	s.prevSlot = slot
	if s.sumBucket() >= s.limit {
		return false
	}
	s.bucket[slot]++
	return true
}

func (s *SlidingWindow) GetBucket() []int {
	res := make([]int, s.number)
	copy(res, s.bucket)
	return res
}

func main() {
	s := NewSlidingWindow(time.Now(), 10*time.Second, 10, 20)

	for i := 1; i <= 10; i++ {
		res := s.Grant()
		fmt.Printf("Grant res%d %v\n", i, res)
		time.Sleep(300 * time.Millisecond)
	}

	fmt.Println("Grant bucket", s.GetBucket())
}
