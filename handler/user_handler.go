package handler

import (
	"net/http"
	"strings"

	"github.com/toolkits/str"

	"github.com/peng19940915/urlooker/web/g"
	"github.com/peng19940915/urlooker/web/http/cookie"
	"github.com/peng19940915/urlooker/web/http/errors"
	"github.com/peng19940915/urlooker/web/http/param"
	"github.com/peng19940915/urlooker/web/http/render"
	"github.com/peng19940915/urlooker/web/model"
	"github.com/peng19940915/urlooker/web/utils"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	if g.Config.Ldap.Enabled {
		errors.Panic("注册已关闭")
	}
	username := param.MustString(c.Request, "username")
	password := param.MustString(c.Request, "password")
	repeat := param.MustString(c.Request, "repeat")

	if password != repeat {
		errors.Panic("两次输入的密码不一致")
	}

	if str.HasDangerousCharacters(username) {
		errors.Panic("用户名不合法，请不要使用非法字符")
	}

	userid, err := model.UserRegister(username, utils.EncryptPassword(password))
	errors.MaybePanic(err)

	render.Data(c, cookie.WriteUser(c, userid, username))
}

func RegisterPage(c *gin.Context) {
	render.HTML(http.StatusOK, c,"auth/register", gin.H{
		"Title": "register",
		"callback": param.String(c.Request, "callback", "/"),
	})
}

func Logout(c *gin.Context) {
	errors.MaybePanic(cookie.RemoveUser(c))
	c.Redirect(http.StatusFound, "/")
}

func LoginPage(c *gin.Context) {
	render.HTML(http.StatusOK, c,"auth/login", gin.H{
		"Title": "login",
		"callback": param.String(c.Request, "callback", "/"),
	})
}

func Login(c *gin.Context) {
	username := param.MustString(c.Request, "username")
	password := param.MustString(c.Request, "password")

	if str.HasDangerousCharacters(username) {
		errors.Panic("用户名不合法，请不要使用非法字符")
	}

	var u *model.User
	var userid int64
	if g.Config.Ldap.Enabled {
		sucess, err := utils.LdapBind(g.Config.Ldap.Addr,
			g.Config.Ldap.BaseDN,
			g.Config.Ldap.BindDN,
			g.Config.Ldap.BindPasswd,
			g.Config.Ldap.UserField,
			username,
			password)

		errors.MaybePanic(err)
		if !sucess {
			errors.Panic("name or password error")
			return
		}

		userAttributes, err := utils.Ldapsearch(g.Config.Ldap.Addr,
			g.Config.Ldap.BaseDN,
			g.Config.Ldap.BindDN,
			g.Config.Ldap.BindPasswd,
			g.Config.Ldap.UserField,
			username,
			g.Config.Ldap.Attributes)
		userSn := ""
		userMail := ""
		userTel := ""
		if err == nil {
			userSn = userAttributes["sn"]
			userMail = userAttributes["mail"]
			userTel = userAttributes["telephoneNumber"]
		}

		arr := strings.Split(username, "@")
		var userName, userEmail string
		if len(arr) == 2 {
			userName = arr[0]
			userEmail = username
		} else {
			userName = username
			userEmail = userMail
		}

		u, err = model.GetUserByName(userName)
		errors.MaybePanic(err)
		if u == nil {
			// 说明用户不存在
			u = &model.User{
				Name:     userName,
				Password: "",
				Cnname:   userSn,
				Phone:    userTel,
				Email:    userEmail,
			}
			errors.MaybePanic(u.Save())
		}
		userid = u.Id
	} else {
		var err error
		userid, err = model.UserLogin(username, utils.EncryptPassword(password))
		errors.MaybePanic(err)
	}

	render.Data(c, cookie.WriteUser(c, userid, username))
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

func UpdateMyProfile(c *gin.Context) {
	me := MeRequired(LoginRequired(c))

	cnname := param.String(c.Request, "cnname", "")
	email := param.String(c.Request, "email", "")
	phone := param.String(c.Request, "phone", "")
	wechat := param.String(c.Request, "wechat", "")

	if str.HasDangerousCharacters(cnname) {
		errors.Panic("中文名不合法")
	}
	if email != "" && !str.IsMail(email) {
		errors.Panic("邮箱不合法")
	}
	if phone != "" && !str.IsPhone(phone) {
		errors.Panic("手机号不合法")
	}
	if str.HasDangerousCharacters(wechat) {
		errors.Panic("微信不合法")
	}

	me.Cnname = cnname
	me.Email = email
	me.Phone = phone
	me.Wechat = wechat
	errors.MaybePanic(me.UpdateProfile())
	render.Data(c, "ok")
}

func ChangeMyPasswd(c *gin.Context) {

	uid, _ := LoginRequired(c)
	me, err := model.GetUserPwById(uid)
	errors.MaybePanic(err)

	oldPasswd := param.MustString(c.Request, "old_password")
	newPasswd := param.MustString(c.Request, "new_password")
	repeat := param.MustString(c.Request, "repeat")

	if newPasswd != repeat {
		errors.Panic("两次输入的密码不一致")
	}

	err = me.ChangePasswd(utils.EncryptPassword(oldPasswd), utils.EncryptPassword(newPasswd))
	if err == nil {
		cookie.RemoveUser(c)
	}

	errors.MaybePanic(err)
	render.Data(c, "ok")
}
