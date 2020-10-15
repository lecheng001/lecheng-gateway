/**
普通模式
开发人员：陈朝能
*/
package proxy

import (
	"util"
	"encoding/json"
	"errors"
	"fmt"
	"gateway/code/basicdata"
	"github.com/valyala/fasthttp"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var SIGNNAME = ""
var ENCRYPTPARAMNAME = ""
var JSONPARAMNAME = ""
var PASSTHROUGHPARA = []string{}
var GATEWAYAPIPATH = ""
var GWPACKAGE GwPackage

type Proxy struct {
	rData     *RequestData
	timestamp int64
	metadata  *basicdata.GatewayMetadata
	config    *basicdata.Syssettingconfig
	ctx       *fasthttp.RequestCtx
	valid     *Valid

	openAPI basicdata.OpenAPI
	nodes   []dispatchNode
	responseContent []byte

	starttime   time.Time
	endtime     time.Time
	elapsedtime time.Duration
	//gwpackage GwPackage
}

type GwPackage struct {
	Version string
	//domain string
	Projecttype string
	Projectid int64
	ProjectName string
	Userid int64
	//validtime int64
}

//调度节点
type dispatchNode struct {
	tag   string
	title string
	applicationid int64
	moduleid      int64
	moduleapiid   int64
	host          string
	serverid      int64
	requesturi    string
	url           string
	timeout              int64
	isDirectReturn       bool
	returnErrorContent   string
	returnTimeoutContent string
	requestBody          []byte

	params map[string]interface{}

	isGet          bool
	isPost         bool
	isPut          bool
	isDelete       bool
	isBodyTransmit bool
	isBodyPayload  bool
	isBodyKeyValue bool

	resbody []byte
	err     error
	code    int64
}

//节点处理结果
type responseStruct struct {
	//节点索引
	index int
	//节点tag
	tag string
	//响应内容
	resBody string
	//错误信息
	err error
	//HTTP状态
	code int64
	//开始时间
	starttime time.Time
	//结束时间
	endtime time.Time
	//耗时
	elapsedtime time.Duration
	////已经写入日志
	//hasLog bool
}

//var M runtime.MemStats
//初始化数据
func (p Proxy) NewProxy(r *RequestData, valid *Valid) *Proxy {
	pro := &Proxy{}
	pro.timestamp = time.Now().Unix()
	pro.ctx = r.Ctx
	pro.config = r.Config
	pro.metadata = r.Metadata
	pro.valid = valid
	pro.rData = r
	return pro
}

func (p *Proxy) RunProxy() {
	util.LogInfo("进入正常网关模式")
	//合法性检测
	uriPath := p.rData.Uri
	util.LogInfo("RunProxy", p.rData.Uri, p.rData.UriScheme, p.rData.UriQuery, p.rData.UriFull, p.rData.UriPath)

	NewOpenapiInfo := basicdata.OpenAPI{PKID: -1}

	//遍历协议适配
	for _, value := range p.metadata.ProtoURL {
		preg := "^" + value.SourceURL + "$"
		reg, _ := regexp.Compile(preg)
		tmplist := reg.FindStringSubmatch(uriPath)
		if len(tmplist) == 0 {
			continue
		}
		//fmt.Println(tmplist)
		//fmt.Println(value)
		//取得转发的具体地址
		util.LogInfo(value.TargetURL)
		targeturl := value.TargetURL
		for k, val := range tmplist {
			if k == 0 {
				continue
			}

			targeturl = strings.Replace(targeturl, "$"+strconv.Itoa(k), val, -1)
		}
		util.LogInfo(targeturl)
		//封装NewOpenapiInfo
		NewOpenapiInfo.PKID = 0
		NewOpenapiInfo.Title = "协议适配"
		//NewOpenapiInfo.Applications = value.Applications
		NewOpenapiInfo.Method = p.rData.Method
		NewOpenapiInfo.Status = value.Status
		NewOpenapiInfo.Results = nil
		NewOpenapiInfo.Timeout = value.Timeout
		util.LogInfo("RunProxy", uriPath)
		NewOpenapiInfo.URL = uriPath
		NewOpenapiInfo.Nodes = map[string][]basicdata.OpenAPINode{}

		util.LogInfo("RunProxy", targeturl)
		tmpModuleAPI := basicdata.ModuleAPI{0, "协议适配微服务接口", basicdata.Status(1), 0, targeturl, NewOpenapiInfo.Method, "",
			"", "", "", value.Timeout}

		datatranslate := ""
		if p.rData.IsGet {
			datatranslate = "url"
		} else if p.rData.IsBodyJson {
			datatranslate = "body_json"
		} else if p.rData.IsBodyTransmit {
			datatranslate = "body_form"
		}
		tmpOpenAPINode := []basicdata.OpenAPINode{{0, "协议适配节点", "", "", p.rData.Method, datatranslate, "", false, p.config.ErrContent,
			p.config.TimeoutContent, value.Timeout, 0, p.rData.ApplicationClient.ApplicationId, value.Module, tmpModuleAPI, value.Servers}}

		NewOpenapiInfo.Nodes[p.rData.Method] = tmpOpenAPINode
		//fmt.Println(NewOpenapiInfo)o
		break
	}
	//}

	//找不到配置的URL
	if NewOpenapiInfo.PKID < 0 {
		util.LogInfo(fmt.Sprintf("地址：%s不存在，请不要非法请求", p.rData.Uri))
		Render{}.RendErr(p.ctx, fmt.Sprintf("地址：%s不存在，请不要非法请求", p.rData.Uri))
		return
	}

	if len(NewOpenapiInfo.Nodes) == 0 {
		Render{}.RendErr(p.ctx, "没有需要请求的节点")
		return
	}

	if _, ok := NewOpenapiInfo.Nodes[p.rData.Method]; !ok {
		Render{}.RendErr(p.ctx, "没有需要请求的节点")
		return
	}
	if len(NewOpenapiInfo.Nodes[p.rData.Method]) == 0 {
		Render{}.RendErr(p.ctx, "没有需要请求的节点")
		return
	}

	//判断服务器情况
	for _, value := range NewOpenapiInfo.Nodes[p.rData.Method] {
		if len(value.Servers) == 0 {
			Render{}.RendErr(p.ctx, value.Module.Title+"没有可用的服务器")
			return
		}
	}

	p.openAPI = NewOpenapiInfo
	err := p.initProxyData()
	if err != nil {
		Render{}.RendErr(p.ctx, err.Error())
		return
	}

	p.start()
	p = nil

}

//实例proxy
func (p *Proxy) initProxyData() error {
	//util.LogDebug("OpenAPI数据：", openapi)
	util.LogInfo("initProxyData", p.openAPI.Nodes)
	dnodes := []dispatchNode{}
	for _, value := range p.openAPI.Nodes[p.rData.Method] {
		dnode := dispatchNode{}
		dnode.title = value.Title
		//中断
		//isbreak := false

		//分配指定服务器
		modulelserver := GetServerHost(value.Servers, value.Module.PKID)
		dnode.applicationid = value.ApplicationID
		dnode.moduleid = value.Module.PKID
		dnode.serverid = modulelserver.PKID
		dnode.moduleapiid = value.Moduleapi.PKID
		dnode.tag = value.Tag
		dnode.host = modulelserver.Host
		dnode.requesturi = string(p.ctx.Request.RequestURI())
		//测试环境，手动配置转跳地址
		if len(p.rData.Test_host) > 0 {
			dnode.url = p.rData.Test_host + value.Moduleapi.URL
		} else {
			dnode.url = "http://" + dnode.host + value.Moduleapi.URL
		}

		//util.LogInfo("initProxyData", dnode.host, value.Moduleapi.URL)
		//util.LogInfo("initProxyData", dnode)
		dnode.isDirectReturn = value.IsMustReturn
		dnode.timeout = value.Moduleapi.Timeout
		//openapi.timeout比重比moduleapi.timeout高
		if p.openAPI.Timeout < dnode.timeout || dnode.timeout == 0 {
			dnode.timeout = p.openAPI.Timeout
		}

		dnode.returnErrorContent = value.ReturnErrorContent
		dnode.returnTimeoutContent = value.ReturnTimeoutContent

		dnodes = append(dnodes, dnode)
	}
	p.nodes = dnodes
	p.starttime = time.Now()

	return nil

}

//监听处理主入口
func (p *Proxy) start() {
	//logger, _ := zap.NewProduction()
	//暂时注释，等待测试结果
	//defer func() {
	//	//释放
	//	p.dispatch = dispatch{}
	//
	//	if err := recover(); err != nil {
	//		Render{}.RendErr(p.ctx, fmt.Sprint(err))
	//		util.LogError("has a error：", err, string(debug.Stack()))
	//	}
	//}()

	//defer logger.Sync()

	time1 := time.Now()

	//5：fasthttp请求，结果处理，异常处理 记录日志
	//用waitgropu处理
	lennodes := len(p.nodes)
	util.LogInfo("需要执行URL总数：", lennodes)
	ch := make(chan responseStruct, lennodes)
	var wg = sync.WaitGroup{}
	wg.Add(lennodes)
	for key, _ := range p.nodes {
		go p.getContentWG(key, ch, &wg)
	}
	wg.Wait()
	util.LogInfo("HTTP请求完成时间：", time.Since(time1))
	p.endtime = time.Now()
	p.elapsedtime = time.Since(p.starttime)
	tmpResponse := []responseStruct{}
	for i := 0; i < lennodes; i++ {
		tmp := <-ch
		tmpResponse = append(tmpResponse, tmp)
	}
	close(ch)
	//fmt.Println("请求响应结果：", tmpResponse)

	//取得最终渲染结果
	p.RenderResponse(tmpResponse)

	tmpResponse = nil
	util.LogInfo("HTTP结果处理时间：", time.Since(time1))
	util.LogDebug("=============deal proxy  stop =========================")
}

//处理请求
func (p *Proxy) getContentWG(index int, ch chan responseStruct, wg *sync.WaitGroup) {
	defer wg.Done()
	node := p.nodes[index]
	err := errors.New("")
	err = nil
	requestBody := []byte{}
	resBody := ""
	url := node.url
	resResult := responseStruct{}
	resResult.index = index
	resResult.tag = node.tag
	//moduleAPILogtype := ""

	util.LogInfo("getContentWG", url)

	//强制返回
	if strings.TrimSpace(url) == "" || node.isDirectReturn {
		resBody = node.returnErrorContent
		err = errors.New("direct return")
		resResult.err = err
	} else {
		resResult.starttime = time.Now()
		util.LogInfo("执行URL_", index+1, "：", url)
		//url := "http://localhost:7090/aa.go"
		//方法一,不可采用，出现no such host  错误
		// req := fasthttp.AcquireRequest()
		//

		//与极速模式请求统一方法
		transmit := ""
		if p.rData.IsBodyJson {
			transmit = "body_json"
		} else if p.rData.IsBodyForm {
			transmit = "body_form"
		} else if !p.rData.IsBodyTransmit {
			transmit = "url"
		} else {
			util.LogError("代码逻辑缺漏")
			resBody = node.returnErrorContent
			err = errors.New("代码逻辑缺漏")
			resResult.err = err
		}

		if err == nil {
			//tmprequstBody := []byte{}
			if p.rData.IsBodyJson {
				requestBody, _ = json.Marshal(node.params)
			} else if p.rData.IsBodyForm {
				tmpStr := ""
				for key, value := range node.params {
					tmpStr += fmt.Sprintf("&%s=%s", key, value)
				}
				requestBody = []byte(strings.Trim(tmpStr, "&"))
			}
			tmpresBody := []byte{}

			tmpresBody, err = getUrlContent(url, p.rData.IP, p.rData.HeaderClient, util.String(p.rData.HeaderVersion), transmit, []byte(p.rData.Method), p.rData.userAgent, requestBody, node.timeout)
			resBody = string(tmpresBody)
			resResult.err = err
		}

	}

	resResult.endtime = time.Now()
	resResult.elapsedtime = time.Since(resResult.starttime)
	//============================这边调整一下，记录返回错误日志
	if err != nil {
		util.LogError("请求失败：", err, ",url:", url)
		//超时
		if err == fasthttp.ErrTimeout {
			resResult.resBody = node.returnTimeoutContent
			resResult.code = fasthttp.StatusRequestTimeout
			//moduleAPILogtype = "timeout"
		} else if err == fasthttp.ErrNoFreeConns {
			util.LogError("处理：fasthttp出现ErrNoFreeConns问题")
			resResult.resBody = "fasthttp出现ErrNoFreeConns问题"
			resResult.code = fasthttp.StatusTooManyRequests
			//moduleAPILogtype = "nofreeconns"
		} else {
			resResult.resBody = node.returnErrorContent
			resResult.code = fasthttp.StatusInternalServerError
			//moduleAPILogtype = "servererror"
		}
	} else {
		resResult.code = fasthttp.StatusOK
		resResult.resBody = resBody
	}

	if util.ConfigGetBool("setting.debuginfo") {
		tmpJSON := map[string]interface{}{}
		erra := json.Unmarshal([]byte(resResult.resBody), &tmpJSON)
		if erra == nil {
			tmpJSON["debugInfo"] = map[string]interface{}{"url": url, "requestBody": string(requestBody), "ip": p.rData.IP, "lecheng-client": p.rData.HeaderClient, "lecheng-version": util.String(p.rData.HeaderVersion)}
			tmpa, _ := json.Marshal(tmpJSON)
			resResult.resBody = string(tmpa)
		} else {
			tmpJSON["debugInfo"] = map[string]interface{}{"responseContent": resResult.resBody, "url": url, "requestBody": string(requestBody), "ip": p.rData.IP, "lecheng-client": p.rData.HeaderClient, "lecheng-version": util.String(p.rData.HeaderVersion)}
			tmpa, _ := json.Marshal(tmpJSON)
			resResult.resBody = string(tmpa)
		}
	}

	//返回
	ch <- resResult
	//util.LogDebug("end:", string(tmp.resbody))
}

//渲染结果
func (p *Proxy) RenderResponse(res []responseStruct) {
	//fmt.Println(p.dispatch.openAPI.Results)
	//只有一个节点，且没有配置转换，直接返回结果
	if len(res) == 1 && len(p.openAPI.Results) == 0 {
		p.ctx.SetStatusCode(fasthttp.StatusOK)
		contenta := res[0].resBody
		p.responseContent = []byte(contenta)
		p.ctx.Response.SetBodyString(contenta)
		return
	}

	p.ctx.Response.SetBodyString("数据处理失败，只支持一个节点")
	return

}