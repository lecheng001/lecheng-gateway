/**
极速模式
开发人员：陈朝能
*/
package proxy

import (
	"util"
	"gateway/code/basicdata"
	"github.com/valyala/fasthttp"
	"regexp"
	"strconv"
	"strings"
)

type Reverse struct {
	metadata *basicdata.GatewayMetadata
	config   *basicdata.Syssettingconfig
	ctx      *fasthttp.RequestCtx
	rData    *RequestData
	//valid    *Valid
}

//节点处理结果
type jiSuStruct struct {
	//响应内容
	resBody []byte
	//错误信息
	err error
}

//初始化数据
func (r Reverse) NewRevers(data *RequestData) *Reverse {
	rev := &Reverse{}
	rev.rData = data
	rev.metadata = data.Metadata
	rev.config = data.Config
	rev.ctx = data.Ctx
	//rev.valid = valid
	return rev
}

func (r *Reverse) ReverseProxyGet() {
	//取得适配协议
	uri := r.rData.Uri
	targeturl := ""
	//遍历协议适配
	for _, value := range r.metadata.ProtoURLGET {
		preg := "^" + value.SourceURL + "$"
		reg, _ := regexp.Compile(preg)
		tmplist := reg.FindStringSubmatch(uri)
		if len(tmplist) == 0 {
			continue
		}
		util.LogDebug(value)
		//取得转发的具体地址
		targeturl = value.TargetURL
		for k, val := range tmplist {
			if k == 0 {
				continue
			}
			targeturl = strings.Replace(targeturl, "$"+strconv.Itoa(k), val, -1)
		}

		break
	}

	if targeturl == "" {
		util.LogDebug("未找到GET转跳协议", uri, r.rData.Method)
		Render{}.RendErr(r.ctx, "未找到GET转跳协议")
		return
	}
	util.LogDebug("转发数据：", targeturl)
	r.ctx.Redirect(targeturl, 302)
	return

}
