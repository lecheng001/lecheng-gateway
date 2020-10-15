/**
网关配置数据
开发人员：陈朝能
20190903
*/
package basicdata

import (
	"util"
	"sync"
)

type Syssettingconfig struct {
	sync.RWMutex
	//IP黑名单
	IpLimitArr []string
	//IP白名单
	IpSafeArr            []string
	CircuteReturnContent string
	Timeout              int64
	ErrContent           string
	TimeoutContent       string
	NotFindContent       string
	ErrPostContent       string
	NoPostContent        string
	 //RequestLongTime      int64
	//EnableModuleAPILog   int64
	//EnableConnNum        int64
	//EnableOpenAPILog     int64
	//RunFastMode          int64
	//监控方式:monitortype  监控周期seconds  请求次数times  处理方式dealtype[temporary,limitfile]   禁止访问时间limitminute
	LimitIPRule []map[string]interface{}

	//BlnConnnum                bool
	//BlnLogmoduleapiresult     bool
	//BlnLogmoduleapi           bool
	//BlnLogopenapi             bool
	//BlnLogopenapiresult       bool
	//BlnGatewdataFromLocalfile bool
	//BlnLimitIP bool

}

var Config Syssettingconfig

func init()  {
	InitGatewayConfig()
}

//初始化网关配置数据
func InitGatewayConfig() {
	Config.Lock()
	defer Config.Unlock()

	js, err := util.JSON_Config("database/data_setting.json")
	if err != nil {
		panic("database/data_setting.json 配置信息丢失！")
		return
	}

	list := js.MustMap()

	for key, _ := range list {
		switch key {
		case "requesttimeout":
			Config.Timeout = util.String_GetInt64(js.GetPath(key, "value").MustString())
		case "txt_requesterrorcontent":
			Config.ErrContent = js.GetPath(key, "value").MustString()
		case "txt_requesttimeoutcontent":
			Config.TimeoutContent = js.GetPath(key, "value").MustString()
		case "txt_request404content":
			Config.NotFindContent = js.GetPath(key, "value").MustString()
		case "txt_nopostbody":
			Config.NoPostContent = js.GetPath(key, "value").MustString()
		case "txt_errpostbody":
			Config.ErrPostContent = js.GetPath(key, "value").MustString()
		}
	}
}
