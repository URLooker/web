package handler

import (
	"net/http"

	"github.com/peng19940915/urlooker/web/g"
	"github.com/peng19940915/urlooker/web/http/errors"
	"github.com/peng19940915/urlooker/web/http/render"
	"github.com/peng19940915/urlooker/web/utils"
	"github.com/gin-gonic/gin"
)

func GetLog(c *gin.Context) {

	AdminRequired(LoginRequired(c))
	appLog, err := utils.ReadLastLine("var/app.log")
	errors.MaybePanic(err)
	render.HTML(http.StatusOK, c,"status/log", gin.H{
		"Log": appLog,
	})
}

func Version(c *gin.Context) {
	render.Data(c, g.VERSION)
}
