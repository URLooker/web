package render

import (
	"net/http"
	"github.com/peng19940915/urlooker/web/http/cookie"
	"github.com/unrolled/render"
	"github.com/gin-gonic/gin"
)


func HTML(status int, c *gin.Context, name string, params gin.H, htmlOpt ...render.HTMLOptions) {
	userid, username, found := cookie.ReadUser(c)
	params["HasLogin"] = found
	params["UserId"] = userid
	params["UserName"] = username
	c.HTML(status, name, params)
}

func MaybeError(c *gin.Context, err error) {
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}

func Data(c *gin.Context, v interface{}, msg ...string) {
	m := ""
	if len(msg) > 0 {
		m = msg[0]
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": m,
		"data": v,
	})
}
