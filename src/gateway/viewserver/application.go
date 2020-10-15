package viewserver

import (
	"util"
	"util/jsondb"
	"strings"
)

type ApplicationController struct {
	Controller
}
func (c ApplicationController) jsonfile()string{
	return "database/application.json"
}

func (c ApplicationController) Getlist() {
	list, _ := jsondb.Query(c.jsonfile(), map[string]interface{}{"cstatus|>=": 0}, map[string]string{"pkid":"desc"})
	c.RenderView("application.gohtml", map[string]interface{}{"list": list})
}

func (c ApplicationController) Getinfo() {
	id := c.GetInt64("id", 0)
	if id == 0 {
		c.RenderError("", "")
		return
	}
	info := jsondb.GetInfo(c.jsonfile(), id)
	c.RenderJSON(map[string]interface{}{"info": info})
	return
}

func (c ApplicationController) Postdata() {
	post := c.GetPostData()
	post["cpath"] = strings.Trim(post["cpath"].(string), "/")
	id := util.String_GetInt64(post["pkid"])
	num := int64(0)
	if id > 0 {
		num, _ = jsondb.UpdateByID(c.jsonfile(), id, post)
	} else {
		num, _ = jsondb.InsertMap(c.jsonfile(), post)
	}
	if num > 0 {
		c.RenderJSON(map[string]interface{}{"msg": "操作成功"})
	} else {
		c.RenderJSON(map[string]interface{}{"msg": "操作失败", "errcode": -1})
	}
}

func (c ApplicationController) DealStatus() {
	post := c.GetPostData()
	id := util.String_GetInt64(post["id"])
	status := util.String_GetInt64(post["status"])
	num, _ := jsondb.UpdateByID(c.jsonfile(), id, map[string]interface{}{"cstatus": status})
	if num > 0 {
		c.RenderJSON(map[string]interface{}{"msg": "操作成功", "url": "location"})
	} else {
		c.RenderJSON(map[string]interface{}{"msg": "操作失败", "errcode": -1})
	}
}
