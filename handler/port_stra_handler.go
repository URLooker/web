package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"github.com/peng19940915/urlooker/web/http/errors"
	"github.com/peng19940915/urlooker/web/http/param"
	"github.com/peng19940915/urlooker/web/http/render"
	"github.com/peng19940915/urlooker/web/model"
	"strconv"
	"github.com/gin-gonic/gin"
)

func AddPortStrategyGet(c *gin.Context) {
	render.HTML(http.StatusOK, c, "strategy/port_create",gin.H{})
}

func AddPortStrategyPost(c *gin.Context) {
	me := MeRequired(LoginRequired(c))
	var msg string
	var err error
	var tagStr string

	tags := strings.Split(param.String(c.Request, "tags", ""), "\n")
	if len(tags) > 0 && tags[0] != "" {
		for _, tag := range tags {
			strs := strings.Split(tag, "=")
			if len(strs) != 2 {
				errors.Panic("tag must be like aaa=bbb")
			}
		}
		tagStr = strings.Join(tags, ",")
	}

	var s = model.PortStrategy{}

	s.Creator = me.Name
	s.Enable = param.MustInt(c.Request, "enable")
	s.Host = param.MustString(c.Request, "host")
	s.Port = strconv.Itoa(param.MustInt(c.Request, "port"))
	timeout, err := strconv.Atoi(param.String(c.Request, "timeout", "6000"))
	if err != nil {
		errors.Panic("timeout 必须是数字")
	}
	s.Timeout = timeout
	maxStep, err := strconv.Atoi(param.String(c.Request, "max_step", "3"))
	if err != nil {
		errors.Panic("告警发送次数必须是数字")
	}
	s.MaxStep = maxStep

	times, err := strconv.Atoi(param.String(c.Request, "times", "3"))
	if err != nil {
		errors.Panic("连续异常次数必须是数字")
	}
	s.Times = times
	s.Teams = param.String(c.Request, "teams", "")
	if s.Teams == "" {
		errors.Panic("请填写正确的告警组")
	}

	s.Note = param.String(c.Request, "note", "")
	s.Keywords = param.String(c.Request, "keywords", "")
	s.Tag = tagStr
	s.IP = param.MustString(c.Request, "ip", "")

	_, err = s.Add()
	if err != nil {
		msg += fmt.Sprintf(" port strategy:%s failed, err:%s", s.Host, err.Error())
	} else {
		msg += fmt.Sprintf("port strategy:%s success :)", s.Host)
	}


	//errors.MaybePanic(err)
	if err != nil {
		errMsg := fmt.Sprintf("%s,err:%v", msg, err)
		errors.Panic(errMsg)
	}
	render.Data(c, msg)
}

func GetPortStrategyById(c *gin.Context) {
	strategy := PortStraRequired(c)
	render.Data(c, strategy)
}

func UpdatePortStrategyGet(c *gin.Context) {
	s := PortStraRequired(c)
	render.HTML(http.StatusOK, c, "strategy/port_edit", gin.H{
		"Id": s.Id,
	})
}

func UpdatePortStrategy(c *gin.Context) {
	s := PortStraRequired(c)
	me := MeRequired(LoginRequired(c))
	var tagStr string

	username := me.Name
	if s.Creator != username && !IsAdmin(username) {
		errors.Panic("没有权限")
	}

	host := param.String(c.Request, "host", "")


	tags := strings.Split(param.String(c.Request, "tags", ""), "\n")
	if len(tags) > 0 && tags[0] != "" {
		log.Println("len:", len(tags))
		for _, tag := range tags {
			strs := strings.Split(tag, "=")
			if len(strs) != 2 {
				errors.Panic("tag must be like aaa=bbb")
			} else if strs[0] == "" || strs[1] == "" {
				errors.Panic("tag must be like aaa=bbb")
			}
		}
		tagStr = strings.Join(tags, ",")
	}
	port := param.MustString(c.Request, "port", "")
	if host == "" || port == "" {
		errors.Panic("host or port is null!")
	}
	s.Creator = username
	s.Host = host
	s.Enable = param.MustInt(c.Request, "enable")
	s.Port = strconv.Itoa(param.MustInt(c.Request, "port"))
	timeout, err := strconv.Atoi(param.String(c.Request, "timeout", "6000"))
	if err != nil {
		errors.Panic("timeout 必须是数字")
	}
	s.Timeout = timeout
	maxStep, err := strconv.Atoi(param.String(c.Request, "max_step", "3"))
	if err != nil {
		errors.Panic("告警发送次数必须是数字")
	}
	s.MaxStep = maxStep

	times, err := strconv.Atoi(param.String(c.Request, "times", "3"))
	if err != nil {
		errors.Panic("连续异常次数必须是数字")
	}
	s.Times = times
	s.Teams = param.String(c.Request, "teams", "")
	s.Note = param.String(c.Request, "note", "")
	s.Keywords = param.String(c.Request, "keywords", "")
	s.IP = param.MustString(c.Request, "ip", "")
	s.Tag = tagStr

	err = s.Update()
	errors.MaybePanic(err)
	render.Data(c, "ok")
}

func DeletePortStrategy(c *gin.Context) {
	me := MeRequired(LoginRequired(c))
	strategy := PortStraRequired(c)
	//teams := strings.Split(strategy.Teams, ",")

	username := me.Name
	if strategy.Creator != username && !IsAdmin(username) {
		errors.Panic("没有权限")
	}
	err := strategy.Delete()
	errors.MaybePanic(err)
	err = model.DeleteOldPortEvent(strategy.Id)
	errors.MaybePanic(err)
	render.Data(c, "ok")
}

func GetPortTeamsOfStrategy(c *gin.Context) {
	MeRequired(LoginRequired(c))
	stra := PortStraRequired(c)
	teams, err := model.GetTeamsByIds(stra.Teams)
	errors.MaybePanic(err)
	render.Data(c, map[string]interface{}{"teams": teams})
}
