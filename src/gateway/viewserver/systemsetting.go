package viewserver

import (
	"util"
	"gateway/code/basicdata"
)

type SystemSettingController struct {
	Controller
}

func (c SystemSettingController) GetList() {
	setting, err := util.JSON_Config("database/data_setting.json") // util.ALLCONFIG.Get("syssetting")
	if setting.Interface() == nil || err != nil {
		util.LogError("数据data_setting.json读取失败")
	}
	tmp := map[string]interface{}{"error": 0, "data": setting.Interface()}
	c.RenderView("systemsetting.gohtml", tmp)
}

func (c SystemSettingController) ConfigWrite() {
	//fmt.Println(ctx.PostArgs())
	postdata := c.GetPostDataString()

	key := postdata["key"]
	if key == "" {
		c.RenderError("key不能为空", "")
		return
	}
	value := postdata["value"]

	if !util.IO_FilePathExists(util.IO_GetRootPath() + "/database/data_setting.json") {
		panic("database/data_setting.json 配置文件丢失")
		//setting := basicdata.DefaultSettingConfig()
		//tmpContent, _ := setting.MarshalJSON()
	}
	err := util.JSON_ConfigSet("database/data_setting.json", key+".value", value)

	if err != nil {
		c.RenderError(err.Error(), "")
		return
	} else {
	}
	basicdata.InitGatewayConfig()

	c.RenderJSON(map[string]interface{}{"msg": "修改成功！"})
}

