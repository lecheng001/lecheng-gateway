/**
文件上传
开发人员：陈朝能
*/
package proxy

import (
	"bytes"
	"util"
	"errors"
	"gateway/code/basicdata"
	"github.com/valyala/fasthttp"
	"io"
	"math/rand"
	"mime/multipart"
	"path"
	"regexp"
	"strings"
	"time"
)

type Upload struct {
	rData    *RequestData
	metadata *basicdata.GatewayMetadata
	config   *basicdata.Syssettingconfig
	ctx      *fasthttp.RequestCtx
}

func (_ Upload) NewUpload(data *RequestData) *Upload {
	u := &Upload{}
	u.rData = data
	u.metadata = data.Metadata
	u.config = data.Config
	u.ctx = data.Ctx

	return u
}

//匹配上传规则
func (u *Upload) GetUploadRule() error {
	if !u.rData.UpFile_HasFile {
		return nil
	}

	if len(u.metadata.UploadURL) == 0 {
		return errors.New("网关上传协议未配置")
	}

	//取得适配协议
	uri := strings.ToLower(u.rData.Uri)
	//targeturl := ""
	//hosts := []basicdata.ModuleServer{}
	uploadInfo := basicdata.UploadURL{}
	//遍历协议适配
	for _, value := range u.metadata.UploadURL {
		if uri == strings.ToLower(value.SourceURL) {
			uploadInfo = value
			break
		}

		//完全匹配
		if strings.Contains(value.SourceURL, "(") {
			//正则匹配
			preg := "^" + value.SourceURL + "$"
			reg, _ := regexp.Compile(preg)
			tmplist := reg.FindStringSubmatch(uri)
			if len(tmplist) == 0 {
				continue
			} else if len(tmplist) > 0 {
				uploadInfo = value
				break
			}
		}
	}
	util.LogDebug(uploadInfo)

	if uploadInfo.TargetURL == "" {
		return errors.New("未找到目标URL")
	}
	u.rData.UpFile_UploadInfo = &uploadInfo
	return nil
}

//图片上传，无需验证,反向代理，下一版最好能在加上安全性验证
func (u *Upload) UploadFile() error {
	ctx := u.ctx

	//fm :=u.rData.MultiForm  // ctx.MultipartForm()
	fiels := u.rData.UpFile // fm.File // u.rData.MultiForm.File

	if len(fiels) == 0 {
		return errors.New("请选择上传文件")
	}
	if u.rData.UpFile_UploadInfo.UpNum>0 && int64(len(fiels)) > u.rData.UpFile_UploadInfo.UpNum {
		return errors.New("最多支持"+util.String(u.rData.UpFile_UploadInfo.UpNum)+"个文件上传")
	}

	//targeturl := ""
	//hosts := []basicdata.ModuleServer{}
	uploadInfo := u.rData.UpFile_UploadInfo
	bodyBufer := &bytes.Buffer{}
	//创建一个multipart文件写入器，方便按照http规定格式写入内容
	bodyWriter := multipart.NewWriter(bodyBufer)
	//增加mulitpart的参数
	if len(u.rData.ParamsMap) > 0 {
		for key, value := range u.rData.ParamsMap {
			bodyWriter.WriteField(key, util.String(value))
		}
	} else if u.rData.IsBase64Encry {
		bodyWriter.WriteField(ENCRYPTPARAMNAME, u.rData.ParamsEnCodeString)
	}

	for key, value := range fiels {
		//if len(value) != 1 {
		//	Render{}.RendErr(ctx, "只支持一个文件上传")
		//	return
		//}

		//判断文件后缀名
		if uploadInfo.FileFormat != "" {
			if !strings.Contains("|"+uploadInfo.FileFormat+"|", "|"+path.Ext(value[0].Filename)+"|") {
				//Render{}.RendString(ctx, []byte())
				return errors.New(uploadInfo.ReturnFormatCotnent)
			}
		}

		//判断文件大小
		if uploadInfo.MaxFileSize > 0 {
			if uploadInfo.MaxFileSize*1024 < value[0].Size {
				//Render{}.RendString(ctx, []byte())
				return errors.New(uploadInfo.ReturnSizeContent)
			}
		}

		file, err := value[0].Open()
		//新建一个缓冲，用于存放文件内容

		//从bodyWriter生成fileWriter,并将文件内容写入fileWriter,多个文件可进行多次
		fileWriter, err := bodyWriter.CreateFormFile(key, value[0].Filename)
		if err != nil {
			util.LogError(err.Error())
			return err
		}
		//file,err := os.Open(uploadFile)
		if err != nil {
			util.LogError(err)
			return err
		}
		//不要忘记关闭打开的文件
		defer file.Close()
		_, err = io.Copy(fileWriter, file)
		if err != nil {
			util.LogError(err.Error())
			return  err
		}

		//关闭bodyWriter停止写入数据
		bodyWriter.Close()

	}

	contentType := bodyWriter.FormDataContentType()

	trueServer := ""
	if len(uploadInfo.Servers) == 0 {
		//Render{}.RendErr(ctx, "微服务服务器未配置")
		return errors.New("微服务服务器未配置")
	} else if len(uploadInfo.Servers) == 1 {
		trueServer = uploadInfo.Servers[0].Host
	} else {
		//随机负载
		trueServer = uploadInfo.Servers[rand.Int()%len(uploadInfo.Servers)].Host
	}

	t:=time.Now()
	u.rData.Time_starttime=t
	//构建request，发送请求
	request := fasthttp.AcquireRequest()
	request.Header.SetContentType(contentType)
	//直接将构建好的数据放入post的body中
	request.SetBody(bodyBufer.Bytes())
	request.Header.SetMethod("POST")
	request.URI().SetHost(trueServer)
	request.URI().SetPath(uploadInfo.TargetURL)
	response := fasthttp.AcquireResponse()
	err := fasthttp.DoTimeout(request, response, time.Second*time.Duration(uploadInfo.Timeout))
	//fmt.Println(string(response.Body()))

	if err != nil {
		ctx.Logger().Printf("error when proxying the request: %s", err)
		if err.Error() == "timeout" {
			response.SetStatusCode(fasthttp.StatusRequestTimeout)
		}
	}
	status := response.StatusCode()
	u.rData.ResponseByte= response.Body()
	u.rData.Response=string(u.rData.ResponseByte)
	if status != 200 {
		util.LogError(string(response.Body()))
		if response.StatusCode() == fasthttp.StatusRequestTimeout {
			//Render{}.RendString(ctx, []byte(uploadInfo.ReturnTimeoutContent))
			u.rData.Time_endtime=time.Now()
			return errors.New(uploadInfo.ReturnTimeoutContent)
		} else if status == 404 {
			//Render{}.RendString(ctx, []byte(u.config.NotFindContent))
			u.rData.Time_endtime=time.Now()
			return  errors.New(u.config.NotFindContent)
		} else {
			//tmpcontent := ""
			//if err != nil {
			//	tmpcontent = err.Error()
			//} else {
			//	tmpcontent = string(response.Body())
			//}
			//Render{}.RendErr(ctx, tmpcontent)
			//Render{}.RendString(ctx, []byte(uploadInfo.ReturnErrorContent))
			u.rData.Time_endtime=time.Now()
			return errors.New(uploadInfo.ReturnErrorContent)
		}
	}
	u.rData.Time_endtime=time.Now()
	u.rData.Time_elapsedtime=time.Since(t)
	return nil
}
