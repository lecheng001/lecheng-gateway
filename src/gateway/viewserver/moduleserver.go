package viewserver

import (
	"util/jsondb"
	"util"
	"strings"
)

type ModuleserverController struct {
	Controller
}

func (c ModuleserverController) jsonfile() string {
	return "database/moduleserver.json"
}

func (c ModuleserverController) GetList() {
	where := map[string]interface{}{"ms.cstatus|>=": 0}
	searchkey := c.GetString("searchkey")
	//moduleid := c.GetInt64("moduleid")

	if len(searchkey) > 0 {
		where["ma.ctitle|like"] = searchkey
	}

	data := map[string]interface{}{}
	modulelist, _ := jsondb.Query("database/module.json", map[string]interface{}{"cstatus|>=": 0}, map[string]string{"pkid":"asc"})
	list, _ := jsondb.Query(c.jsonfile(),  map[string]interface{}{"cstatus|>=": 0}, map[string]string{"pkid":"desc"})
	for key, value := range list {
		item := util.Map_getValues(modulelist, "pkid", value["cmoduleid"].(string))
		if item != nil {
			list[key]["moduletitle"] = item["ctitle"]
		}
	}
	data["list"] = list
	data["modules"] = modulelist
	c.RenderView("moduleserver.gohtml", data)
}



func (c ModuleserverController) GetInfo() {
	id := c.GetInt64("id", 0)
	if id == 0 {
		c.RenderError("", "")
		return
	}
	info:=jsondb.GetInfo(c.jsonfile(),id)
	c.RenderJSON(map[string]interface{}{"info": info})
}
func (c ModuleserverController) PostData() {
	post := c.GetPostData()
	id:=util.String_GetInt64(post["pkid"])

	post["chost"] = strings.Replace(util.String(post["chost"]), "：", ":", -1)
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
func (c ModuleserverController) DealStatus() {
	post := c.GetPostData()
	id := util.String_GetInt64(post["id"])
	status := util.String_GetInt64(post["status"])
	num,_:=jsondb.UpdateByID(c.jsonfile(),id, map[string]interface{}{"cstatus": status})
	if num > 0 {
		c.RenderJSON(map[string]interface{}{"msg": "操作成功", "url": "location"})
	} else {
		c.RenderJSON(map[string]interface{}{"msg": "操作失败", "errcode": -1})
	}
}
