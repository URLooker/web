package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/toolkits/str"
	"github.com/peng19940915/urlooker/web/model"
	"github.com/toolkits/web"
	"net/http"
	"github.com/peng19940915/urlooker/web/http/param"
	"github.com/peng19940915/urlooker/web/http/errors"
	"github.com/peng19940915/urlooker/web/http/render"
	"github.com/peng19940915/urlooker/web/utils"
)

type RenderEvent struct {
	Url           string
	Id            int
	Ip            string
	RespTime      int
	Status        string
	AlarmCount    int
	AlarmTime     string
	RespCode      string
}

func GetURLAlarm(c *gin.Context) {
	MeRequired(LoginRequired(c))

	query := param.String(c.Request, "q", "")
	//onlyAlarm := param.Int(c.Request, "only_alarm", 1)
	if str.HasDangerousCharacters(query) {
		errors.Panic("查询字符不合法")
	}

	limit := param.Int(c.Request, "limit", 20)
	alarmTotal, err := model.GetAllEventCounter(query)
	errors.MaybePanic(err)
	pager := web.NewPaginator(c.Request, limit, alarmTotal)

	events, err := model.GetAllEvent(limit, pager.Offset(), query)
	errors.MaybePanic(err)

	var tmpEvents = make([] *RenderEvent, 0)
	for i, event := range events {
		tmpEvents = append(tmpEvents,&RenderEvent{
			Id: i,
			Url: event.Url,
			Ip:  event.Ip,
			RespTime: event.RespTime,
			AlarmTime: utils.TranUnix2String(event.EventTime),
			RespCode: event.RespCode,
			Status: event.Status,
		})
	}
	render.HTML(http.StatusOK, c, "alarm/url_alarm", gin.H{
		"Events": tmpEvents,
		"Pager": pager,
		"Query": query,
	})
}


type RenderPortEvent struct {
	Host          string
	Id            int
	Port          string
	Ip            string
	RespTime      int
	Status        string
	AlarmCount    int
	AlarmTime     string
}

func GETPortAlarm(c *gin.Context){
	MeRequired(LoginRequired(c))

	query := param.String(c.Request, "q", "")
	if str.HasDangerousCharacters(query) {
		errors.Panic("查询字符不合法")
	}

	limit := param.Int(c.Request, "limit", 20)
	alarmTotal, err := model.GetAllPortEventCounter(query)
	errors.MaybePanic(err)
	pager := web.NewPaginator(c.Request, limit, alarmTotal)

	events, err := model.GetAllPortEvent(limit, pager.Offset(), query)
	errors.MaybePanic(err)

	var tmpEvents = make([] *RenderPortEvent, 0)
	for i, event := range events {
		tmpEvents = append(tmpEvents,&RenderPortEvent{
			Id: i,
			Host: event.Host,
			Ip:  event.Ip,
			RespTime: event.RespTime,
			AlarmTime: utils.TranUnix2String(event.EventTime),
			Port: event.Port,
			Status: event.Status,
		})
	}

	render.HTML(http.StatusOK, c, "alarm/port_alarm", gin.H{
		"Events": tmpEvents,
		"Pager": pager,
		"Query": query,
	})
}