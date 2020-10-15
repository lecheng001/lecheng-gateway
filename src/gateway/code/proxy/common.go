package proxy

import (
	"util"
	"gateway/code/basicdata"
	"github.com/valyala/fasthttp"
	"time"
)

//请求HTTP
func getUrlContent(url, ip, lechengClient, lechengVersion, transmit string, method, userAgent, requestBody []byte, timeout int64) ([]byte, error) {
util.LogInfo("getUrlContent:",url, ip, lechengClient, lechengVersion, transmit , method, userAgent )
	req := &fasthttp.Request{}
	req.Header.Set("X-Real-Ip", ip)
	req.Header.Set("X-Forwarded-For", ip)
	//req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Lecheng-Client", lechengClient)
	req.Header.Set("Lecheng-Version", lechengVersion)
	req.Header.SetMethodBytes(method)
	req.Header.SetUserAgentBytes(userAgent)
	req.SetRequestURI(url)

	if transmit == "url" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else if transmit == "body_form" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else if transmit == "body_json" {
		req.Header.Set("Content-Type", "application/json")
	}
	req.SetBody(requestBody)
	resp := &fasthttp.Response{}
	err := fasthttp.DoTimeout(req, resp, time.Millisecond*time.Duration(timeout))
	resBody := resp.Body()

	//释放
	fasthttp.ReleaseResponse(resp)
	fasthttp.ReleaseRequest(req)
	resp = nil
	req = nil

	if err != nil {
		return nil, err
	} else {
		return resBody, nil
	}

}

//轮询服务器
func GetServerHost(moduleserver []basicdata.ModuleServer, moduleid int64) basicdata.ModuleServer {
	//tmplist := _ModuleServerMapList[moduleid]
	tmplen := len(moduleserver)
	if tmplen == 0 {
		util.LogError("服务器信息未配置")
		return basicdata.ModuleServer{}
	}
	tmpServer := basicdata.ModuleServer{}
	if tmplen >= 1 {
		tmpServer = moduleserver[0]
	}

	//记录使用信息
	util.LogDebug("使用服务器：", tmpServer.Title, tmpServer.PKID, tmpServer.Host)

	return tmpServer
}
