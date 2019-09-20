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
			if flag, tgtTiket:=verifyTicket(ticket); flag{
				userEmail, cnName := getUserInfo(tgtTiket)
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

func getUserInfo(ticket string) (loginEmail string, chineseName string) {
	url := g.Config.SSO.ServerUrl+"/api/v2/info"
	client := http.Client{}

	req, _ := http.NewRequest("GET",url, nil)
	req.Header.Set("ticket", ticket)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("",err)
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	jsonObj, err := simplejson.NewJson(b)
	if err != nil{
		log.Println("tran to json failed when get user info, detail: ",err, "data:",string(b))
	}
	infoData := jsonObj.Get("data")
	loginEmail, err = infoData.Get("LoginEmail").String()
	if err != nil{
		log.Println("get login email failed, detail:", err)
	}
	chineseName, err = infoData.Get("DisplayName").String()
	if err != nil {
		log.Println("get display name failed,detail: ",err)
	}

	return loginEmail, chineseName
}

func verifyTicket(ticket string) (bool, string){
	validateUrl := g.Config.SSO.ServerUrl+"/api/v2/validate"
	client := http.Client{}
	req, _ := http.NewRequest("GET",validateUrl, nil)

	req.Header.Set("s-ticket", ticket)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	b, _:= ioutil.ReadAll(resp.Body)

	jsonObj, err := simplejson.NewJson(b)
	if err != nil{
		log.Println("to json failed "+err.Error())
		log.Println(string(b))
		return false, ""
	}
	errorCode,err := jsonObj.Get("errorCode").Int()
	if err != nil {
		log.Println("get errorCode failed:", err)
		return false, ""
	}
	if errorCode == 0 || errorCode == 1 {
		log.Println("verfyTicket success.")
		tgtTicket, err := jsonObj.Get("data").Get("Ticket").String()
		if err != nil {
			log.Println("get tgt ticket failed, detail: %s", err.Error())
			return false, ""
		}
		return true, tgtTicket
	}else{
		log.Println("verfy ticket failed:")
	}

	return false, ""
}
