package http

import (
	"github.com/peng19940915/urlooker/web/handler"
	"github.com/gin-gonic/gin"
	"path/filepath"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func ConfigRouter(r *gin.Engine) {

	configStraRoutes(r)
	configApiRoutes(r)
	configDomainRoutes(r)
	configAuthRoutes(r)
	configUserRoutes(r)
	configTeamRoutes(r)
	configPortStraRoutes(r)
	configProcRoutes(r)
	configPortDomainRoutes(r)
}

func configDomainRoutes(r *gin.Engine) {
	r.GET("/url", handler.UrlStatus)
}

func configPortDomainRoutes(r *gin.Engine) {
	r.GET("/port", handler.PortStatus)
}
func configChartRoutes(r *gin.Engine) {
	r.GET("/chart", handler.UrlStatus)
}

func configApiRoutes(r *gin.Engine) {
	r.GET("/api/item/{hostname}", handler.GetHostIpItem)
}

func configStraRoutes(r *gin.Engine) {
	r.GET("/", handler.HomeIndex)
	r.GET("/strategy/add", handler.AddStrategyGet)
	r.POST("/strategy/add", handler.AddStrategyPost)
	r.POST("/strategy", handler.GetStrategyById)
	r.POST("/strategy/delete", handler.DeleteStrategy)
	r.GET("/strategy/edit", handler.UpdateStrategyGet)
	r.POST("/strategy/edit", handler.UpdateStrategy)
	r.GET("/strategy/teams", handler.GetTeamsOfStrategy)
}

func configPortStraRoutes(r *gin.Engine){
	r.GET("/tcp_port_scan", handler.QueryPortScan)
	r.GET("/port_strategy/add", handler.AddPortStrategyGet)
	r.POST("/port_strategy/add", handler.AddPortStrategyPost)

	r.POST("/port_strategy", handler.GetPortStrategyById)
	r.POST("/port_strategy/delete", handler.DeletePortStrategy)
	r.GET("/port_strategy/edit", handler.UpdatePortStrategyGet)
	r.POST("/port_strategy/edit", handler.UpdatePortStrategy)
	r.GET("/port_strategy/teams", handler.GetPortTeamsOfStrategy)
}

func configAuthRoutes(r *gin.Engine) {
	r.GET("/auth/register", handler.RegisterPage)
	r.POST("/auth/register", handler.Register)

	r.GET("/auth/logout", handler.Logout)
	r.POST("/auth/login", handler.Login)
	r.GET("/auth/login", handler.LoginPage)
}

func configUserRoutes(r *gin.Engine) {
	r.GET("/me.json", handler.MeJson)
	r.POST("/me/profile", handler.UpdateMyProfile)
	r.POST("/me/chpwd", handler.ChangeMyPasswd)
	r.GET("/users/query", handler.UsersJson)
}

func configTeamRoutes(r *gin.Engine) {
	r.GET("/teams", handler.TeamsPage)
	r.GET("/teams/query", handler.TeamsJson)
	r.GET("/team/create", handler.CreateTeamGet)
	r.POST("/team/create", handler.CreateTeamPost)
	r.GET("/team/{tid:[0-9]+}/edit", handler.UpdateTeamGet)
	r.POST("/team/{tid:[0-9]+}/edit", handler.UpdateTeamPost)
	r.GET("/team/{tid:[0-9]+}/users", handler.GetUsersOfTeam)
}

func configProcRoutes(r *gin.Engine) {
	//r.HandleFunc("/log", handler.GetLog).Methods("GET")
	r.GET("/version", handler.Version)
}
func ConfigStaticRoutes(r *gin.Engine) {
	r.Static("/css", "./static")
	r.Static("/fonts","./static")
	r.Static("/js","./static")
	r.Static("/img", "./static")
	r.Static("/lib", "./static")
}
/*
func configStaticRoutes(r *gin.Engine) {
	r.GET("/css*", http.FileServer(http.Dir("./static"))
	r.GET("/fonts*", http.FileServer(http.Dir("./static")))
	r.GET("/js*").(http.FileServer(http.Dir("./static")))
	r.GET("/img*").Handler(http.FileServer(http.Dir("./static")))
	r.GET("/lib*").Handler(http.FileServer(http.Dir("./static")))
}
*/
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Error(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}