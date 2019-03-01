package sender

import (
	"github.com/peng19940915/urlooker/web/g"
	"github.com/peng19940915/urlooker/web/utils"

	"github.com/toolkits/container/list"
)

var (
	NodeRing   *ConsistentHashNodeRing          // 服务节点的一致性哈希环
	SendQueues map[string]*list.SafeListLimited // 发送缓存队列,减少发起连接次数
	PortSendQueues map[string]*list.SafeListLimited // 发送缓存队列,减少发起连接次数
)

func initRing() {
	NodeRing = NewConsistentHashNodeRing(g.Config.Alarm.Replicas, utils.KeysOfMap(g.Config.Alarm.Cluster))
}

func initSendQueues() {
	SendQueues = make(map[string]*list.SafeListLimited)
	PortSendQueues = make(map[string]*list.SafeListLimited)
	for node, _ := range g.Config.Alarm.Cluster {
		Q1 := list.NewSafeListLimited(10240)
		Q2 := list.NewSafeListLimited(10240)
		SendQueues[node] = Q1
		PortSendQueues[node] = Q2
	}
}

func startSendTasks() {
	for node, _ := range g.Config.Alarm.Cluster {
		queue1 := SendQueues[node]
		queue2 := PortSendQueues[node]
		go SendToAlarm("url", queue1, node)
		go SendToAlarm("port", queue2, node)
	}
}

func Init() {
	initRing()
	initSendQueues()

	startSendTasks()
}
