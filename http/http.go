package http

import (

	"github.com/gin-gonic/gin"
	"html/template"
	"github.com/peng19940915/urlooker/web/utils"
	"github.com/peng19940915/urlooker/web/http/middleware"
	"github.com/peng19940915/urlooker/web/http/helper"
	"net/http"
)


func StartGin(port string, r *gin.Engine) {
	r.Use(utils.CORS())
	r.SetFuncMap(template.FuncMap{
		"Times1000":       helper.Times1000,
		"MailsOfTeam":     helper.MailsOfTeam,
		"TeamsOfStrategy": helper.TeamsOfStrategy,
		"HumenTime":       helper.HumenTime,
		"GetFirst":        helper.GetFirst,
	})
	r.LoadHTMLGlob("./views/**/*")
	//r.Delims("{[{", "}]}")
	r.StaticFS("/lib", http.Dir("./static/lib"))
	r.StaticFS("/js",  http.Dir("./static/js"))
	r.StaticFS("/img", http.Dir("./static/img"))
	r.StaticFS("/css", http.Dir("./static/css"))

	r.Use(middleware.Recovery())

	ConfigRouter(r)
	r.Run(port)

}