package cm

const (
	ConfigFileName        = "zoo.cfg"
	ConfigVolumeName      = "zookeeper-internal-config"
	ConfigVolumeMountPath = "/mnt/zookeeper/zookeeper-config"

	DynamicConfigFileVolumeName      = "zookeeper-internal-dynamic-config"
	DynamicConfigFileVolumeMountPath = "/mnt/zookeeper/dynamic-config"
	DynamicConfigFile                = "zoo_replicated1.cfg.dynamic"

	DataVolumeName = "zookeeper-data"
	DataMountPath  = "/data"

	LogsVolumeName = "zookeeper-logs"
	LogsMountPath  = "/logs"

	ScriptsVolumeName      = "zookeeper-internal-scripts"
	ScriptsVolumeMountPath = "/mnt/zookeeper/scripts"

	AgentVolumeName      = "zookeeper-agent-config"
	AgentVolumeMountPath = "/mnt/zookeeper/agent"
)
