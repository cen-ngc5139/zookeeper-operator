package cm

import (
	"fmt"
	"strings"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/utils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ZkConfigMap struct {
	Data map[string]string `json:"data"`
	Name string            `json:"name"`
}

func (c *CM) GenerateConfigMap() ([]*corev1.ConfigMap, error) {
	name := c.Workload.Name
	namespace := c.Workload.Namespace
	var list []*ZkConfigMap

	//生成agent启动配置文件
	agent := GenZkAgentConfig()

	agentData := map[string]string{
		AgentConfigKey: agent,
	}

	list = append(list, &ZkConfigMap{
		Data: agentData,
		Name: genConfigMapName(name, AgentVolumeName),
	})

	//生成初始化脚本
	fsScript, err := RenderPrepareFsScript()
	if err != nil {
		return nil, err
	}

	scriptData := map[string]string{
		PrepareFsScriptConfigKey: fsScript,
	}

	list = append(list, &ZkConfigMap{
		Data: scriptData,
		Name: genConfigMapName(name, ScriptsVolumeName),
	})

	//生成zookeeper启动配置文件
	config := GenZkConfig()

	configData := map[string]string{
		ConfigFileName: config,
	}

	list = append(list, &ZkConfigMap{
		Data: configData,
		Name: genConfigMapName(name, ConfigVolumeName),
	})

	//生成动态配置文件
	podNames := utils.PodNames(*c.ExpectSts)
	podIpArray := utils.GetPodIp(c.Workload, podNames)

	var hosts string
	if podIpArray != nil {
		hosts = strings.Join(podIpArray, "\n")
	}

	dynamicConfigData := map[string]string{
		DynamicConfigFile: hosts,
	}

	list = append(list, &ZkConfigMap{
		Data: dynamicConfigData,
		Name: genConfigMapName(name, DynamicConfigFileVolumeName),
	})

	cmList := GenConfigMap(list, namespace, c.Labels)

	return cmList, nil
}

func GenConfigMap(listArray []*ZkConfigMap, namespace string, labels map[string]string) []*corev1.ConfigMap {
	var cmList []*corev1.ConfigMap
	for _, data := range listArray {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      data.Name,
				Namespace: namespace,
				Labels:    labels,
			},
			Data: data.Data,
		}
		cmList = append(cmList, cm)
	}
	return cmList
}

func genConfigMapName(name, cmType string) string {
	return fmt.Sprintf("%s-%s", name, cmType)
}
