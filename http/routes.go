package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urlooker/web/handler"
)

func ConfigRouter(r *mux.Router) {
	configStraRoutes(r)
	configStaticRoutes(r)
	configApiRoutes(r)
	configDomainRoutes(r)
	configAuthRoutes(r)
	configUserRoutes(r)
	configTeamRoutes(r)
}

func configDomainRoutes(r *mux.Router) {
	r.HandleFunc("/url", handler.UrlStatus).Methods("GET")
}

func configChartRoutes(r *mux.Router) {
	r.HandleFunc("/chart", handler.UrlStatus).Methods("GET")
}

func configApiRoutes(r *mux.Router) {
	r.HandleFunc("/api/item/{hostname}", handler.GetHostIpItem).Methods("GET")
}

func configStraRoutes(r *mux.Router) {
	r.HandleFunc("/", handler.HomeIndex).Methods("GET")
	r.HandleFunc("/strategy/add", handler.AddStrategyGet).Methods("GET")
	r.HandleFunc("/strategy/add", handler.AddStrategyPost).Methods("POST")
	r.HandleFunc("/strategy/{id:[0-9]+}", handler.GetStrategyById).Methods("POST")
	r.HandleFunc("/strategy/{id:[0-9]+}/delete", handler.DeleteStrategy).Methods("POST")
	r.HandleFunc("/strategy/{id:[0-9]+}/edit", handler.UpdateStrategyGet).Methods("GET")
	r.HandleFunc("/strategy/{id:[0-9]+}/edit", handler.UpdateStrategy).Methods("POST")
	r.HandleFunc("/strategy/{id:[0-9]+}/teams", handler.GetTeamsOfStrategy).Methods("GET")
}

func configAuthRoutes(r *mux.Router) {
	r.HandleFunc("/auth/register", handler.RegisterPage).Methods("GET")
	r.HandleFunc("/auth/register", handler.Register).Methods("POST")
	r.HandleFunc("/auth/logout", handler.Logout).Methods("GET")
	r.HandleFunc("/auth/login", handler.Login).Methods("POST")
	r.HandleFunc("/auth/login", handler.LoginPage).Methods("GET")
}

func configUserRoutes(r *mux.Router) {
	r.HandleFunc("/me.json", handler.MeJson).Methods("GET")
	r.HandleFunc("/me/profile", handler.UpdateMyProfile).Methods("POST")
	r.HandleFunc("/me/chpwd", handler.ChangeMyPasswd).Methods("POST")
	r.HandleFunc("/users/query", handler.UsersJson).Methods("GET")
}

func configTeamRoutes(r *mux.Router) {
	r.HandleFunc("/teams", handler.TeamsPage).Methods("GET")
	r.HandleFunc("/teams/query", handler.TeamsJson).Methods("GET")
	r.HandleFunc("/team/create", handler.CreateTeamGet).Methods("GET")
	r.HandleFunc("/team/create", handler.CreateTeamPost).Methods("POST")
	r.HandleFunc("/team/{tid:[0-9]+}/edit", handler.UpdateTeamGet).Methods("GET")
	r.HandleFunc("/team/{tid:[0-9]+}/edit", handler.UpdateTeamPost).Methods("POST")
	r.HandleFunc("/team/{tid:[0-9]+}/users", handler.GetUsersOfTeam).Methods("GET")
}

func configStaticRoutes(r *mux.Router) {
	r.PathPrefix("/css").Handler(http.FileServer(http.Dir("./static")))
	r.PathPrefix("/fonts").Handler(http.FileServer(http.Dir("./static")))
	r.PathPrefix("/js").Handler(http.FileServer(http.Dir("./static")))
	r.PathPrefix("/img").Handler(http.FileServer(http.Dir("./static")))
	r.PathPrefix("/lib").Handler(http.FileServer(http.Dir("./static")))
}
