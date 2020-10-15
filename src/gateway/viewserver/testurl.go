package viewserver

import (
	"util/jsondb"
	"util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
)

type TestUrlController struct {
	Controller
}

func (c TestUrlController) jsonfile() string {
	return "database/testurl.json"
}

func (c TestUrlController) TestUrl() {

	data := map[string]interface{}{}
	where := map[string]interface{}{}
	//where["cstatus|>="] = 0
	//where["ctype"] = "test"
	tmplist,_:=jsondb.Query(c.jsonfile(),where,nil)
	data["loglist"] = tmplist
	c.RenderView("testurl.gohtml", data)
}

func (c TestUrlController) GetList() {

	data := map[string]interface{}{}
	where := map[string]interface{}{}
	//where["cstatus|>="] = 0
	//where["ctype"] = "test"
	tmplist,_:=jsondb.Query(c.jsonfile(),where,nil)
	data["loglist"] = tmplist
	c.RenderJSON( data)
}

//保存请求数据
func (c TestUrlController) PostData() {
	postdata := c.GetPostData()
	id:=util.String_GetInt64(postdata["pkid"])
	err:=errors.New("")
	if id>0{
	_,err=	jsondb.UpdateByID(c.jsonfile(),id,postdata)
	}else {
		_,err=	jsondb.InsertMap(c.jsonfile(),postdata)
	}
	if err != nil {
		c.RenderJSON(map[string]interface{}{"msg": "保存失败" + err.Error()})
		return
	}

	c.RenderJSON(map[string]interface{}{"msg": "保存成功", "url": "location"})
}

func (c TestUrlController) Delete() {
	postdata := c.GetPostData()
	id:=util.String_GetInt64(postdata["pkid"])
	err:=errors.New("")
	if id>0{
		//_,err=	jsondb.UpdateByID(c.jsonfile(),id,postdata)
		_,err=jsondb.Delete(c.jsonfile(),"pkid",id)
	}

	if err != nil {
		c.RenderJSON(map[string]interface{}{"msg": "保存失败" + err.Error()})
		return
	}

	c.RenderJSON(map[string]interface{}{"msg": "保存成功", "url": "location"})
}

// 取得URL请求数据
func (c TestUrlController) TestGetUrlContent() {
	st := time.Now()
	err, body, status := c.getUrlContent(c.GetPostDataString())
	if err != nil {
		// handle error
		fmt.Println(err.Error())
	}
	ut := time.Since(st)
	//fmt.Println(string(body))
	c.RenderJSON(map[string]interface{}{"status": status, "runtime": ut.String(), "responsecontent": string(body)})

}

//返回URL数据，这个很可能异常，加入defer+recover
func (c TestUrlController) getUrlContent(postdata map[string]string, ) (err error, body []byte, status string) {
	err = errors.New("")
	body = []byte{}
	status = ""
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(util.String(r))
			body = []byte("请求失败，没有响应")
			status = " 404 "
		}
	}()

	client := &http.Client{}
	query := postdata["cquery"]
	if _, ok := postdata["cqueryrsa"]; ok {
		query = postdata["cqueryrsa"]
	}
	util.LogDebug(strings.ToUpper(postdata["cmethod"]), postdata["curl"], query)
	req, err := http.NewRequest(strings.ToUpper(postdata["cmethod"]), postdata["curl"], strings.NewReader(query))
	if err != nil {
		// handle error
	}

	head := postdata["chead"]
	if head != "" {
		tmpmp := map[string]interface{}{}
		json.Unmarshal([]byte(head), &tmpmp)
		if tmpmp != nil {
			for key, value := range tmpmp {
				req.Header.Set(key, util.String(value))
			}
		}
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)
	status = resp.Status
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	return
}

func (c TestUrlController) RsaCode() {
	postdata := c.GetPostData()
	query := util.String(postdata["query"])
	queryrsa := util.String(postdata["queryrsa"])
	head := util.String(postdata["head"])
	//method := util.String(postdata["method"])

	//取得rsa公钥
	tmp := map[string]interface{}{}
	json.Unmarshal([]byte(head), &tmp)
	aescode := ""
	applicationClient := map[string]interface{}{}
	if value, ok := tmp["lecheng-client"].(string); ok && len(value) > 0 {
		applicationClient=jsondb.QueryRow("database/applicationclient.json",map[string]interface{}{"cclientagent": value},nil)
	}

	if applicationClient["pkid"]==nil{
		c.RenderError( "没有找到对应的应用项目客户端，头部需要存在合法的lecheng-client值","")
		return
	}

	if len(query) > 10 {
		tmpContent := query
		tmpResult := tmpContent
		data := map[string]interface{}{}
		md5:=""
		if util.String_GetInt64(applicationClient["cenablesign"]) > 0 {
			sign:=""
			sign,md5, _ = c.getSign(applicationClient["cencryptcode"].(string), tmpContent)
			data["md5sign"] = sign
			tmpResult = sign
			fmt.Println("md5:", tmpResult)
		}

		tmp := map[string]interface{}{}
		json.Unmarshal([]byte(tmpContent), &tmp)
		tmp[util.ConfigGetString("setting.signname")]= md5
		tmp2,_:= json.Marshal(tmp)
		data["value"] = string(tmp2)
		c.RenderJSON(data)
		return
	}
	if len(queryrsa) > 10 {
		tmprsa, _ := util.AES_DecryptBase64(queryrsa, []byte( aescode))
		c.RenderJSON(map[string]interface{}{"value": tmprsa})
		return
	}
	c.RenderJSON(map[string]interface{}{"value": "参数有误"})
}

//MD5验签,成功后删除lcsign属性
func (d TestUrlController) getSign(encryptcode string, paramsStr string) (string,string, error) {
	isbodyjson := false
	params, err := util.String_GetJSON(paramsStr)
	if err != nil {

		params, err = util.String_AnalysisURI(paramsStr)
		if err != nil {
			return "","", err
		}
	} else {
		isbodyjson = true
		//转换成单层
		params = d.changeToSingleJSON(params, paramsStr)
	}

	util.LogDebug(params)
	delete(params, "lcsign")
	keys := []string{}

	//key排序
	for key, _ := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	content := ""
	for _, value := range keys {
		str := util.String(params[value])
		content += value + "=" + str + "&"
	}
	content += "SECRET=" + encryptcode

	util.LogDebug("MD5：", content)
	md5 := util.MD5(content)
	params[util.ConfigGetString("setting.signname")] = md5

	result := ""
	if isbodyjson {
		resulta, _ := json.Marshal(params)
		result = string(resulta)
	} else {
		tmp := ""
		for key, value := range params {
			tmp += key + "=" + util.String(value) + "&"
		}
		result = strings.Trim(tmp, "&")
	}

	return result,md5, nil
}

//转换成单层参数模式
func (v TestUrlController) changeToSingleJSON(params map[string]interface{}, paramsStr string) map[string]interface{} {
	result := map[string]interface{}{}
	for key, value := range params {
		v._changeToSingleJSON(key, value, &result, paramsStr)
	}

	return result
}

//递归取得验签参数
func (d TestUrlController) _changeToSingleJSON(pkey string, pvalue interface{}, result *map[string]interface{}, paramsStr string) {
	if val, ok := pvalue.(map[string]interface{}); ok {
		//空对象
		if len(val) == 0 {
			(*result)[pkey] = "null"
			return
		}
		for key, value := range val {
			//(*result)[pkey+"."+k] = v
			d._changeToSingleJSON(pkey+"."+key, value, result, paramsStr)
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
