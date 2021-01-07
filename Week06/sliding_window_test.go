package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func parallelRun(s *SlidingWindow, cnt int) int32 {
	var runCnt int32
	wg := sync.WaitGroup{}
	for j := 1; j <= cnt; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res := s.Grant()
			if res {
				atomic.AddInt32(&runCnt, 1)
			}
		}()
	}
	wg.Wait()
	return runCnt
}

func TestRolling(t *testing.T) {
	s := NewSlidingWindow(time.Now(), 5*time.Second, 5, 10)

	var runCnt int32
	for i := 0; i < 10; i++ {
		if i == 4 {
			runCnt += parallelRun(s, 3)
		}
		if i >= 7 {
			runCnt += parallelRun(s, 1)
		}
		time.Sleep(1 * time.Second)
	}

	t.Logf("runCnt %d", runCnt)
	if runCnt != 6 {
		t.Fatalf("error runCnt %d", runCnt)
	}
	expectBucket := []int{0, 0, 1, 1, 1}
	if CompareSlice(s.GetBucket(), expectBucket) == false {
		t.Fatalf("error Bucket %v\n", s.GetBucket())
	}
}

func TestOverload(t *testing.T) {
	s := NewSlidingWindow(time.Now(), 5*time.Second, 5, 5)

	var runCnt int32
	for i := 0; i < 10; i++ {
		if i == 4 || i == 5 {
			runCnt += parallelRun(s, 3)
		}
		time.Sleep(1 * time.Second)
	}

	t.Logf("runCnt %d", runCnt)
	if runCnt != 5 {
		t.Fatalf("error runCnt %d", runCnt)
	}
	expectBucket := []int{2, 0, 0, 0, 3}
	if CompareSlice(s.GetBucket(), expectBucket) == false {
		t.Fatalf("error Bucket %v\n", s.GetBucket())
	}
}

func TestSeparateSlot(t *testing.T) {
	s := NewSlidingWindow(time.Now(), 5*time.Second, 5, 5)

	var runCnt int32
	for i := 0; i < 15; i++ {
		if i == 2 {
			runCnt += parallelRun(s, 3)
		}
		if i == 12 {
			runCnt += parallelRun(s, 3)
		}
		time.Sleep(1 * time.Second)
	}

	t.Logf("runCnt %d", runCnt)
	if runCnt != 6 {
		t.Fatalf("error runCnt %d", runCnt)
	}
	expectBucket := []int{0, 0, 3, 0, 0}
	if CompareSlice(s.GetBucket(), expectBucket) == false {
		t.Fatalf("error Bucket %v\n", s.GetBucket())
	}
}

func TestRolling2(t *testing.T) {
	s := NewSlidingWindow(time.Now(), 10*time.Second, 4, 10)

	var runCnt int32
	for i := 0; i < 10; i++ {
		if i == 4 {
			runCnt += parallelRun(s, 3)
		}
		if i >= 7 {
			runCnt += parallelRun(s, 1)
		}
		time.Sleep(1 * time.Second)
	}

	t.Logf("runCnt %d", runCnt)
	if runCnt != 6 {
		t.Fatalf("error runCnt %d", runCnt)
	}
	expectBucket := []int{0, 3, 1, 2}
	if CompareSlice(s.GetBucket(), expectBucket) == false {
		t.Fatalf("error Bucket %v\n", s.GetBucket())
	}
}

func CompareSlice(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	for key, value := range a {
		if value != b[key] {
			return false
		}
	}

	return true
}
