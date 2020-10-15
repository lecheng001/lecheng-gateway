package viewserver

import (
	"util/jsondb"
	"util"
)

type UploadController struct {
	Controller
}

func (c UploadController) jsonfile() string {
	return "database/upload.json"
}

func (c UploadController) GetList() {
	where := map[string]interface{}{"cstatus|>=": 0}
	data := map[string]interface{}{}
	data["list"], _ =jsondb.Query(c.jsonfile(),where, map[string]string{"csort":"asc"})
	modules, _ := jsondb.Query("database/module.json", map[string]interface{}{"cstatus|>=": 0}, map[string]string{"pkid":"asc"})
	moduleservers, _ := jsondb.Query("database/moduleserver.json", map[string]interface{}{"cstatus|>=": 0},map[string]string{"pkid":"desc"})
	for key, value := range modules {
		item := util.Map_getValues(moduleservers, "cmoduleid", value["pkid"].(string))
		if item != nil {
			modules[key]["chost"] = item["chost"]
		}
	}
	data["modules"] = modules
	c.RenderView("upload.gohtml", data)
}

func (c UploadController) GetInfo() {
	id := c.GetInt64("id", 0)
	if id == 0 {
		c.RenderError("", "")
		return
	}
	info:=jsondb.GetInfo(c.jsonfile(),id)
	c.RenderJSON(map[string]interface{}{"info": info})
}
func (c UploadController) PostData() {
	post := c.GetPostData()
	id:=util.String_GetInt64(post["pkid"])
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
func (c UploadController) DealStatus() {
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
