package handler

import (
	"github.com/peng19940915/urlooker/web/g"
	"github.com/peng19940915/urlooker/web/http/render"
	"github.com/gin-gonic/gin"
)

func Version(c *gin.Context) {
	render.Data(c, g.VERSION)
}
