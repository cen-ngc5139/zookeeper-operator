package zk

import (
	"context"
	"time"
)

// ServerStats is the information pulled from the Zookeeper `stat` command.
type ServerStats struct {
	Sent        int64
	Received    int64
	NodeCount   int64
	MinLatency  int64
	AvgLatency  int64
	MaxLatency  int64
	Connections int64
	Outstanding int64
	Epoch       int32
	Counter     int32
	BuildTime   time.Time
	Mode        Mode
	Version     string
	Error       error
}

// Mode is used to build custom server modes (leader|follower|standalone).
type Mode uint8

type ClusterStats struct {
	AvailableNodes int
	LeaderNode     string
}

func (c *BaseClient) GetClusterStatus(ctx context.Context) (ServerStats, error) {
	var result ServerStats
	return result, c.Get(ctx, "/status", &result)
}

func (c *BaseClient) GetClusterUp(ctx context.Context) (bool, error) {
	var result bool
	return result, c.Get(ctx, "/runok", &result)
}
