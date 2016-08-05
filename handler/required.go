package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/toolkits/str"
	"github.com/urlooker/web/http/cookie"
	"github.com/urlooker/web/http/errors"
	"github.com/urlooker/web/model"
)

func StraRequired(r *http.Request) *model.Strategy {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	errors.MaybePanic(err)

	obj, err := model.GetStrategyById(id)
	errors.MaybePanic(err)
	if obj == nil {
		panic(errors.BadRequestError("no such item"))
	}
	return obj
}

func HostnameRequired(r *http.Request) string {
	vars := mux.Vars(r)
	hostname := vars["hostname"]

	if str.HasDangerousCharacters(hostname) {
		errors.Panic("hostname不合法")
	}

	return hostname
}

func LoginRequired(w http.ResponseWriter, r *http.Request) (int64, string) {
	userid, username, found := cookie.ReadUser(r)
	if !found {
		panic(errors.NotLoginError())
	}

	return userid, username
}

func MeRequired(id int64, name string) *model.User {
	user, err := model.GetUserById(id)
	if err != nil {
		panic(errors.InternalServerError(err.Error()))
	}

	if user == nil {
		panic(errors.NotLoginError())
	}

	return user
}

func TeamRequired(r *http.Request) *model.Team {
	vars := mux.Vars(r)
	tid, err := strconv.ParseInt(vars["tid"], 10, 64)
	errors.MaybePanic(err)

	team, err := model.GetTeamById(tid)
	errors.MaybePanic(err)
	if team == nil {
		panic(errors.BadRequestError("no such team"))
	}

	return team
}

func UserMustBeMemberOfTeam(uid, tid int64) {
	is, err := model.IsMemberOfTeam(uid, tid)
	errors.MaybePanic(err)
	if is {
		return
	}

	team, err := model.GetTeamById(tid)
	errors.MaybePanic(err)
	if team != nil && team.Creator == uid {
		return
	}

	panic(errors.BadRequestError("用户不是团队的成员"))
}
