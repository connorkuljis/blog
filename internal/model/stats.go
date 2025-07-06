package model

import (
	"runtime"
	"time"
)

type NerdStats struct {
	PageCount int
	Platform  string
	Arch      string
	Version   string

	StartTime time.Time
}

func NewNerdStats(startTime time.Time) *NerdStats {
	return &NerdStats{
		StartTime: startTime,
		Platform:  runtime.GOOS,
		Arch:      runtime.GOARCH,
		Version:   runtime.Version(),
	}
}
