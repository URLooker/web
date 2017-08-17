package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/urlooker/web/http/errors"
	"github.com/urlooker/web/http/param"
	"github.com/urlooker/web/http/render"
	"github.com/urlooker/web/model"
	"github.com/urlooker/web/utils"
)

func AddStrategyGet(w http.ResponseWriter, r *http.Request) {
	idcs := utils.GetAllIdc()
	render.Data(r, "Idcs", idcs)
	render.HTML(r, w, "strategy/create")
}

func AddStrategyPost(w http.ResponseWriter, r *http.Request) {
	me := MeRequired(LoginRequired(w, r))
	var msg string
	var err error
	var tagStr string

	urls := strings.Split(param.String(r, "url", ""), "\n")
	for _, url := range urls {
		err := utils.CheckUrl(url)
		if err != nil {
			errors.Panic(fmt.Sprintf("url:%s %s", url, err.Error()))
		}
	}

	tags := strings.Split(param.String(r, "tags", ""), "\n")
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
		s.Url = url
		s.ExpectCode = param.String(r, "expect_code", "200")
		s.Timeout = param.Int(r, "timeout", 3000)
		s.MaxStep = param.Int(r, "max_step", 3)
		s.Teams = param.String(r, "teams", "")
		s.Times = param.Int(r, "times", 3)
		s.Note = param.String(r, "note", "")
		s.Keywords = param.String(r, "keywords", "")
		s.Data = param.String(r, "data", "")
		monitorIdc := param.String(r, "monitor_idc", "")
		if monitorIdc == "" {
			monitorIdc = utils.GetAllIdc()
		}

		s.MonitorIdc = monitorIdc
		s.Tag = tagStr

		_, err = s.Add()
		if err != nil {
			msg += fmt.Sprintf("strategy:%s failed, err:%s", url, err.Error())
		} else {
			msg += fmt.Sprintf("strategy:%s success :)")
		}
	}

	render.AutoJSON(w, err, msg)
}

func GetStrategyById(w http.ResponseWriter, r *http.Request) {
	strategy := StraRequired(r)
	render.AutoJSON(w, nil, strategy)
}

func UpdateStrategyGet(w http.ResponseWriter, r *http.Request) {
	s := StraRequired(r)
	render.Data(r, "Id", s.Id)
	render.HTML(r, w, "strategy/edit")
}

func UpdateStrategy(w http.ResponseWriter, r *http.Request) {
	s := StraRequired(r)
	me := MeRequired(LoginRequired(w, r))
	var tagStr string

	username := me.Name
	if s.Creator != username && s.Creator != "qinyening" {
		errors.Panic("没有权限")
	}

	url := param.String(r, "url", "")
	err := utils.CheckUrl(url)
	if err != nil {
		errors.Panic(fmt.Sprintf("url:%s %s", url, err.Error()))
	}

	tags := strings.Split(param.String(r, "tags", ""), "\n")
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

	s.Creator = username
	s.Url = url
	s.ExpectCode = param.String(r, "expect_code", "200")
	s.Timeout = param.Int(r, "timeout", 3000)
	s.MaxStep = param.Int(r, "max_step", 3)
	s.Teams = param.String(r, "teams", "")
	s.Times = param.Int(r, "times", 3)
	s.Note = param.String(r, "note", "")
	s.Keywords = param.String(r, "keywords", "")
	s.Data = param.String(r, "data", "")
	s.MonitorIdc = param.String(r, "monitor_idc", "")
	s.Tag = tagStr

	err = s.Update()
	render.AutoJSON(w, err)
}

func DeleteStrategy(w http.ResponseWriter, r *http.Request) {
	me := MeRequired(LoginRequired(w, r))
	strategy := StraRequired(r)
	//teams := strings.Split(strategy.Teams, ",")

	username := me.Name
	if strategy.Creator != username && strategy.Creator != "qinyening" {
		errors.Panic("没有权限")
	}

	err := strategy.Delete()
	render.AutoJSON(w, err)
}

func GetTeamsOfStrategy(w http.ResponseWriter, r *http.Request) {
	MeRequired(LoginRequired(w, r))
	stra := StraRequired(r)
	teams, err := model.GetTeamsByIds(stra.Teams)
	render.AutoJSON(w, err, map[string]interface{}{"teams": teams})
}
