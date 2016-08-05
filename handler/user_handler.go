package handler

import (
	"fmt"
	"net/http"

	"github.com/toolkits/str"
	//"github.com/toolkits/web"

	"github.com/urlooker/web/http/cookie"
	"github.com/urlooker/web/http/errors"
	"github.com/urlooker/web/http/param"
	"github.com/urlooker/web/http/render"
	"github.com/urlooker/web/model"
	"github.com/urlooker/web/utils"
)

func Register(w http.ResponseWriter, r *http.Request) {
	username := param.MustString(r, "username")
	password := param.MustString(r, "password")
	repeat := param.MustString(r, "repeat")

	if password != repeat {
		errors.Panic("两次输入的密码不一致")
	}

	if str.HasDangerousCharacters(username) {
		errors.Panic("用户名不合法，请不要使用非法字符")
	}

	userid, err := model.UserRegister(username, utils.EncryptPassword(password))
	errors.MaybePanic(err)

	render.AutoJSON(w, cookie.WriteUser(w, userid, username))
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	render.Data(r, "Title", "register")
	render.Data(r, "callback", param.String(r, "callback", "/"))
	render.HTML(r, w, "auth/register")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	errors.MaybePanic(cookie.RemoveUser(w))
	http.Redirect(w, r, "/", 302)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	render.Data(r, "Title", "login")
	render.Data(r, "callback", param.String(r, "callback", "/"))
	render.HTML(r, w, "auth/login")
}

func Login(w http.ResponseWriter, r *http.Request) {
	username := param.MustString(r, "username")
	password := param.MustString(r, "password")

	if str.HasDangerousCharacters(username) {
		errors.Panic("用户名不合法，请不要使用非法字符")
	}

	userid, err := model.UserLogin(username, utils.EncryptPassword(password))
	errors.MaybePanic(err)

	render.AutoJSON(w, cookie.WriteUser(w, userid, username))
}

func MeJson(w http.ResponseWriter, r *http.Request) {
	render.AutoJSON(w, nil, MeRequired(LoginRequired(w, r)))
}

func UsersJson(w http.ResponseWriter, r *http.Request) {
	MeRequired(LoginRequired(w, r))
	query := param.String(r, "query", "")
	limit := param.Int(r, "limit", 10)
	if str.HasDangerousCharacters(query) {
		render.AutoJSON(w, fmt.Errorf("query invalid"), nil)
		return
	}

	users, err := model.QueryUsers(query, limit)
	errors.MaybePanic(err)

	render.JSON(w, map[string]interface{}{"users": users})
}

func UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	me := MeRequired(LoginRequired(w, r))

	cnname := param.String(r, "cnname", "")
	email := param.String(r, "email", "")
	phone := param.String(r, "phone", "")
	wechat := param.String(r, "wechat", "")

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

	me.Email = email
	me.Phone = phone
	me.Wechat = wechat
	errors.MaybePanic(me.UpdateProfile())
	render.AutoJSON(w, nil)
}

func ChangeMyPasswd(w http.ResponseWriter, r *http.Request) {

	me := MeRequired(LoginRequired(w, r))

	oldPasswd := param.MustString(r, "old_password")
	newPasswd := param.MustString(r, "new_password")
	repeat := param.MustString(r, "repeat")

	if newPasswd != repeat {
		errors.Panic("两次输入的密码不一致")
	}

	err := me.ChangePasswd(utils.EncryptPassword(oldPasswd), utils.EncryptPassword(newPasswd))
	if err == nil {
		cookie.RemoveUser(w)
	}

	render.AutoJSON(w, err)
}
