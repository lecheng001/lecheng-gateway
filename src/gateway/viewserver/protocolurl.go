package viewserver

import (
	"util/jsondb"
	"util"
	"strings"
)

type ProtocolUrlController struct {
	Controller
}

func (c ProtocolUrlController) jsonfile() string {
	return "database/protocolurl.json"
}

func (c ProtocolUrlController) GetList() {
	where := map[string]interface{}{"cstatus|>=": 0}
	searchkey := c.GetString("searchkey")
	appid := c.GetInt64("appid")
	if appid > 0 {
		where["oap.capplicationid"] = appid
	}
	if len(searchkey) > 0 {
		where["oa.ctitle|like"] = searchkey
	}

	data := map[string]interface{}{}
	urllist, _ := jsondb.Query(c.jsonfile(), where, map[string]string{"csort":"asc","cstatus":"desc"})
	applications, _ := jsondb.Query("database/application.json", map[string]interface{}{"cstatus|>=": 0}, map[string]string{"pkid":"asc"})
	modules, _ := jsondb.Query("database/module.json", map[string]interface{}{"cstatus|>=": 0}, map[string]string{"pkid":"asc"})
	moduleservers, _ := jsondb.Query("database/moduleserver.json", map[string]interface{}{"cstatus|>=": 0}, map[string]string{"pkid":"desc"})
	for key, value := range modules {
		item := util.Map_getValues(moduleservers, "cmoduleid", value["pkid"].(string))
		if item != nil {
			modules[key]["chost"] = item["chost"]
		}
	}

	data["modules"] = modules
	//data["moduleservers"]=moduleservers
	data["applications"] = applications
	for key, value := range urllist {
		//clients := ""
		if util.String_GetInt64(value["ccontainallapp"]) == 1 {
			urllist[key]["applications"] = "全部项目"
		} else {
			tmpApp := ""
			appids := strings.Split(util.String(value["capplicationids"]), ",")
			for _, v := range appids {
				for _, v2 := range applications {

					if util.String_GetInt64(v) == util.String_GetInt64(v2["pkid"]) {
						tmpApp += "," + v2["ctitle"].(string)
					}
				}
			}
			if len(tmpApp) > 0 {
				tmpApp = tmpApp[1:]
			}
			urllist[key]["applications"] = tmpApp

			//list[key]["applications"] = strings.Trim(clients, ",")
		}
	}
	data["list"] = urllist
	c.RenderView("protocolurl.gohtml", data)
}

func (c ProtocolUrlController) GetInfo() {
	id := c.GetInt64("id", 0)
	//if id == 0 {
	//	return c.RenderErr("", "")
	//}
	info := jsondb.GetInfo(c.jsonfile(), id)
	data := map[string]interface{}{"info": info}
	c.RenderJSON(data)
}
func (c ProtocolUrlController) PostData() {
	post := c.GetPostData()
	id := util.String_GetInt64(post["pkid"])
	//applicationids := post["applicationids"]
	//delete(post, "applicationids")

	//tmp := c.checkUrl(post["csourceurl"].(string), util.String_GetInt64(post["pkid"]))
	//if tmp != nil {
	//	c.RenderJSON(map[string]interface{}{"msg": "地址重复，无法提交", "errcode": -1})
	//	return
	//}

	num := int64(0)
	if id > 0 {
		num, _ = jsondb.UpdateByID(c.jsonfile(), id, post)
	} else {
		num, _ = jsondb.InsertMap(c.jsonfile(), post)
	}
	if num > 0 {
		c.RenderJSON(map[string]interface{}{"msg": "操作成功", "url": "/#/protocolurl"})
	} else {
		c.RenderJSON(map[string]interface{}{"msg": "操作失败", "errcode": -1, "url": "location"})
	}
}
func (c ProtocolUrlController) DealStatus() {
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

