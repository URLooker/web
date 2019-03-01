package middleware

import (
	"github.com/gin-gonic/gin"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/peng19940915/urlooker/web/http/errors"
	"net/http"
	"github.com/peng19940915/urlooker/web/http/render"
)


func Recovery() gin.HandlerFunc {
	return func (context *gin.Context) {
		defer func(){
			if err := recover(); err != nil {
				if customError, ok := err.(errors.Error); ok {

					msg := fmt.Sprintf("[%s:%d] %s %s", customError.File, customError.Line, customError.Time, customError.Msg)
					log.Error(msg)
					if isAjax(context) {
						context.JSON(http.StatusOK, map[string]string{"msg": customError.Msg})
						return
					}

					if customError.Code == http.StatusUnauthorized || customError.Code == http.StatusForbidden {
						context.Redirect(http.StatusFound, "auth/login")
						return
					}
					render.HTML(http.StatusInternalServerError, context,"inc/error",gin.H{})
					return
				}
				//stack := make([]byte, 1024 * 8)
				//stack = stack[:runtime.Stack(stack, false)]
				//context.AbortWithStatus(500)
				//log.Errorf("[Recovery] %s panic recovered:\n%s\n%s", err, stack)
			}
		}()

		context.Next()
	}
}

func isAjax(c *gin.Context) bool {
	return c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest"
}
