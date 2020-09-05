package provision

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/model"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/observer"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/utils"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/zk"
	zkcli "github.com/samuel/go-zookeeper/zk"
	"gopkg.in/fatih/set.v0"
)

type AddMember struct {
	Record string `json:"record"`
}

func (p Provision) ReConfig() error {
	name := p.Workload.GetName()
	namespace := p.Workload.GetNamespace()

	//获取目前正在允许的pod信息
	currentPods, err := utils.GetCurrentPods(p.Client, p.Workload, p.Labels, p.Log)
	if err != nil {
		p.Log.Info(
			"Unable to get current pods.",
			"error", err,
			"namespace", namespace,
			"zk_name", name,
		)
		return errors.New("reconfig need requeue")
	}

	//获取期望 zk config内容
	podNames := utils.PodNames(*p.ExpectSts)
	AddMemberRecords := utils.GetPodIp(p.Workload, podNames)
	expectConfigRecord := set.New(set.ThreadSafe)

	for _, v := range AddMemberRecords {
		expectConfigRecord.Add(v)
	}

	if len(currentPods) == 0 {
		p.Log.Info("Cluster is building ,pls hold on.")
		return nil
	}
	//生成zk connect client
	randomPod := currentPods[rand.Intn(len(currentPods))]
	podIp := randomPod.Status.PodIP

	zkAgentUrl := fmt.Sprintf("http://%s:%d", podIp, model.AgentPort)

	cli := &zk.BaseClient{
		HTTP:      &http.Client{},
		Endpoint:  zkAgentUrl,
		Transport: &http.Transport{},
	}

	//获取实际 zk config内容
	ctx, cancel := context.WithCancel(context.Background())
	timeoutCtx, cancel := context.WithTimeout(ctx, observer.DefaultSettings.RequestTimeout)
	defer cancel()
	var config AddMember
	if err := cli.Get(timeoutCtx, "/get", &config); err != nil {
		p.Log.Info(
			"Unable to get zk config.",
			"error", err,
			"url", zkAgentUrl,
			"pod", randomPod.Name,
		)
		return err
	}

	//返回config字符串转数组
	currentConfigs := strings.Split(config.Record, "\n")

	//去掉config 末尾 version=100000104 项
	currentConfigs = currentConfigs[:len(currentConfigs)-1]

	//config数组转set interface{}
	currentConfigRecord := set.New(set.ThreadSafe)
	for _, v := range currentConfigs {
		currentConfigRecord.Add(v)
	}

	//比较实际和期待 zk config
	//needAdd为需要添加配置项
	//needDel为需要删除配置项
	needAdd := set.Difference(expectConfigRecord, currentConfigRecord)
	needDel := set.Difference(currentConfigRecord, expectConfigRecord)

	//set interface转[]string
	needAddArray := set.StringSlice(needAdd)
	needDelArray := set.StringSlice(needDel)

	//如果待添加和待删除数组均为0，退出reconfig流程
	if len(needAddArray) == 0 && len(needDelArray) == 0 {
		p.Log.Info(
			"Don't need to update zk config.",
			"url", zkAgentUrl,
			"pod", randomPod.Name,
		)
		return nil
	}

	//调用zk-agent接口添加新节点到集群中reconfig
	for _, record := range needAddArray {
		member := &AddMember{
			Record: record,
		}
		ctx, cancel := context.WithCancel(context.Background())
		timeoutCtx, cancel := context.WithTimeout(ctx, observer.DefaultSettings.RequestTimeout)
		defer cancel()

		var result zkcli.Stat

		if err := cli.Post(timeoutCtx, "/add", member, &result); err != nil {
			p.Log.Info(
				"Unable to add member to zk.",
				"error", err,
				"url", zkAgentUrl,
				"pod", randomPod.Name,
			)
			continue
		}
		//fmt.Println(result)
	}

	//调用zk-agent接口删除新节点到集群中reconfig
	for _, record := range needDelArray {
		member := &AddMember{
			Record: record,
		}
		ctx, cancel := context.WithCancel(context.Background())
		timeoutCtx, cancel := context.WithTimeout(ctx, observer.DefaultSettings.RequestTimeout)
		defer cancel()

		var result zkcli.Stat

		if err := cli.Post(timeoutCtx, "/del", member, &result); err != nil {
			p.Log.Info(
				"Unable to add member to zk.",
				"error", err,
				"url", zkAgentUrl,
				"pod", randomPod.Name,
			)
			continue
		}
		//fmt.Println(result)
	}

	return nil
}
