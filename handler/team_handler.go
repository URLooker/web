package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/toolkits/str"
	"github.com/toolkits/web"

	"github.com/peng19940915/urlooker/web/http/errors"
	"github.com/peng19940915/urlooker/web/http/param"
	"github.com/peng19940915/urlooker/web/http/render"
	"github.com/peng19940915/urlooker/web/model"
	"github.com/gin-gonic/gin"
	"github.com/peng19940915/urlooker/web/utils"
)

func TeamsPage(c *gin.Context) {
	me := MeRequired(LoginRequired(c))
	query := param.String(c.Request, "q", "")
	if str.HasDangerousCharacters(query) {
		errors.Panic("查询字符不合法")
	}
	limit := param.Int(c.Request, "limit", 10)

	total, err := model.TeamCountOfUser(query)
	errors.MaybePanic(err)
	pager := web.NewPaginator(c.Request, limit, total)
	teams, err := model.Teams(query, limit, pager.Offset())
	errors.MaybePanic(err)
	for _, team := range teams {
		user, err := model.GetUserById(team.Creator)
		if err == nil && user != nil {
			if user.Cnname == ""{
				team.CreatorName = user.Name
			}else {
				team.CreatorName = user.Cnname
			}
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

	render.HTML(http.StatusOK, c,"team/create", gin.H{
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

	emailsStr := param.String(c.Request, "emails", "")
	if str.HasDangerousCharacters(emailsStr) {
		errors.Panic("emails不合法")
	}
	emailSlice := strings.Split(emailsStr, ",")
	emails := make([]string, 0)
	for _, u := range emailSlice {
		if u == "" {
			continue
		}
		if !utils.CheckEmail(u){
			errors.Panic(u+": email address is wrong!")
		}
		u = strings.TrimSpace(u)
		emails = append(emails, u)

	}
	_, err := model.AddTeam(name, resume, me.Id, emails)
	render.MaybeError(c, err)
}

func UpdateTeamGet(c *gin.Context) {
	team := TeamRequired(c)
	me := MeRequired(LoginRequired(c))
	if !IsAdmin(me.Name) {

	}

	emails := make([]string, 0)
	relTeamUsers, err := model.MailsOfTeam(team.Id)
	errors.MaybePanic(err)
	for _, user := range relTeamUsers {
		emails = append(emails, user.Email)
	}
	render.HTML(http.StatusOK, c,"team/edit",gin.H{
		"Team": team,
		"Emails": strings.Join(emails, ","),
		"Me": me,
		"Title": "Team",
	})
}

func UpdateTeamPost(c *gin.Context) {
	me := MeRequired(LoginRequired(c))
	team := TeamRequired(c)
	if me.Id != team.Creator && !IsAdmin(me.Name) {
		errors.Panic("权限不足")
	}

	team.Resume = param.String(c.Request, "resume", "")
	if str.HasDangerousCharacters(team.Resume) {
		errors.Panic("resume不合法")
	}
	emailsStr := param.String(c.Request, "emails", "")
	if str.HasDangerousCharacters(emailsStr) {
		errors.Panic("emails不合法")
	}
	emailsSlice := strings.Split(emailsStr, ",")
	emails := make([]string, 0)
	for _, emailStr := range emailsSlice {
		if emailStr == "" {
			continue
		}
		if !utils.CheckEmail(emailStr){
			errors.Panic(emailStr+": email address is wrong!")
		}
		emails = append(emails, emailStr)
	}

	render.Data(c, team.Update(emails))
}

func GetUsersOfTeam(c *gin.Context) {
	MeRequired(LoginRequired(c))
	team := TeamRequired(c)

	users, err := model.MailsOfTeam(team.Id)
	errors.MaybePanic(err)
	render.Data(c, users)
}

func DeleteTeam(c *gin.Context) {
	user := MeRequired(LoginRequired(c))
	team := TeamRequired(c)
	if user.Id != team.Creator && ! IsAdmin(user.Name){
		errors.Panic("权限不足")
	}
	err := model.RemoveTeamById(team.Id)
	if err != nil {
		errors.Panic("delete failed: "+team.Name+"detail: "+err.Error())
	}
	render.Data(c, nil)
}