package handler

import (
	"net/http"

	"github.com/urlooker/web/g"
	"github.com/urlooker/web/http/render"
)

func GetHostIpItem(w http.ResponseWriter, r *http.Request) {
	hostname := HostnameRequired(r)
	ipItem, exists := g.DetectedItemMap.Get(hostname)
	if !exists {
		render.JSON(w, "")
		return
	}
	render.JSON(w, ipItem)
}
