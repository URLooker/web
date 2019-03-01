package api

import (
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/peng19940915/urlooker/web/g"
	"github.com/peng19940915/urlooker/web/model"
	"github.com/peng19940915/urlooker/web/sender"
	"github.com/peng19940915/urlooker/web/utils"
)

type SendResultReq struct {
	Hostname     string
	CheckResults []*g.CheckResult
}

func (this *Web) SendResult(req SendResultReq, reply *string) error {
	for _, arg := range req.CheckResults {
		itemStatus := model.ItemStatus{
			Ip:       arg.Ip,
			Sid:      arg.Sid,
			RespTime: arg.RespTime,
			RespCode: arg.RespCode,
			PushTime: arg.PushTime,
			Result:   arg.Status,
		}

		relSidIp := model.RelSidIp{
			Sid: arg.Sid,
			Ip:  arg.Ip,
			Ts:  time.Now().Unix(),
		}

		err := relSidIp.Save()
		if err != nil {
			log.Errorf("save sid_ip error, detail: %v", err)
			*reply = "save sid_ip error:" + err.Error()
			return nil
		}

		err = itemStatus.Save()
		if err != nil {
			log.Errorf("save item error, detail: %v", err.Error())
			*reply = "save item error:" + err.Error()
			return nil
		}

		if g.Config.Alarm.Enable {
			node, err := sender.NodeRing.GetNode(itemStatus.PK())
			if err != nil {
				log.Errorf("Get Node from node ring error, detail:%v:", err.Error())
				*reply = "get node error:" + err.Error()
				return nil
			}

			Q := sender.SendQueues[node]
			isSuccess := Q.PushFront(itemStatus)
			if !isSuccess {
				log.Println("push itemStatus error:", itemStatus)
				*reply = "push itemStatus error"
				return nil
			}
		}

	}

	if g.Config.Falcon.Enable {
		if len(req.CheckResults) > 0 {
			utils.PushFalcon(req.CheckResults, req.Hostname)
		}
	}

	*reply = ""
	return nil
}

type SendPortResultReq struct {
	Hostname     string
	CheckResults []*g.CheckPortResult
}

func (this *Web) SendPortResult(req SendPortResultReq, reply *string) error {
	for _, arg := range req.CheckResults {

		portStatus := model.PortStatus{
			Sid:      arg.Sid,
			RespTime: arg.RespTime,
			Result:   arg.Result,
			PushTime: arg.PushTime,
		}


		err := portStatus.Save()
		if err != nil {
			log.Errorf("save item error, detail: %v", err)
			*reply = "save item error:" + err.Error()
			return nil
		}

		if g.Config.Alarm.Enable {
			node, err := sender.NodeRing.GetNode(portStatus.PK())
			if err != nil {
				log.Errorf("Get Node from node ring error, detail:%v:", err)
				*reply = "get node error:" + err.Error()
				return nil
			}

			Q := sender.PortSendQueues[node]
			isSuccess := Q.PushFront(portStatus)
			if !isSuccess {
				log.Error("push item Port Status error:", portStatus)
				*reply = "push itemStatus error"
				return nil
			}
		}

	}

	if g.Config.Falcon.Enable {
		if len(req.CheckResults) > 0 {
			utils.PushPort2Falcon(req.CheckResults, req.Hostname)
		}
	}

	*reply = ""
	return nil
}



type GetItemResponse struct {
	Message string
	Data    []*g.DetectedItem
}

func (this *Web) GetItem(hostname string, resp *GetItemResponse) error {
	items, exists := g.DetectedItemMap.Get(hostname)
	if !exists {
		resp.Message = "no found item assigned to " + hostname
	}
	resp.Data = items
	return nil
}

type GetPortItemResponse struct {
	Message string
	Data [] *g.DetectedPortItem
}

func (this *Web) GetPortItem(hostname string, resp *GetPortItemResponse) error {

	items, exists := g.DetectedPortItemMap.Get(hostname)
	if !exists {
		resp.Message = "no found item assigned to " + hostname
	}
	resp.Data = items
	return nil
}
