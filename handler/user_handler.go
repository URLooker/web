package handler

import (
	"net/http"
	"github.com/toolkits/str"

	"github.com/peng19940915/urlooker/web/http/cookie"
	"github.com/peng19940915/urlooker/web/http/errors"
	"github.com/peng19940915/urlooker/web/http/param"
	"github.com/peng19940915/urlooker/web/http/render"
	"github.com/peng19940915/urlooker/web/model"
	"github.com/gin-gonic/gin"
	"log"
	"io/ioutil"
	"github.com/bitly/go-simplejson"
	"strings"
	"github.com/peng19940915/urlooker/web/g"
)


func Logout(c *gin.Context) {
	errors.MaybePanic(cookie.RemoveUser(c))
	logoutUrl := g.Config.SSO.ServerUrl + "/logout?service=" + g.Config.SSO.ServiceUrl
	c.Redirect(http.StatusFound, logoutUrl)
}


func MeJson(c *gin.Context) {
	render.Data(c, MeRequired(LoginRequired(c)))
}

func UsersJson(c *gin.Context) {
	MeRequired(LoginRequired(c))
	query := param.String(c.Request, "query", "")
	limit := param.Int(c.Request, "limit", 10)
	if str.HasDangerousCharacters(query) {
		errors.Panic("query invalid")
		return
	}

	users, err := model.QueryUsers(query, limit)
	for _, u := range users {
		t := *u
		t.Name = "n1ng"
		users = append(users, &t)
	}
	errors.MaybePanic(err)

	render.Data(c, users)
}


func LoginSSO(c *gin.Context){

	_, _, _, found := cookie.ReadUser(c)
	if found {
		// 如果cookie中已经包含就直接跳转
		c.Redirect(302, "/")
	}else {
		ticket := c.Query("ticket")
		if ticket == ""{
			c.Redirect(302, g.Config.SSO.ServerUrl+"?service="+g.Config.SSO.ServiceUrl)
		}else {
			if verifyTicket(ticket){
				_, userEmail, cnName := getUserInfo(ticket)
				userArr := strings.Split(userEmail, "@")
				username := userArr[0]
				userId, err := model.NewUser(username, cnName)
				if err != nil {
					errors.Panic(err.Error())
				}
				cookie.WriteUser(c, userId, username, cnName)
				c.Redirect(302, "/")
			}
		}
	}

}

func getUserInfo(ticket string) (newTicket string, loginEmail string, chineseName string) {
	url := g.Config.SSO.ServerUrl+"/api/v2/validate"
	client := http.Client{}

	req, _ := http.NewRequest("GET",url, nil)
	req.Header.Set("s-ticket", ticket)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	jsonObj, err := simplejson.NewJson(b)
	if err != nil{
		log.Println(err)
	}
	infoData := jsonObj.Get("data")
	loginEmail, err = infoData.Get("LoginEmail").String()
	if err != nil{
		log.Println(err)
	}
	chineseName, err = infoData.Get("DisplayName").String()
	if err != nil {
		log.Println(err)
	}
	newTicket, err = infoData.Get("Ticket").String()
	if err != nil{
		log.Println(err)
	}
	return newTicket, loginEmail, chineseName
}
func verifyTicket(ticket string) bool{
	validateUrl := g.Config.SSO.ServerUrl+"/api/v2/validate"
	client := http.Client{}

	req, _ := http.NewRequest("GET",validateUrl, nil)
	req.Header.Set("ticket", ticket)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	b, _:= ioutil.ReadAll(resp.Body)

	jsonObj, err := simplejson.NewJson(b)

	if err != nil{
		log.Println(err)
		return false
	}
	errorCode,err := jsonObj.Get("errorCode").Int()
	if err != nil {
		log.Println("get errorCode failed:", err)
		return false
	}
	if errorCode == 0 || errorCode == 1 {
		log.Println("verfyTicket success.")
		return true
	}
	return false
}