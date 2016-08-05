package handler

import (
	"log"
	"net/http"

	"github.com/urlooker/web/g"
	"github.com/urlooker/web/http/render"
)

func GetHostIpItem(w http.ResponseWriter, r *http.Request) {
	hostname := HostnameRequired(r)
	ipItem, exists := g.DetectedItemMap.Get(hostname)
	log.Println(ipItem)
	if !exists {
		render.JSON(w, "")
		return
	}
	render.JSON(w, ipItem)
}
