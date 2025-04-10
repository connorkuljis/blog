package model

import (
	"runtime"
	"time"
)

type NerdStats struct {
	PageCount int

	StartTime  time.Time
	FinishTime time.Time
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

func (n *NerdStats) CalculateDuration() time.Duration {
	return n.FinishTime.Sub(n.StartTime)
}
