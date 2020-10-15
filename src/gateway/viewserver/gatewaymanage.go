package viewserver

import (
	"gateway/code/basicdata"
	"gateway/global"
	"runtime/debug"
	"strings"
	"sync"
	"util"
)

type GatewayManageController struct {
	Controller
}

func (c GatewayManageController) GetList() {
	debuginfo := util.ConfigGetBool("setting.debuginfo")
	c.RenderView("gatewaymanage.gohtml", map[string]interface{}{"debuginfo": debuginfo})
}

func (c GatewayManageController) Refresh() {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	refreshData(c.Controller)
	//c.RenderJSON(map[string]interface{}{"msg": "刷新成功"})
}

func (c GatewayManageController) RefreshConfig() {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	util.Config_RefreshConfig()
	c.RenderJSON(map[string]interface{}{"msg": "刷新成功"})
}

func (c GatewayManageController) ChangeDebug() {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	debuginfo := util.ConfigGetBool("setting.debuginfo")
	util.ConfigSet("setting.debuginfo", !debuginfo)
	c.RenderJSON(map[string]interface{}{"msg": "更改成功"})
}

func (c GatewayManageController) Certificate() {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	tmp, err := global.CheckSoftPackage()
	if tmp.Projectid > 0 {
		c.RenderJSON(map[string]interface{}{"msg": "刷新成功，项目名称：" + tmp.ProjectName + "，版本号：" + tmp.Version})
	} else {
		c.RenderError(err.Error(), "")
	}
}

//func (c GatewayManageController) LocalData() {
//	syncData(c.Controller)
//	//c.RenderJSON(map[string]interface{}{"msg": "刷新成功"})
//}

//刷新数据
func refreshData(c Controller) {
	defer func() {
		if err := recover(); err != nil {
			util.LogError("刷新失败", err, string(debug.Stack()))
		}
	}()

	tmpbody := strings.Replace(string(c.Request.Body()), "\"", "", -1)
	if "choice147852" != tmpbody {
		c.RenderError("基础数据刷新失败", "")
		return
	}

	basicdata.InitGatewayConfig()
	//初始化基础数据
	basicdata.LoadDataSource()
	//util.LogDebug(basicdata.DataDao)
	//解析所有节点
	basicdata.InitGatewayMetadata()
	tmp := map[string]interface{}{"error": 0, "msg": "基础数据已经刷新"}
	c.RenderJSON(tmp)
}
