package viewserver

import (
	"util/jsondb"
	"util"
	"time"
)

type LoginController struct {
	Controller
}

func (c LoginController) jsonfile() string {
	return "database/user.json"
}

func (c LoginController) LoginPage() {
	data:= map[string]interface{}{"aa":123}
	c.RenderView("login.gohtml", data)
}

func (t LoginController) _login(username, userpwd string) map[string]interface{} {
	info := jsondb.QueryRow(t.jsonfile(), map[string]interface{}{"cusername": username, "cuserpwd": userpwd},nil)
	if info == nil {
		return nil
	}

	return info
}

func (c LoginController) Login() {
	post := c.GetPostData()
	loginname := util.String(post["loginname"])
	loginpwd := util.String(post["loginpwd"])

	info := c._login(loginname, loginpwd)
	if info == nil {
		c.RenderError("登录失败", "")
		return
	}
	token := util.Token_Get(map[string]string{"pkid": util.String(info["pkid"]), "loginname": info["cusername"].(string)})

	util.Cookie_SetFasthttp(c.RequestCtx, "token", token, time.Hour*365*24)
	c.RenderJSON(map[string]interface{}{"msg": "登录成功", "token": token, "url": "/html/home"})
	return

}
