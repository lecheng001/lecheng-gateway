package proxy

import (
	"bytes"
	"util"
	"encoding/json"
	"errors"
	"gateway/code/basicdata"
	"github.com/valyala/fasthttp"
	"mime/multipart"
	"net/url"
	"strings"
	"time"
)

type RequestData struct {
	Metadata          *basicdata.GatewayMetadata
	Config            *basicdata.Syssettingconfig
	Ctx               *fasthttp.RequestCtx
	ApplicationClient *basicdata.ApplicationClient
	//Header  fasthttp.RequestHeader
	Header map[string]string
	//请求方式
	HeaderContentType string
	userAgent         []byte
	//只在测试模式下有效
	Test_host string
	//  eg:/aa/bb.html?aa=bb
	Uri string
	//  eg:http://www.xxx.com/aa/bb.html?aa=bb
	UriFull string
	// eg: /aa/bb.html
	UriPath string
	// eg: aa=bb
	UriQuery  string
	UriScheme []string
	Method    string
	//已解密
	RequestBody string
	//密文
	RequestBodyEncode       []byte
	RequestBodyEncodeString string
	//已解密
	RequestQuery string
	//密文
	RequestQueryEncode       []byte
	RequestQueryEncodeString string

	//文件
	//MultiForm *multipart.Form
	UpFile_HasFile    bool
	UpFile            map[string][]*multipart.FileHeader
	UpFile_UploadInfo *basicdata.UploadURL

	//参数加密
	ParamsEnCode       []byte
	ParamsEnCodeString string
	//参数
	Params string
	//请求MAP数据
	ParamsMap map[string]interface{}

	IP       string
	IsGet    bool
	IsPost   bool
	IsPut    bool
	IsDelete bool
	//对GET处理事务有用
	IsRESTfull          bool
	IsBodyMultipartForm bool
	IsBodyJson          bool
	IsBodyForm          bool
	//是否base64密文传输
	IsBase64Encry bool
	//数据提交方式只有两种：post提交数据，GET参数提交
	IsBodyTransmit bool

	HeaderClient  string
	HeaderVersion int64

	Response     string
	ResponseByte []byte
	//开始时间
	Time_starttime time.Time
	//结束时间
	Time_endtime time.Time
	//耗时
	Time_elapsedtime time.Duration
}

//初始化执行参数，只是基本复制，不错业务判断
func (r RequestData) GetInitRequestData(ctx *fasthttp.RequestCtx) (*RequestData, error) {
	//支持RESTfull请求，GET、POST、PUI、DELET，
	//其中POST、PUT有body

	rData := RequestData{}
	rData.IsGet = ctx.IsGet()
	rData.IsPost = ctx.IsPost()
	rData.IsPut = ctx.IsPut()
	rData.IsDelete = ctx.IsDelete()
	if !rData.IsGet && !rData.IsPost && !rData.IsPut && !rData.IsDelete {
		return &rData, errors.New("不支持改请求方式：" + string(ctx.Method()))
	}

	rData.Ctx = ctx
	rData.Method = string(ctx.Method())
	rData.UriFull = ctx.Request.URI().String()
	rData.Uri = string(ctx.Request.RequestURI())
	rData.UriPath = string(ctx.Request.URI().Path())
	rData.UriQuery = string(ctx.Request.URI().QueryString())
	rData.UriScheme = strings.Split(strings.Trim(rData.UriPath, "/"), "/")

	//系统后台业务操作接口，与网关无关系
	if "system" == rData.UriScheme[0] {
		return &rData, nil
	}

	//测试模式
	if util.ConfigGetString("env") != "pro" {
		rData.Test_host = string(ctx.Request.Header.Peek("lecheng-host"))
		if !util.ConfigGetBool("setting.enable_lecheng_host") {
			rData.Test_host = ""
		}
	}

	tmpHeader := fasthttp.RequestHeader{}
	ctx.Request.Header.CopyTo(&tmpHeader)
	//rData.Header = tmpHeader
	rData.Header = r.ansyHead(tmpHeader.Header())
	rData.HeaderContentType = string(ctx.Request.Header.ContentType())
	rData.userAgent = tmpHeader.UserAgent()
	rData.HeaderClient = strings.TrimSpace(rData.Header["Lecheng-Client"])
	if rData.HeaderClient == "" {
		rData.HeaderClient = strings.TrimSpace(rData.Header["lecheng-client"])
	}

	rData.IsGet = ctx.IsGet()
	rData.IsPost = ctx.IsPost()
	rData.IsPut = ctx.IsPut()
	rData.IsDelete = ctx.IsDelete()

	rData.RequestQueryEncode = ctx.QueryArgs().QueryString()
	rData.RequestQueryEncodeString = string(rData.RequestQueryEncode)
	//BODY_json都是RESTFull模式
	rData.RequestBodyEncode = ctx.Request.Body()
	rData.RequestBodyEncodeString = string(rData.RequestBodyEncode)
	if len(rData.RequestBodyEncode) > 0 {
		rData.IsRESTfull = true
	}

	//判断RESTfull
	if rData.IsGet {
		lvcontent := ctx.QueryArgs().Peek(ENCRYPTPARAMNAME)
		if len(lvcontent) > 0 {
			rData.IsRESTfull = true
		}
	} else if rData.IsPost || rData.IsPut {
		rData.IsRESTfull = true
	} else if rData.IsDelete {
		rData.IsRESTfull = true
	}

	if rData.IsPost {
		if strings.Contains(rData.Header["Content-Type"], "multipart") {
			rData.IsBodyMultipartForm = true

			rData.RequestBodyEncode = ctx.Request.Body()
			rData.RequestBodyEncodeString = string(rData.RequestBodyEncode)

		}
	}

	//判断头部信息
	if rData.IsRESTfull {
		if rData.HeaderClient == "" && !rData.IsBodyMultipartForm {
			return &rData, errors.New("请求客户端头部数据丢失")
		}
	}

	//版本号
	rData.HeaderVersion = util.String_GetInt64(rData.Header["Lecheng-Version"])
	if rData.HeaderVersion == 0 {
		rData.HeaderVersion = util.String_GetInt64(rData.Header["Lecheng-version"])
	}
	//IP
	rData.IP = util.Request_GetIP_Fasthttp(ctx)

	return &rData, nil
}

//判断请求URL地址和client是否匹配
func (r *RequestData) CheckApplication() error {
	//上传的如果没有数据去掉验证
	if r.IsRESTfull {
		if r.HeaderClient == "" && r.IsBodyMultipartForm {
			return nil
		}
	}

	path := r.UriScheme[0]

	for _, value := range r.Metadata.ApplicationClient {
		if r.HeaderClient == value.Agentkey {
			if path == GATEWAYAPIPATH {
				r.ApplicationClient = &value
				return nil
			} else if value.Application.Path == path {
				r.ApplicationClient = &value
				return nil
			} else {
				//请求地址配置的Applicaton和client无法匹配，Header信息与APplication不匹配
				return errors.New("请求入口H&A有误，请正确配置")
			}
		}
	}

	return errors.New("请求入口有误，找不到可用的数据")
}

//分析请求基础数据.GET,POST ,PUT ,DELTE
func (r *RequestData) AnsyRequestData() error {
	//处理请求数据 querystring，ispayload，payloaddata
	//r.RequestQueryEncode = r.Ctx.QueryArgs().QueryString()
	//r.RequestBodyEncode = r.Ctx.Request.Body()

	r.ParamsMap = map[string]interface{}{}
	if len(r.RequestBodyEncode) > 0 {
		r.IsBodyTransmit = true
		r.ParamsEnCode = r.RequestBodyEncode
		r.ParamsEnCodeString = string(r.ParamsEnCode)
	} else if r.IsRESTfull {

		if r.Ctx.QueryArgs().Len() == 0 {
			return errors.New("请求参数丢失！")
		}
		tmpArgs := map[string]interface{}{}
		//标注body_form请求，头部包含"Content-Type":"application/x-www-form-urlencoded"
		r.Ctx.QueryArgs().VisitAll(func(key, value []byte) {
			tmpArgs[string(key)] = string(value)
		})
		tmpValue := tmpArgs[ENCRYPTPARAMNAME]
		if tmpValue != nil {
			if util.String_IsBase64(tmpValue.(string)) {
				r.ParamsEnCodeString = tmpValue.(string)
				r.ParamsEnCode = []byte(r.ParamsEnCodeString)
				return nil
			}
			//delete(tmpArgs, ENCRYPTPARAMNAME)
		}

		r.ParamsMap = tmpArgs

		//if strings.Contains(r.RequestQueryEncodeString, ENCRYPTPARAMNAME+"=") {
		//	tmp, err := url.QueryUnescape(strings.Replace(r.RequestQueryEncodeString, ENCRYPTPARAMNAME+"=", "", -1))
		//	if err != nil {
		//		return err
		//	}
		//	if util.String_IsBase64(tmp) {
		//		r.IsBase64Encry = true
		//		r.ParamsEnCodeString = tmp
		//		r.ParamsEnCode = []byte(tmp)
		//	}
		//} else {
		//	r.ParamsEnCode = r.RequestQueryEncode
		//	r.ParamsEnCodeString = string(r.ParamsEnCode)
		//}
		return nil
	}

	if len(r.ParamsEnCode) == 0 {
		return errors.New("请求参数丢失！")
	}

	//BODY 有数据
	if r.IsBodyTransmit {
		//1：严格请求
		//2：非严格
		//body_form 1:明文，2：加密:lccontent=xxxxxxxxxxxxx
		//body_json 1:明文，2：加密:xxxxxxxxxxxxxxxxxxxxxxxxxxxx 或：{"data":"xxxxxxxxxxxxxxxxxxxxxxxxxx"}
		//POST-multipart 1:明文aa=bb&cc=dd，2：加密：lccontent =xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
		if r.IsBodyMultipartForm || r.Ctx.PostArgs().Len() > 0 {
			r.IsBodyForm = true
			if r.IsBodyMultipartForm {
				multiForm, err := r.Ctx.MultipartForm()
				if err != nil {
					return errors.New("文件上传失败！" + err.Error())
				}
				r.UpFile = multiForm.File
				if len(r.UpFile) > 0 {
					r.UpFile_HasFile = true
				}
				//r.UpFile = mf.File
				//tmpmap:= map[string]interface{}{}
				for key, value := range multiForm.Value {
					//tmpmap[key]=strings.Join(value,",")
					tmpp:=value[0]
					tmpp=strings.Replace(tmpp, key+"=","",-1)
					tmpp,_= url.QueryUnescape(tmpp)
					r.ParamsMap[key] =tmpp // strings.Join(tmpp, ",")
				}

				r.ParamsEnCode, _ = json.Marshal(r.ParamsMap)
				r.ParamsEnCodeString = string(r.ParamsEnCode)
			}

			if r.Ctx.PostArgs().Len() > 0 {

				//标注body_form请求，头部包含"Content-Type":"application/x-www-form-urlencoded"
				r.Ctx.PostArgs().VisitAll(func(key, value []byte) {
					r.ParamsMap[string(key)] = string(value)
				})
			}

			//只有一个长度是判断是否是加密
			if len(r.ParamsMap) == 1 {
				if value, ok := r.ParamsMap[ENCRYPTPARAMNAME]; ok {
					if util.String_IsBase64(value.(string)) {
						r.ParamsMap = map[string]interface{}{}
						r.IsBase64Encry = true
						r.ParamsEnCode = []byte(value.(string))
						r.ParamsEnCodeString = value.(string)
						return nil
					}

					return errors.New("请求参数有误")
				}
			}

			return nil
		}

		//1:判断base64，
		// 是：一定是json，
		// 否：解析json
		//       解析成功：body_json 明文或者{"data":"xxxxxxxxxxxxxxxxx"}
		//       解析失败：keyvalue内容解析（加密：base64 和非加密） body_form

		if util.String_IsBase64(util.String(r.RequestBodyEncode)) {
			r.IsBodyJson = true
			r.ParamsEnCode = r.RequestBodyEncode
			return nil
		}

		//非base64加密：
		tmpAA := map[string]interface{}{}
		err := json.Unmarshal(r.RequestBodyEncode, &tmpAA)
		//body_json格式数据(未加密、加密（data））
		if err == nil {
			r.IsBodyJson = true
			//约定：{"data":"xxxxxxxxxxxxxxxxx"}，或者xxxxxxxxxxxxxxxxxxxxxxxxxxx
			if value, ok := tmpAA[JSONPARAMNAME]; ok {
				if val, ok2 := value.(string); ok2 {
					r.IsBase64Encry = util.String_IsBase64(val)
					if r.IsBase64Encry {
						r.RequestBodyEncode = []byte(val)
						r.ParamsEnCode = r.RequestBodyEncode
						r.ParamsEnCodeString=string(r.ParamsEnCode)
						return nil
					}
				}
			}
			r.ParamsMap = tmpAA

			return nil
		}

		r.IsBodyForm = true

		//json解析失败，body_form类型
		r.ParamsMap, err = util.String_AnalysisURI(string(r.RequestBodyEncode))
		if err != nil {
			return err
		} else {
			if value, ok := r.ParamsMap[ENCRYPTPARAMNAME]; ok {
				r.IsBase64Encry = true
				r.ParamsEnCodeString=value.(string)
				r.ParamsEnCode = []byte(r.ParamsEnCodeString)
				r.ParamsMap = map[string]interface{}{}
				return nil
			}
		}

		return nil
	}

	return nil
}

//分析Header头部信息
func (r RequestData) ansyHead(content []byte) map[string]string {
	//fmt.Println(string(content))
	tmp := content
	tmp = bytes.ReplaceAll(tmp, []byte{'\r'}, []byte{})
	//result := []string{""}
	tmpmap := map[string]string{}
	index := 0
	for {
		a := bytes.IndexByte(tmp, '\n')
		if a <= 0 {
			break
		} else {
			aa := string(tmp[index:a])
			if aa != "" {
				tmparr := strings.Split(aa, ":")
				if len(tmparr) > 1 {
					tmpmap[tmparr[0]] = tmparr[1]
				}
			}
			tmp = tmp[a+1:]
		}
	}
	return tmpmap
}
