package render

import (
	"html/template"
	"net/http"

	"github.com/gorilla/context"
	"github.com/unrolled/render"

	"github.com/urlooker/web/g"
	"github.com/urlooker/web/http/cookie"
	"github.com/urlooker/web/http/helper"
)

var Render *render.Render

var funcMap = template.FuncMap{
	"Times1000":       helper.Times1000,
	"UsersOfTeam":     helper.UsersOfTeam,
	"TeamsOfStrategy": helper.TeamsOfStrategy,
	"HumenTime":       helper.HumenTime,
	"GetFirst":        helper.GetFirst,
}

func Init() {
	debug := g.Config.Debug
	Render = render.New(render.Options{
		Directory:     "views",
		Extensions:    []string{".html"},
		Delims:        render.Delims{"{{", "}}"},
		Funcs:         []template.FuncMap{funcMap},
		IndentJSON:    false,
		IsDevelopment: debug,
	})
}

func Data(r *http.Request, key string, val interface{}) {
	m, ok := context.GetOk(r, "DATA_MAP")
	if ok {
		mm := m.(map[string]interface{})
		mm[key] = val
		context.Set(r, "DATA_MAP", mm)
	} else {
		context.Set(r, "DATA_MAP", map[string]interface{}{key: val})
	}
}

func HTML(r *http.Request, w http.ResponseWriter, name string, htmlOpt ...render.HTMLOptions) {
	userid, username, found := cookie.ReadUser(r)

	Data(r, "Debug", g.Config.Debug)
	Data(r, "HasLogin", found)
	Data(r, "UserId", userid)
	Data(r, "UserName", username)
	Render.HTML(w, http.StatusOK, name, context.Get(r, "DATA_MAP"), htmlOpt...)
}

func JSON(w http.ResponseWriter, v interface{}, statusCode ...int) {
	code := http.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	Render.JSON(w, code, v)
}

func AutoJSON(w http.ResponseWriter, err error, v ...interface{}) {
	if err != nil {
		JSON(w, map[string]interface{}{"msg": err.Error()})
		return
	}

	if len(v) > 0 {
		JSON(w, map[string]interface{}{"msg": "", "data": v[0]})
	} else {
		JSON(w, map[string]interface{}{"msg": ""})
	}
}

func Text(w http.ResponseWriter, v string) {
	Render.Text(w, http.StatusOK, v)
}
