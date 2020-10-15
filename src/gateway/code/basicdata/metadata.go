/**
网关运行需要的元数据
开发人员：陈朝能
20190903
*/
package basicdata

import (
	"util/deepcopy"
	"util"
	"strings"
	"sync"
)

//对元数据进行封装到map，便于快速读取
var _ApplicationMap map[int64]Application
var _ApplicationClientMap map[int64][]ApplicationClient
var _ModuleMap map[int64]Module
var _ModuleServerMapList map[int64][]ModuleServer
//var _ProtocolURLApplicationMapList map[int64][]Application

//网关元数据结构
type GatewayMetadata struct {
	sync.RWMutex
	//APIList           map[string]OpenAPI
	Application       []Application
	ApplicationClient []ApplicationClient
	Module            []Module
	ModuleServer      []map[string]interface{}
	ProtoURL          []ProtocolURL
	ProtoURLGET       []ProtocolURLGET
	UploadURL         []UploadURL
}

//网关元数据：运行网关需要的执行数据
var GWMetaList GatewayMetadata

func _InitData() {
	if Config.Timeout == 0 || Config.ErrContent == "" || Config.TimeoutContent == "" {
		util.LogError("global params is not config:requesttimeout、requesterrorcontent、requesttimeoutcontent")
	}

	//加载应用客户端
	_ApplicationClientMap = map[int64][]ApplicationClient{}
	for _, value := range DataDao.applicationclient {
		applicationID := util.String_GetInt64(value["capplicationid"])
		tmp := _ApplicationClientMap[applicationID]
		tmpApplication := Application{}

		//创建application
		for _, value2 := range DataDao.application {
			applicaitonPKID := util.String_GetInt64(value2["pkid"])
			if applicaitonPKID == applicationID {
				tmpApplication = Application{util.String_GetInt64(value2["pkid"]), value2["ctitle"].(string), nil,
					Status(util.String_GetInt64(value2["cstatus"])), value2["cpath"].(string)}
				break
			}
		}

		tmp = append(tmp, ApplicationClient{util.String_GetInt64(value["pkid"]), applicationID, tmpApplication, value["ctitle"].(string),
			value["cclient"].(string), Status(util.String_GetInt64(value["cstatus"])), strings.TrimSpace(value["cclientagent"].(string)), util.String_GetInt64(value["cenablesign"]),
			value["cencrypttype"].(string), value["cencryptcode"].(string), util.String_GetInt64(value["cenableaes"])})
		_ApplicationClientMap[applicationID] = tmp
	}

	//加载应用
	_ApplicationMap = map[int64]Application{}
	for _, value := range DataDao.application {
		applicaitonPKID := util.String_GetInt64(value["pkid"])
		_ApplicationMap[applicaitonPKID] = Application{util.String_GetInt64(value["pkid"]), value["ctitle"].(string), _ApplicationClientMap[applicaitonPKID],
			Status(util.String_GetInt64(value["cstatus"])), value["cpath"].(string)}
	}

	_ModuleMap = map[int64]Module{}
	for _, value := range DataDao.module {
		_ModuleMap[util.String_GetInt64(value["pkid"])] = Module{util.String_GetInt64(value["pkid"]), value["ctitle"].(string),
			Status(util.String_GetInt64(value["cstatus"]))}
	}

	_ModuleServerMapList = map[int64][]ModuleServer{}
	for _, value := range DataDao.moduleserver {
		tmp := _ModuleServerMapList[util.String_GetInt64(value["cmoduleid"])]
		tmp = append(tmp, ModuleServer{util.String_GetInt64(value["pkid"]), value["ctitle"].(string), Status(util.String_GetInt64(value["cstatus"])),
			util.String_GetInt64(value["cmoduleid"]), value["chost"].(string)})
		_ModuleServerMapList[util.String_GetInt64(value["cmoduleid"])] = tmp
	}
}

//初始化网关元数据
func InitGatewayMetadata() {
	DataDao.Lock()
	defer DataDao.Unlock()

	_InitData()

	GWMetaList.Lock()
	defer GWMetaList.Unlock()

	GWMetaList.ModuleServer = deepcopy.Copy(DataDao.moduleserver).([]map[string]interface{})

	tmpApplication := []Application{}
	for _, value := range _ApplicationMap {
		tmpApplication = append(tmpApplication, value)
	}
	GWMetaList.Application = tmpApplication

	tmpApplicationClient := []ApplicationClient{}
	for _, value := range _ApplicationClientMap {
		for _, val := range value {
			tmpApplicationClient = append(tmpApplicationClient, val)
		}
	}
	GWMetaList.ApplicationClient = tmpApplicationClient

	tmpModule := []Module{}
	for _, value := range _ModuleMap {
		tmpModule = append(tmpModule, value)
	}
	GWMetaList.Module = tmpModule

	tmplist := []ProtocolURL{}
	for _, value := range DataDao.protocolurl {
		tmptimeout := util.String_GetInt64(value["ctimeout"])
		if tmptimeout <= 0 {
			tmptimeout = Config.Timeout
		}

		ContainAllApp:=util.String_GetInt64(value["ccontainallapp"])== 1
		applications:= []Application{}

		//tmplist,_:=jsondb.Query("database/protocolurl.json", map[string]interface{}{"cstatus":1})
		for _,v1:=range _ApplicationMap {
			 if ContainAllApp{
			 	applications=append(applications,v1)
			 }else{
			 	tmpids:=strings.Split(util.String(value["capplicationids"]),",")
				 for _, v2 := range tmpids {
					 if v1.PKID==util.String_GetInt64(v2){
						 applications=append(applications,v1)
					 }
				 }
			 }
		}
		tmplist = append(tmplist, ProtocolURL{util.String_GetInt64(value["pkid"]), value["ctitle"].(string), Status(util.String_GetInt64(value["cstatus"])),
			value["csourceurl"].(string), value["ctargeturl"].(string), tmptimeout,
			ContainAllApp, applications,
			_ModuleMap[util.String_GetInt64(value["cmoduleid"])],
			_ModuleServerMapList[util.String_GetInt64(value["cmoduleid"])]})
	}
	GWMetaList.ProtoURL = tmplist

	tmplistget := []ProtocolURLGET{}
	for _, value := range DataDao.protocolurlget {
		tmplistget = append(tmplistget, ProtocolURLGET{util.String_GetInt64(value["pkid"]), value["ctitle"].(string), Status(util.String_GetInt64(value["cstatus"])),
			util.String_GetInt64(value["ccontainallapp"]) == 1, value["csourceurl"].(string), value["ctargeturl"].(string)})
	}
	GWMetaList.ProtoURLGET = tmplistget

	tmplist2 := []UploadURL{}
	for _, value := range DataDao.uploadlurl {
		tmplist2 = append(tmplist2, UploadURL{util.String_GetInt64(value["pkid"]), value["ctitle"].(string), Status(util.String_GetInt64(value["cstatus"])),
			util.String_GetInt64(value["cvalid"]) == 1, util.String_GetInt64(value["ccontainallapp"]) == 1, strings.ToLower(value["csourceurl"].(string)), value["ctargeturl"].(string), util.String_GetInt64(value["ctimeout"]), util.String_GetInt64(value["cmaxfilesize"]),
			value["cfileformat"].(string), util.String_GetInt64(value["cupnum"]), value["creturnerrorcontent"].(string), value["creturnsizecontent"].(string),
			value["creturnformatcontent"].(string), value["creturntimeoutcontent"].(string),
			_ModuleMap[util.String_GetInt64(value["cmoduleid"])],
			_ModuleServerMapList[util.String_GetInt64(value["cmoduleid"])]})
	}
	GWMetaList.UploadURL = tmplist2

}
