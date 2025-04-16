package model

import (
	"runtime"
	"time"
)

type NerdStats struct {
	PageCount int
	StartTime time.Time
}

func NewNerdStats() *NerdStats {
	return &NerdStats{
		StartTime: time.Now(),
	}
}

func (n *NerdStats) SetPageCount(count int) {
	n.PageCount = count
}

func (n *NerdStats) Platform() string {
	return runtime.GOOS
}

func (n *NerdStats) Arch() string {
	return runtime.GOARCH
}

func (n *NerdStats) Version() string {
	return runtime.Version()
}
