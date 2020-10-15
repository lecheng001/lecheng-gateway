package viewserver

import (
	"util/jsondb"
	"util"
	"strings"
)

type ModuleController struct {
	Controller
}

func (c ModuleController) jsonfile()string{
	return "database/module.json"
}

func (c ModuleController) Getlist() {
	list,_:=jsondb.Query(c.jsonfile(),map[string]interface{}{"cstatus|>=": 0}, map[string]string{"pkid":"desc"})
	c.RenderView("module.gohtml", map[string]interface{}{"list": list})
}

func (c ModuleController) Getinfo() {
	id := c.GetInt64("id", 0)
	if id == 0 {
		c.RenderError("", "")
		return
	}
	info:=jsondb.GetInfo(c.jsonfile(),id)
	c.RenderJSON(map[string]interface{}{"info": info})
}

func (c ModuleController) Postdata() {
	post := c.GetPostData()
	id:=util.String_GetInt64(post["pkid"])
	post["curl"] = strings.Trim(post["curl"].(string), "/")
	num:=int64(0)
	if id>0{
		num,_=jsondb.UpdateByID(c.jsonfile(),id,post)
	}else {
		num,_=jsondb.InsertMap(c.jsonfile(),post)
	}
	if num > 0 {
		c.RenderJSON(map[string]interface{}{"msg": "操作成功", "url": "location"})
	} else {
		c.RenderJSON(map[string]interface{}{"msg": "操作失败", "errcode": -1})
	}
}

func (c ModuleController) DealStatus() {
	post := c.GetPostData()
	id := util.String_GetInt64(post["id"])
	status := util.String_GetInt64(post["status"])
	num,_:=jsondb.UpdateByID(c.jsonfile(),id,map[string]interface{}{"cstatus": status})
	if num > 0 {
		c.RenderJSON(map[string]interface{}{"msg": "操作成功", "url": "location"})
	} else {
		c.RenderJSON(map[string]interface{}{"msg": "操作失败", "errcode": -1})
	}
}
