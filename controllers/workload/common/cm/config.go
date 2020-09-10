package cm

import (
	"fmt"
)

// PodName returns the name of the pod with the given ordinal for this StatefulSet.
func PodName(ssetName string, ordinal int32) string {
	return fmt.Sprintf("%s-%d", ssetName, ordinal)
}

func GenZkConfig() string {
	return `tickTime=2000
initLimit=10
skipACL=yes
syncLimit=5
dataDir=/data
maxClientCnxns=300
dataLogDir=/logs
reconfigEnabled=true
standaloneEnabled=false
autopurge.snapRetainCount=20
autopurge.purgeInterval=24
4lw.commands.whitelist=cons, envi, conf, crst, srvr, stat, mntr, ruok
dynamicConfigFile=/conf/zoo.cfg.dynamic
`
}
