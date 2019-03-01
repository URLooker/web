package handler

import (
	"net/http"

	"github.com/peng19940915/urlooker/web/g"
	"github.com/peng19940915/urlooker/web/http/errors"
	"github.com/peng19940915/urlooker/web/http/param"
	"github.com/peng19940915/urlooker/web/model"
	"github.com/peng19940915/urlooker/web/utils"
	"github.com/gin-gonic/gin"
	"github.com/peng19940915/urlooker/web/http/render"
)

type Url struct {
	Ip     string              `json:"ip"`
	Status []*model.ItemStatus `json:"status"`
}

func UrlStatus(c *gin.Context) {
	sid := param.MustInt64(c.Request, "id")

	sidIpIndex, err := model.RelSidIpRepo.GetBySid(sid)
	errors.MaybePanic(err)

	urlArr := make([]Url, 0)
	idx := 0
	var ts int64 = 0
	for i, index := range sidIpIndex {
		url := Url{
			Ip: index.Ip,
		}
		url.Status, err = model.ItemStatusRepo.GetByIpAndSid(index.Ip, index.Sid)
		errors.MaybePanic(err)

		if len(url.Status) > 0 && ts < url.Status[len(url.Status)-1].PushTime {
			ts = url.Status[len(url.Status)-1].PushTime
			idx = i
		}

		urlArr = append(urlArr, url)
	}

	//绘图使用，时间轴
	var timeData []string
	if len(urlArr) > 0 {
		for _, item := range urlArr[idx].Status {
			t := utils.TimeFormat(item.PushTime)
			timeData = append(timeData, t)
		}
	}

	events, err := model.EventRepo.GetByStrategyId(sid, g.Config.Past*60)

	errors.MaybePanic(err)

	strategy, err := model.GetStrategyById(sid)
	errors.MaybePanic(err)

	render.HTML(http.StatusOK, c,"chart/index", gin.H{
		"AlarmOn": g.Config.Alarm.Enable,
		"TimeData": timeData,
		"Id": sid,
		"Url": strategy.Url,
		"Events": events,
		"Data": urlArr,
	})
}

func PortStatus(c *gin.Context) {
	sid := param.MustInt64(c.Request, "id")
	var timeData [] string
	portArr, err := model.PortStatusRepo.GetBySid(sid)
	errors.MaybePanic(err)
	// 绘图使用，时间轴
	if len(portArr) > 0 {
		for _, item := range(portArr){
			t := utils.TimeFormat(item.PushTime)
			timeData = append(timeData, t)
		}
	}
	events, err := model.PortEventRepo.GetByPortStrategyId(sid, g.Config.Past * 60)
	errors.MaybePanic(err)

	strategy, err := model.GetPortStrategyById(sid)
	errors.MaybePanic(err)
	render.HTML(http.StatusOK, c,"chart/port_index", gin.H{
		"AlarmOn": g.Config.Alarm.Enable,
		"TimeData": timeData,
		"Id": sid,
		"Host": strategy.Host,
		"Port": strategy.Port,
		"Events": events,
		"Data": portArr,
	})
}
