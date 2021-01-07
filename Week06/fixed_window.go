package main

// 固定窗口算法（计数器）
type FixedWindow struct {
	slot     int
	interval int
	count    int
	limit    int
}
