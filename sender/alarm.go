package sender

import (
	log "github.com/sirupsen/logrus"
	"time"
	"github.com/toolkits/container/list"

	"github.com/peng19940915/urlooker/web/backend"
	"github.com/peng19940915/urlooker/web/g"
)

func SendToAlarm(alarmType string, Q *list.SafeListLimited, node string) {
	cfg := g.Config
	batch := cfg.Alarm.Batch
	addr := cfg.Alarm.Cluster[node]

	//todo：rpc 当数据量增大时，rpc调用改为并行方式
	for {
		items := Q.PopBackBy(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(time.Duration(cfg.Alarm.SleepTime) * time.Second)
			continue
		}

		var resp string
		var err error
		sendOk := false
		for i := 0; i < 3; i++ {
			rpcClient := backend.NewRpcClient(addr)

			if alarmType == "port"{
				err = rpcClient.Call("Alarm.SendPort", items, &resp)
			}else if alarmType == "url" {
				err = rpcClient.Call("Alarm.Send", items, &resp)
			}

			if err == nil {
				sendOk = true
				break
			}
			time.Sleep(1)
		}

		if !sendOk {
			log.Printf("send alarm %s:%s fail: %v", node, addr, err)
		}
		log.Debug("<=", resp)
	}
}
