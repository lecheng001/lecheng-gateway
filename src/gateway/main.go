/**
开发人员：陈朝能
网关主入口
*/
package main

import (
	"util"
	"encoding/json"
	"fmt"
	"gateway/code/basicdata"
	"gateway/code/proxy"
	"gateway/global"
	"gateway/viewserver"
	"github.com/valyala/fasthttp"
	"runtime"
	"runtime/debug"
	"strings"
)
import _ "net/http/pprof"


func main() {
	//监测可用性
	//go main2()

	util.LogError("启动proxy")
	global.CheckSoftPackage()
	//设置CPU个数,有证书CPU全开，否则只启用一个CPU
	if proxy.GWPACKAGE.Userid > 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(1)
	}

	//初始化基础数据
	basicdata.LoadDataSource()

	//解析所有节点
	basicdata.InitGatewayMetadata()

	//获取网关签名名称
	proxy.SIGNNAME = util.ConfigGetString("setting.signname")
	if proxy.SIGNNAME == "" {
		panic("CONFIG配置文件，setting.signname 未配置数据")
	}
	proxy.ENCRYPTPARAMNAME = util.ConfigGetString("setting.encryptparamnname")
	if proxy.ENCRYPTPARAMNAME == "" {
		panic("CONFIG配置文件，setting.encryptparamnname 未配置数据")
	}
	proxy.JSONPARAMNAME = util.ConfigGetString("setting.jsonparamnname")
	if proxy.JSONPARAMNAME == "" {
		panic("CONFIG配置文件，setting.jsonparamnname 未配置数据")
	}
	proxy.GATEWAYAPIPATH = util.ConfigGetString("setting.gatewayapipath")
	if proxy.GATEWAYAPIPATH == "" {
		panic("CONFIG配置文件，setting.gatewayapipath 未配置数据")
	}

	//获取网关透传参数（仍需要参与验签，不受业务参数限制）
	proxy.PASSTHROUGHPARA = util.ConfigGetSlice("setting.passthroughpara")

	//util.LogDebug("GWMetaList:", basicdata.GWMetaList)
	hostport := util.ConfigGetString("hostport")
	util.LogError("端口号：" + hostport)
	//fasthttp.ListenAndServe(hostport, ServeFastHTTP)

	m := func(ctx *fasthttp.RequestCtx) {
		lvPath := string(ctx.Path())
		fmt.Println(lvPath)
		if lvPath=="/"{
			ctx.Response.SetBody([]byte(`{"errcode":-1,"msg":"请不要非法请求"}`))
			return
		}

		if lvPath=="/version"{
			ctx.Response.SetBody([]byte(`网关版本号：`+global.GW_VERISON))
			return
		}

		if strings.Index(lvPath, "/html") == 0 {
			//载入静态资源
			if strings.Index(lvPath, "/html/res/") == 0 {
				tmpPath := strings.Replace(lvPath, "/html/res/", "/static/", 1)
				fasthttp.ServeFile(ctx, util.IO_GetRootPath()+tmpPath)
				return
			}
			ctx.Response.Header.Set("Content-Type", "text/html ")
			//viewserver.InitRouter(ctx)
			HTMLFastHTTP(ctx)
			return
		}  else {
			//没有证书，40%请求失败
			if proxy.GWPACKAGE.Userid <= 0 {
				lvRand := util.Rand_Int64(1, 10)
				if lvRand < 5 {
					proxy.Render{}.RendJson(ctx, "未配置证书，请从官网下载网关证书")
					return
				}
			}

			ServeFastHTTP(ctx)
			return
		}
	}
	fasthttp.ListenAndServe(hostport, m)

	////阻止线程关闭
	//select {}
}

//监听处理主入口
func ServeFastHTTP(ctx *fasthttp.RequestCtx) {
	fmt.Println("ServeFastHTTP")
	util.LogDebug("========================start fasthttp ======================")
	defer func() {
		if err := recover(); err != nil {
			util.LogError("has a no know error：", err, string(debug.Stack()))

			tmp, err := json.Marshal(map[string]interface{}{"errcode": -9999, "msg": fmt.Sprint(err), "content": nil})
			if err != nil {
				fmt.Fprint(ctx, err)
			}
			ctx.Error(string(tmp), 200)
		}
	}()

	//处理跨域请求
	cros := util.ALLCONFIG.Get("cros").MustMap()
	if cros != nil {
		for key, value := range cros {
			ctx.Response.Header.Set(key, value.(string))
		}
	}

	if string(ctx.Request.Header.Method()) == "OPTIONS" {
		ctx.Response.SetBody([]byte("ok"))
		ctx.Response.SetStatusCode(200)
		ctx.Response.SetConnectionClose()

		return
	}

	//fmt.Println(string( ctx.Request.Body()))
	//取得全部参数
	rData, err := proxy.RequestData{}.GetInitRequestData(ctx)
	if err != nil {
		proxy.Render{}.RendErr(ctx, err.Error())
		return
	}
	//处理网关系统接口
	if "system" == rData.UriScheme[0] {
		util.LogDebug("to run system.Run")
		//system.Run(ctx)
		return
	}
	//处理网关系统接口
	if "html" == rData.UriScheme[0] {
		util.LogDebug("to run system.Run")
		//system.Run(ctx)
		return
	}

	//载入网关运行需要的数据，以指针形式传递并失败，避免读写异常。
	rData.Metadata = basicdata.GetGatewayMetaData()
	rData.Config = basicdata.GetGatewaySettingConfig()

	//业务判断，GET 302转发，302转跳没有头部信息，无需校验一致性
	if rData.IsGet && !rData.IsRESTfull {
		r := proxy.Reverse{}.NewRevers(rData)
		r.ReverseProxyGet()
		return
	}

	//判断入口地址与头部信息一致性
	err = rData.CheckApplication()
	if err != nil {
		proxy.Render{}.RendErr(ctx, err.Error())
		return
	}

	//处理RESTFull请求数据，分析出请求参数（未加密和验签的请求）
	err = rData.AnsyRequestData()
	if err != nil {
		proxy.Render{}.RendErr(ctx, err.Error())
		return
	}

	//执行网关业务
	valid := proxy.Valid{}.NewValid(rData)
	blnNeedValid := true
	u := proxy.Upload{}.NewUpload(rData)
	//处理图片问题
	if rData.UpFile_HasFile {
		err := u.GetUploadRule()
		if err != nil {
			proxy.Render{}.RendErr(ctx, err.Error())
			return
		}
		blnNeedValid = rData.UpFile_UploadInfo.Valid
	}

	//不需要校验数据，现阶段只有图片上传
	if blnNeedValid {
		//解析请求参数，解密失败将推出
		//解密去验签，得到明文的POST数据,请求数据放在paramsMap
		err = valid.CheckValid()
		if err != nil {
			errMsg := err.Error()
			util.LogDebug("非法请求：", rData.Uri, errMsg, rData.Header)

			proxy.Render{}.RendErr(ctx, errMsg)
			return
		}
	}

	//需要安全性校验，将参数一并传到后端服务器
	if rData.UpFile_HasFile {
		util.LogDebug("进入图片上传")
		err := u.UploadFile()
		if err != nil {
			//日志
			tmp := err.Error()
			rData.Response = tmp
			proxy.Render{}.RendErr(ctx, tmp)
		} else {
			proxy.Render{}.RendByte(ctx, rData.ResponseByte)
		}

		return
	}

	//处理RESTfull请求，支持 GET,POST,PUT,DELETE
	//执行正常网关逻辑
	p := proxy.Proxy{}.NewProxy(rData, valid)
	p.RunProxy()

}

func HTMLFastHTTP(ctx *fasthttp.RequestCtx) {
	defer func() {
		if err := recover(); err != nil {
			util.LogError("has a no know error：", err, string(debug.Stack()))

			//tmp, err := json.Marshal(map[string]interface{}{"errcode": -9999, "msg": fmt.Sprint(err), "content": nil})
			//if err != nil {
			//	fmt.Fprint(ctx, err)
			//}
			//ctx.Error(string(tmp), 200)
			//
			ctx.Response.Header.Set("Content-Type", "text/html ")
			cc := viewserver.Controller{ctx}
			cc.RenderView("error.gohtml", map[string]interface{}{"msg":err})
		}
	}()

	fmt.Println("HTMLFastHTTP")
	fmt.Println(string(ctx.RequestURI()))
	fmt.Println(string(ctx.PostBody()))
	fmt.Println(string(ctx.QueryArgs().QueryString()))
	ctx.Response.Header.Set("Content-Type", "text/html ")
	//RenderView(ctx,"index.gohtml",nil)

	viewserver.InitRouter(ctx)
}
