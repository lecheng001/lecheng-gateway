/**
开发人员：陈朝能
*/
package proxy

import (
	"util"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gateway/code/basicdata"
	"github.com/buger/jsonparser"
	"github.com/valyala/fasthttp"
	"sort"
	"strings"
)

type Valid struct {
	rData    *RequestData
	metadata *basicdata.GatewayMetadata
	config   *basicdata.Syssettingconfig
	ctx      *fasthttp.RequestCtx
}

//初始化数据
func (v Valid) NewValid(data *RequestData) *Valid {
	valid := &Valid{}
	valid.rData = data
	valid.metadata = data.Metadata
	valid.config = data.Config
	valid.ctx = data.Ctx
	return valid
}

//验证合法性，包含：验签、usertoken校验
func (v *Valid) CheckValid() error {
	//验证Client
	clientItem := v.rData.ApplicationClient

	if clientItem.Salt == "" {
		util.LogError("_LCC_ err::", clientItem)
		return errors.New("请求客户端头部数据配置有误,验签盐值不能为空")
	}
	//判断头部client与applicationClient是否一致
	//if v.rData.UriScheme[0]!=v.metadata

	//如果没有验签与AES加密，返回原请求内容，不做JSON判断
	if clientItem.EnableSign == 0 && clientItem.EnableBASE64 == 0 {
		if len(v.rData.ParamsMap) == 0 {
			return errors.New("请求参数不允许为空")
		}
		return nil
	}

	//判断BAse64加密，并解密,结果存在 v.rData.Params
	if clientItem.EnableBASE64 == 1 {
		_, err := v.CheckBase64(clientItem)
		if err != nil {
			util.LogError("Base64解密失败:", err.Error(), "__请求参数：", string(v.ctx.Request.Body()))
			return err
		}
		paramsStr := v.rData.Params
		//err := errors.New("")
		if v.rData.IsBodyJson {
			params := map[string]interface{}{}
			err = json.Unmarshal([]byte(paramsStr), &params)
			if err != nil {
				util.LogError("POST数据必须是json格式:", v.rData.ParamsEnCode, paramsStr)
				return errors.New("POST数据必须是json格式")
			} else {
				v.rData.ParamsMap = params
			}
		} else if v.rData.IsBodyForm {
			v.rData.ParamsMap, err = util.String_AnalysisURI(paramsStr)
			if err != nil {
				return err
			}
		}

	} else {
		//未加密，
		v.rData.Params = string(v.rData.ParamsEnCode)
	}

	//将参数字符串解析成map数据
	v._paramsToMap()

	//MD5验签
	if clientItem.EnableSign == 1 {
		//验签参数不存在
		if _, ok := v.rData.ParamsMap[SIGNNAME]; !ok {
			util.LogError("必需参数丢失", v.rData.Params)
			return errors.New("必需参数丢失")
		}

		err := v.CheckSign(clientItem.Salt, v.rData.ParamsMap, v.rData.Params)
		if err != nil {
			return err
		}
		return nil
	}

	//不需要验签则返回正式
	return nil
}

//AES验签
func (v *Valid) CheckBase64(clientApplication *basicdata.ApplicationClient) (bool, error) {
	unCodeStr, err := base64.StdEncoding.DecodeString(v.rData.ParamsEnCodeString)
	//if err != nil {
	//	return false
	//}
	//body, err := util.String_IsBase64(v.rData.ParamsEnCodeString, []byte( clientApplication.AESCode))
	if err != nil {
		return false, errors.New("解密失败，请核对密钥")
	}
	if unCodeStr == nil {
		return false, errors.New("解密失败，解密结果为空")
	}
	v.rData.Params = string(unCodeStr)
	return true, nil
}

//MD5验签,成功后删除lcsign属性
func (v *Valid) CheckSign(encryptcode string, params map[string]interface{}, paramsStr string) error {
	util.LogDebug(params)
	sign := params[SIGNNAME].(string)
	delete(params, SIGNNAME)
	keys := []string{}

	//转换成单层
	params2 := v.changeToSingleJSON(params, paramsStr)
	//key排序
	for key, _ := range params2 {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	content := ""
	for _, value := range keys {
		str := util.String(params2[value])
		content += value + "=" + str + "&"
	}
	content += "SECRET=" + encryptcode

	util.LogDebug("MD5：", content)
	md5 := util.MD5(content)
	if md5 != sign {
		msg := ""
		if util.ConfigGetString("env") != "pro" {
			msg = fmt.Sprintf("oldSign:%s,newSign: %s,MD5：%s", sign, md5, content)
		}
		tmp, _ := json.Marshal(params)
		util.LogError("验签失败 "+msg, "__请求数据：", string(tmp), "__lecheng-client:", string(v.ctx.Request.Header.Peek("lecheng-client")))
		return errors.New("验签失败 " + msg)
	}

	delete(params, SIGNNAME)
	v.rData.ParamsMap = params
	//_, err := json.Marshal(params)
	//if err != nil {
	//	return errors.New( "json转换失败")
	//}
	return nil
}

//转换成单层参数模式
func (v *Valid) changeToSingleJSON(params map[string]interface{}, paramsStr string) map[string]interface{} {
	result := map[string]interface{}{}
	for key, value := range params {
		v._changeToSingleJSON(key, value, &result, paramsStr)
	}

	return result
}

//递归取得验签参数
func (v *Valid) _changeToSingleJSON(pkey string, pvalue interface{}, result *map[string]interface{}, paramsStr string) {
	if val, ok := pvalue.(map[string]interface{}); ok {
		//空对象
		if len(val) == 0 {
			(*result)[pkey] = "null"
			return
		}
		for key, value := range val {
			//(*result)[pkey+"."+k] = v
			v._changeToSingleJSON(pkey+"."+key, value, result, paramsStr)
		}
	} else if _, ok := pvalue.([]interface{}); ok {
		//特殊处理数组对象，不能对数组对象进行排序，会导致签名失败
		tmpa, tmp, tmpc, tmpd := jsonparser.Get([]byte(paramsStr), strings.Split(pkey, ".")...)
		util.LogDebug(tmpa, tmp, tmpc, tmpd)

		(*result)[pkey] = tmpa
	} else {
		(*result)[pkey] = pvalue
	}
}

//将最后需要解析的参数字符串拆解成
func (v *Valid) _paramsToMap() {
	if len(v.rData.ParamsMap) > 0 {
		return
	}

	if v.rData.IsBodyJson {
		tmpmap, err := util.String_GetJSON(v.rData.Params)
		if err != nil {
			v.rData.ParamsMap = tmpmap
		}
	} else {
		tmpArr := strings.Split(v.rData.Params, "&")
		for _, value := range tmpArr {
			tmpArr2 := strings.Split(value, "=")
			if len(tmpArr2) > 1 {
				v.rData.ParamsMap[tmpArr2[0]] = tmpArr2[1]
			} else {
				v.rData.ParamsMap[tmpArr2[0]] = ""
			}
		}
	}
}
