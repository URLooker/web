package handler

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"

	"github.com/peng19940915/urlooker/web/http/errors"
	"github.com/peng19940915/urlooker/web/http/param"
	"github.com/peng19940915/urlooker/web/http/render"
	"github.com/peng19940915/urlooker/web/model"
	"github.com/peng19940915/urlooker/web/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func AddStrategyGet(c *gin.Context) {
	render.HTML(http.StatusOK, c, "strategy/create", gin.H{})
}

func AddStrategyPost(c *gin.Context) {
	me := MeRequired(LoginRequired(c))
	var msg string
	var err error
	var tagStr string

	urls := strings.Split(param.String(c.Request, "url", ""), "\n")
	for _, url := range urls {
		err := utils.CheckUrl(url)
		if err != nil {
			errors.Panic(fmt.Sprintf("url:%s %s", url, err.Error()))
		}
	}

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

	for _, url := range urls {
		var s = model.Strategy{}
		s.Creator = me.Name
		s.Enable = param.MustInt(c.Request, "enable")
		s.Url = url
		s.ExpectCode = param.String(c.Request, "expect_code","200")

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
		s.Data = param.String(c.Request, "data", "")
		s.Tag = tagStr
		s.IP = param.String(c.Request, "ip", "")

		_, err = s.Add()
		if err != nil {
			msg += fmt.Sprintf("strategy:%s failed, err:%s", url, err.Error())
		} else {
			msg += fmt.Sprintf("strategy:%s success :)", url)
		}
	}

	//errors.MaybePanic(err)
	if err != nil {
		errMsg := fmt.Sprintf("%s,err:%v", msg, err)
		errors.Panic(errMsg)
	}
	render.Data(c, msg)
}

func GetStrategyById(c *gin.Context) {
	strategy := StraRequired(c)
	render.Data(c, strategy)
}

func UpdateStrategyGet(c *gin.Context) {
	s := StraRequired(c)
	render.HTML(http.StatusOK, c,"strategy/edit", gin.H{
		"Id": s.Id,
	})
}

func UpdateStrategy(c *gin.Context) {
	s := StraRequired(c)
	me := MeRequired(LoginRequired(c))
	var tagStr string

	username := me.Name
	if s.Creator != username && !IsAdmin(username) {
		errors.Panic("没有权限")
	}

	url := param.String(c.Request, "url", "")
	err := utils.CheckUrl(url)
	if err != nil {
		errors.Panic(fmt.Sprintf("url:%s %s", url, err.Error()))
	}

	tags := strings.Split(param.String(c.Request, "tags", ""), "\n")
	if len(tags) > 0 && tags[0] != "" {
		log.Info("len: %d", len(tags))
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

	s.Creator = username
	s.Url = url
	s.Enable = param.MustInt(c.Request, "enable")
	s.ExpectCode = param.String(c.Request, "expect_code", "200")
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
	s.Data = param.String(c.Request, "data", "")
	s.IP = param.String(c.Request, "ip", "")
	s.Tag = tagStr

	err = s.Update()
	errors.MaybePanic(err)
	render.Data(c, "ok")
}

func DeleteStrategy(c *gin.Context) {
	me := MeRequired(LoginRequired(c))
	strategy := StraRequired(c)
	//teams := strings.Split(strategy.Teams, ",")

	username := me.Name
	if strategy.Creator != username && !IsAdmin(username) {
		errors.Panic("没有权限")
	}

	err := strategy.Delete()
	errors.MaybePanic(err)
	err = model.DeleteOldEvent(strategy.Id)
	errors.MaybePanic(err)
	render.Data(c, "ok")
}

func GetTeamsOfStrategy(c *gin.Context) {
	MeRequired(LoginRequired(c))
	stra := StraRequired(c)
	teams, err := model.GetTeamsByIds(stra.Teams)
	errors.MaybePanic(err)
	render.Data(c, map[string]interface{}{"teams": teams})
}
