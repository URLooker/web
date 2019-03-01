package handler

import (
	"github.com/peng19940915/urlooker/web/g"
	"github.com/peng19940915/urlooker/web/http/render"
	"github.com/gin-gonic/gin"
)

func GetHostIpItem(c *gin.Context) {
	hostname := HostnameRequired(c)
	ipItem, exists := g.DetectedItemMap.Get(hostname)
	if !exists {
		render.Data(c, "")
		return
	}
	render.Data(c, ipItem)
}
