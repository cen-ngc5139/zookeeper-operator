package volume

const (
	DataVolClaimName = "zookeeper-data"
	AffinityOff      = "none"
)

const (
	ConfigFileName        = "zoo.cfg"
	ConfigVolumeName      = "zookeeper-internal-config"
	ConfigVolumeMountPath = "/conf"

	DynamicConfigFileVolumeName      = "zookeeper-internal-dynamic-config"
	DynamicConfigFileVolumeMountPath = "/mnt/zookeeper/dynamic-config"
	DynamicConfigFile                = "zoo_replicated1.cfg.dynamic"

	DataVolumeName = "zookeeper-data"
	DataMountPath  = "/data"

	LogsVolumeName = "zookeeper-logs"
	LogsMountPath  = "/logs"

	ScriptsVolumeName      = "zookeeper-internal-scripts"
	ScriptsVolumeMountPath = "/mnt/zookeeper/scripts"
)
