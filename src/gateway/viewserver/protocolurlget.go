package viewserver

import (
	"util/jsondb"
	"util"
)

type ProtocolUrlGetController struct {
	Controller
}

func (c ProtocolUrlGetController) jsonfile() string {
	return "database/protocolurlget.json"
}


func (c ProtocolUrlGetController) GetList() {
	where := map[string]interface{}{"cstatus|>=": 0}
	data := map[string]interface{}{}
	list,_:=jsondb.Query(c.jsonfile(),where, map[string]string{"csort":"asc"})
	data["list"] = list
	c.RenderView("protocolurlget.gohtml", data)
}

func (c ProtocolUrlGetController) GetInfo() {
	id := c.GetInt64("id", 0)
	info:=jsondb.GetInfo(c.jsonfile(),id)
	data := map[string]interface{}{"info": info}
	c.RenderJSON(data)
}
func (c ProtocolUrlGetController) PostData() {
	post := c.GetPostData()

	id:=util.String_GetInt64(post["pkid"])
	num:=int64(0)
	if id>0{
		num,_=jsondb.UpdateByID(c.jsonfile(),id,post)
	}else {
		num,_=jsondb.InsertMap(c.jsonfile(),post)
	}
	if num > 0 {
		c.RenderJSON(map[string]interface{}{"msg": "操作成功", "url": "/#/protocolurlget"})
	} else {
		c.RenderJSON(map[string]interface{}{"msg": "操作失败", "errcode": -1, "url": "location"})
	}
}
func (c ProtocolUrlGetController) DealStatus() {
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
