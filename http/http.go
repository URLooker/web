package http

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	"github.com/urlooker/web/g"
	//"github.com/xpharos/web/http/cookie"
	"github.com/urlooker/web/http/middleware"
	"github.com/urlooker/web/http/render"
)

func Start() {
	render.Init()
	//cookie.Init()

	r := mux.NewRouter().StrictSlash(false)
	ConfigRouter(r)

	n := negroni.New()
	n.Use(middleware.NewRecovery())
	n.UseHandler(r)
	n.Use(middleware.NewLogger())
	n.Run(g.Config.Http.Listen)
}
