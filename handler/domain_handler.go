package handler

import (
	"net/http"

	"github.com/urlooker/web/g"
	"github.com/urlooker/web/http/errors"
	"github.com/urlooker/web/http/param"
	"github.com/urlooker/web/http/render"
	"github.com/urlooker/web/model"
	"github.com/urlooker/web/utils"
)

type Url struct {
	Ip     string              `json:"ip"`
	Status []*model.ItemStatus `json:"status"`
}

func UrlStatus(w http.ResponseWriter, r *http.Request) {
	sid := param.MustInt64(r, "id")

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

	render.Put(r, "AlarmOn", g.Config.Alarm.Enable)
	render.Put(r, "TimeData", timeData)
	render.Put(r, "Id", sid)
	render.Put(r, "Url", strategy.Url)
	render.Put(r, "Events", events)
	render.Put(r, "Data", urlArr)
	render.HTML(r, w, "chart/index")
}
