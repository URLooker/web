package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/toolkits/str"
	"github.com/toolkits/web"

	"github.com/peng19940915/urlooker/web/http/errors"
	"github.com/peng19940915/urlooker/web/http/param"
	"github.com/peng19940915/urlooker/web/http/render"
	"github.com/peng19940915/urlooker/web/model"
	"github.com/gin-gonic/gin"
)

func TeamsPage(c *gin.Context) {
	me := MeRequired(LoginRequired(c))

	query := param.String(c.Request, "q", "")
	if str.HasDangerousCharacters(query) {
		errors.Panic("查询字符不合法")
	}
	limit := param.Int(c.Request, "limit", 10)

	total, err := model.TeamCountOfUser(query, me.Id)
	errors.MaybePanic(err)
	pager := web.NewPaginator(c.Request, limit, total)
	teams, err := model.TeamsOfUser(query, me.Id, limit, pager.Offset())
	errors.MaybePanic(err)
	for _, team := range teams {
		user, err := model.GetUserById(team.Creator)
		if err == nil && user != nil {
			team.CreatorName = user.Name
		}
	}
	render.HTML(http.StatusOK, c, "team/index", gin.H{
		"Teams": teams,
		"Query": query,
		"Pager": pager,
		"Me": me,
		"Title": "Team",
	})
}

func TeamsJson(c *gin.Context) {
	MeRequired(LoginRequired(c))

	query := param.String(c.Request, "query", "")
	limit := param.Int(c.Request, "limit", 10)

	if str.HasDangerousCharacters(query) {
		render.Data(c, fmt.Errorf("query invalid"))
		return
	}

	teams, err := model.QueryTeams(query, limit)
	errors.MaybePanic(err)

	render.Data(c, map[string]interface{}{"teams": teams})
}

func CreateTeamGet(c *gin.Context) {
	me := MeRequired(LoginRequired(c))

	c.HTML(http.StatusOK, "team/create", gin.H{
		"Me": me,
		"Title": "Team",
	})
}

func CreateTeamPost(c *gin.Context) {
	me := MeRequired(LoginRequired(c))

	name := param.MustString(c.Request, "name")
	if str.HasDangerousCharacters(name) {
		errors.Panic("team名称不合法")
	}
	resume := param.String(c.Request, "resume", "")
	if str.HasDangerousCharacters(resume) {
		errors.Panic("resume不合法")
	}

	uidsStr := param.String(c.Request, "users", "")
	if str.HasDangerousCharacters(uidsStr) {
		errors.Panic("users不合法")
	}
	uidSlice := strings.Split(uidsStr, ",")

	isci := false
	uids := make([]int64, 0)
	for _, u := range uidSlice {
		if u == "" {
			continue
		}
		uid, err := strconv.ParseInt(u, 10, 64)
		errors.MaybePanic(err)
		uids = append(uids, uid)
		if uid == me.Id {
			isci = true
		}
	}
	if !isci {
		// creator is member of team
		uids = append(uids, me.Id)
	}

	_, err := model.AddTeam(name, resume, me.Id, uids)
	render.MaybeError(c, err)
}

func UpdateTeamGet(c *gin.Context) {
	team := TeamRequired(c)
	me := MeRequired(LoginRequired(c))
	if !IsAdmin(me.Name) {
		UserMustBeMemberOfTeam(me.Id, team.Id)
	}

	uids := make([]string, 0)
	users, err := model.UsersOfTeam(team.Id)
	errors.MaybePanic(err)
	for _, user := range users {
		uids = append(uids, strconv.FormatInt(user.Id, 10))
	}
	render.HTML(http.StatusOK, c,"team/edit",gin.H{
		"Team": team,
		"Uids": strings.Join(uids, ","),
		"Me": me,
		"Title": "Team",
	})
}

func UpdateTeamPost(c *gin.Context) {
	me := MeRequired(LoginRequired(c))
	team := TeamRequired(c)
	if !IsAdmin(me.Name) {
		UserMustBeMemberOfTeam(me.Id, team.Id)
	}

	team.Resume = param.String(c.Request, "resume", "")
	if str.HasDangerousCharacters(team.Resume) {
		errors.Panic("resume不合法")
	}
	uidsStr := param.String(c.Request, "users", "")
	if str.HasDangerousCharacters(uidsStr) {
		errors.Panic("users不合法")
	}
	uidsSlice := strings.Split(uidsStr, ",")
	uids := make([]int64, 0)
	for _, uidStr := range uidsSlice {
		if uidStr == "" {
			continue
		}
		uid, err := strconv.ParseInt(uidStr, 10, 64)
		errors.MaybePanic(err)
		uids = append(uids, uid)
	}

	render.Data(c, team.Update(uids))
}

func GetUsersOfTeam(c *gin.Context) {
	MeRequired(LoginRequired(c))
	team := TeamRequired(c)

	users, err := model.UsersOfTeam(team.Id)
	errors.MaybePanic(err)
	render.Data(c, users)
}
