package cache

import "time"

// var (
// 	defaultSyncInterval =
// 	defaultMonitorInterval =
// )

const (
	defaultSyncInterval    = time.Minute * 1
	defaultMonitorInterval = time.Second * 1
)

type Options struct {
	syncInterval    time.Duration
	monitorInterval time.Duration
}
