package viewserver

import (
	"util"
	"util/jsondb"
)

type ApplicationClientController struct {
	Controller
}
func (c ApplicationClientController) jsonfile()string{
	return "database/applicationclient.json"
}

func (c ApplicationClientController) Getlist() {
	appid := c.GetInt64("appid")
	list,_:=jsondb.Query(c.jsonfile(), map[string]interface{}{"cstatus|>=": 0, "capplicationid": util.String(appid)},nil)
	info:= jsondb.GetInfo("database/application.json",appid)
	c.RenderView("applicationclient.gohtml", map[string]interface{}{"list": list, "info": info})
}

func (c ApplicationClientController) Getinfo() {
	id := c.GetInt64("id", 0)
	if id == 0 {
		c.RenderError("数据丢失", "/html/application/list")
		return
	}
	info:= jsondb.GetInfo(c.jsonfile(),id)
	c.RenderJSON(map[string]interface{}{"info": info})
}

func (c ApplicationClientController) Postdata() {
	post := c.GetPostData()
	id:=util.String_GetInt64(post["pkid"])

	//判断重复
	tmplist,_:= jsondb.Query(c.jsonfile(), map[string]interface{}{"cclientagent": post["cclientagent"]},nil)
	if len(tmplist)>0 && util.String_GetInt64( tmplist[0]["pkid"])!=id{
		c.RenderJSON(map[string]interface{}{"msg": "操作失败,lecheng-client不能重复", "errcode": -1})
		return
	}

	num:=int64(0)
	if id>0{
		num,_=jsondb.UpdateByID(c.jsonfile(),id,post)
	}else{
		num,_=jsondb.InsertMap(c.jsonfile(),post)
	}
	if num > 0 {
		  c.RenderJSON(map[string]interface{}{"msg": "操作成功", "url": "location"})
	} else {
		  c.RenderJSON(map[string]interface{}{"msg": "操作失败", "errcode": -1})
	}
}

func (c ApplicationClientController) DealStatus() {
	post := c.GetPostData()
	id := util.String_GetInt64(post["id"])
	status := util.String_GetInt64(post["status"])
	num, _ := jsondb.UpdateByID("database/applicationclient.json",id, map[string]interface{}{"cstatus": status})
	if num > 0 {
		  c.RenderJSON(map[string]interface{}{"msg": "操作成功", "url": "location"})
	} else {
		  c.RenderJSON(map[string]interface{}{"msg": "操作失败", "errcode": -1})
	}
}
