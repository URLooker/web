package handler

import (
	"strconv"

	"github.com/toolkits/str"

	"github.com/peng19940915/urlooker/web/g"
	"github.com/peng19940915/urlooker/web/http/cookie"
	"github.com/peng19940915/urlooker/web/http/errors"
	"github.com/peng19940915/urlooker/web/model"
	"github.com/gin-gonic/gin"
)

func StraRequired(c *gin.Context) *model.Strategy {
	idTmp := c.Query("id")
	if idTmp != "" {
		id, err := strconv.ParseInt(idTmp, 10, 64)
		errors.MaybePanic(err)
		obj, err := model.GetStrategyById(id)
		errors.MaybePanic(err)
		if obj == nil {
			panic(errors.BadRequestError("no such item"))
		}
		return obj
	}else {
		panic(errors.BadRequestError("plz make sure strategy is right."))
	}
}

func PortStraRequired(c *gin.Context) *model.PortStrategy {
	idTmp := c.Query("id")

	if idTmp != "" {
		id, err := strconv.ParseInt(idTmp, 10, 64)
		errors.MaybePanic(err)
		obj, err := model.GetPortStrategyById(id)
		errors.MaybePanic(err)
		if obj == nil {
			panic(errors.BadRequestError("no such item"))
		}
		return obj
	}else {
		panic(errors.BadRequestError("plz make sure strategy is right."))
	}
}

func HostnameRequired(c *gin.Context) string {
	hostname := c.Query("hostname")

	if hostname != "" {
		if str.HasDangerousCharacters(hostname) {
			errors.Panic("hostname不合法")
		}
		return hostname
	}else {
		panic(errors.BadRequestError("hostname is null"))
	}
}

func LoginRequired(c *gin.Context) (int64, string) {
	userid, username, _,found := cookie.ReadUser(c)
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

func TeamRequired(c *gin.Context) *model.Team {
	tidTmp := c.Query("tid")
	if tidTmp == "" {
		panic(errors.BadRequestError("tid is null"))
	}
	tid, err := strconv.ParseInt(tidTmp, 10, 64)
	errors.MaybePanic(err)

	team, err := model.GetTeamById(tid)
	errors.MaybePanic(err)
	if team == nil {
		panic(errors.BadRequestError("no such team"))
	}

	return team
}

func IsAdmin(username string) bool {
	for _, admin := range g.Config.Admins {
		if username == admin {
			return true
		}
	}
	return false
}
