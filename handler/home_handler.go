package handler

import (
	"net/http"

	"github.com/toolkits/str"
	"github.com/toolkits/web"

	"github.com/peng19940915/urlooker/web/http/errors"
	"github.com/peng19940915/urlooker/web/http/param"
	"github.com/peng19940915/urlooker/web/model"
	"github.com/gin-gonic/gin"
	"github.com/peng19940915/urlooker/web/http/render"
)

func HomeIndex(c *gin.Context) {
	me := MeRequired(LoginRequired(c))
	username := me.Name
	mine := param.Int(c.Request, "mine", 1)
	query := param.String(c.Request, "q", "")
	if str.HasDangerousCharacters(query) {
		errors.Panic("查询字符不合法")
	}

	limit := param.Int(c.Request, "limit", 10)
	total, err := model.GetAllStrategyCount(mine, query, username)
	errors.MaybePanic(err)
	pager := web.NewPaginator(c.Request, limit, total)

	strategies, err := model.GetAllStrategy(mine, limit, pager.Offset(), query, username)

	errors.MaybePanic(err)
	render.HTML(http.StatusOK, c, "home/index", gin.H{
		"Strategies": strategies,
		"Pager": pager,
		"Mine": mine,
		"Query": query,
	})
}

func QueryPortScan(c *gin.Context) {
	me := MeRequired(LoginRequired(c))
	username := me.Name
	mine := param.Int(c.Request, "mine", 1)
	query := param.String(c.Request, "q", "")
	if str.HasDangerousCharacters(query) {
		errors.Panic("查询字符不合法")
	}

	limit := param.Int(c.Request, "limit", 10)
	total, err := model.GetAllPortStrategyCount(mine, query, username)
	errors.MaybePanic(err)
	pager := web.NewPaginator(c.Request, limit, total)

	strategies, err := model.GetAllPortStrategy(mine, limit, pager.Offset(), query, username)

	errors.MaybePanic(err)
	render.HTML(http.StatusOK, c, "home/tcp_port_scan", gin.H{
		"Strategies": strategies,
		"Pager": pager,
		"Mine": mine,
		"Query": query,
	})
}