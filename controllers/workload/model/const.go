package model

const (
	AppLabel                      = "zookeeper"
	RoleName                      = "zookeeper"
	ContainerName                 = "zookeeper"
	DefaultImageRepository string = "ghostbaby/zookeeper"
)

const (
	ZooKeeperStatefulset         = "ZooKeeperStatefulset"
	ZooKeeperService             = "ZooKeeperService"
	ZooKeeperConfigMap           = "ZooKeeperConfigMap"
	ZooKeeperServiceMonitor      = "ZooKeeperServiceMonitor"
	ZooKeeperPodDisruptionBudget = "ZooKeeperPodDisruptionBudget"
	ZooKeeperUserSecret          = "ZooKeeperUserSecret"
	ZooKeeperKeySecret           = "ZooKeeperKeySecret"
	ZooKeeperPrometheusRules     = "ZooKeeperPrometheusRules"

	ZooKeeperDownScaling = "ZooKeeperDownScaling"
	ZooKeeperUpScaling   = "ZooKeeperUpScaling"

	MessageZooKeeperStatefulset         = "ZooKeeper  Statefulset %s already created."
	MessageZooKeeperService             = "ZooKeeper  Service %s already created."
	MessageZooKeeperConfigMap           = "ZooKeeper  ConfigMap %s already created."
	MessageZooKeeperServiceMonitor      = "ZooKeeper  ServiceMonitor %s already created."
	MessageZooKeeperPrometheusRules     = "ZooKeeper  PrometheusRules %s already created."
	MessageZooKeeperPodDisruptionBudget = "ZooKeeper  PodDisruptionBudget %s already created."
	MessageZooKeeperUserSecret          = "ZooKeeper  User Secret %s already created."
	MessageZooKeeperKeySecret           = "ZooKeeper  Key Secret %s already created."

	UpdateMessageZooKeeperStatefulset = "ZooKeeper  Statefulset %s already update."
	UpdateMessageZooKeeperConfigMap   = "ZooKeeper  ConfigMap %s already update."

	MessageZooKeeperDownScaling = "ZooKeeper downscale from %d to %d"
	MessageZooKeeperUpScaling   = "ZooKeeper upscale from %d to %d"
)

const (
	ExporterPort                 = 9114
	ExporterPortName             = "http-metrics"
	ExporterContainerName        = "zk-exporter"
	ExporterDefaultRequestCPU    = "25m"
	ExporterDefaultLimitCPU      = "50m"
	ExporterDefaultRequestMemory = "50Mi"
	ExporterDefaultLimitMemory   = "100Mi"
)

const (
	AgentPort                 = 1988
	AgentPortName             = "http-agent"
	AgentContainerName        = "zk-agent"
	AgentDefaultRequestCPU    = "25m"
	AgentDefaultLimitCPU      = "50m"
	AgentDefaultRequestMemory = "100Mi"
	AgentDefaultLimitMemory   = "200Mi"
)

const (
	ClientPort         = 2181
	ServerPort         = 2888
	LeaderElectionPort = 3888
	AgentHTTPPort      = 1988

	ClusterIPServiceType = "ClusterIP"
	HeadlessServiceType  = "Headless"
)

const (
	ServiceMonitorInterval = "30s"
	ServiceMonitorPort     = "http-metrics"
	ServiceMonitorCrdName  = "servicemonitors.monitoring.coreos.com"
	PrometheusRulesCrdName = "prometheusrules.monitoring.coreos.com"
	StorageLowAlertName    = "存储空间低"
)
