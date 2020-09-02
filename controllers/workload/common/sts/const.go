package sts

const (
	// prepareFilesystemContainerName is the name of the container that prepares the filesystem
	PrepareFilesystemContainerName = "zookeeper-internal-init-filesystem"
)

const (
	// EnvPodName and EnvPodIP are injected as env var into the ZK pod at runtime,
	// to be referenced in ZK configuration file
	EnvPodName = "POD_NAME"
	EnvPodIP   = "POD_IP"
)
