package api

import (
	"log"
	"time"

	"github.com/urlooker/web/g"
	"github.com/urlooker/web/model"
	"github.com/urlooker/web/sender"
	"github.com/urlooker/web/utils"
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
			log.Println("save sid_ip error:", err)
			*reply = "save sid_ip error:" + err.Error()
			return nil
		}

		err = itemStatus.Save()
		if err != nil {
			log.Println("save item error:", err)
			*reply = "save item error:" + err.Error()
			return nil
		}

		if g.Config.Alarm.Enable {
			node, err := sender.NodeRing.GetNode(itemStatus.PK())
			if err != nil {
				log.Println("error:", err)
				*reply = "get node error:" + err.Error()
				return nil
			}

			Q := sender.SendQueues[node]
			isSuccess := Q.PushFront(itemStatus)
			if !isSuccess {
				log.Println("error:", err)
				*reply = "save item error:" + err.Error()
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

type GetItemResponse struct {
	Message string
	Data    []*g.DetectedItem
}

func (this *Web) GetItem(hostname string, resp *GetItemResponse) error {
	items, exists := g.DetectedItemMap.Get(hostname)
	if !exists {
		resp.Message = "no found"
	}
	resp.Data = items
	return nil
}
