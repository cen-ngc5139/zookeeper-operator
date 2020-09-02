package zk

import "time"

const (
	// DefaultVotingConfigExclusionsTimeout is the default timeout for setting voting exclusions.
	DefaultVotingConfigExclusionsTimeout = "30s"
	// DefaultReqTimeout is the default timeout used when performing HTTP calls against Elasticsearch
	DefaultReqTimeout = 3 * time.Minute
)
