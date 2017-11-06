package handler

import (
	"net/http"

	"github.com/urlooker/web/g"
	"github.com/urlooker/web/http/errors"
	"github.com/urlooker/web/http/render"
	"github.com/urlooker/web/utils"
)

func GetLog(w http.ResponseWriter, r *http.Request) {

	AdminRequired(LoginRequired(w, r))
	appLog, err := utils.ReadLastLine("var/app.log")
	errors.MaybePanic(err)

	render.Put(r, "Log", appLog)

	render.HTML(r, w, "status/log")
}

func Version(w http.ResponseWriter, r *http.Request) {
	render.Data(w, g.VERSION)
}
