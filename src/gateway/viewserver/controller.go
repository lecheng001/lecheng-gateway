package viewserver

import (
	"util"
	"encoding/json"
	"gateway/code/proxy"
	"github.com/valyala/fasthttp"
	"html/template"
	"strconv"
	"strings"
	"time"
)

type Controller struct {
	*fasthttp.RequestCtx
}

//取得POST数据
func (c Controller) GetPostData() map[string]interface{} {
	tmp := map[string]interface{}{}
	content := c.PostBody()
	//content, _ := ioutil.ReadAll(body)
	//fmt.Println("Post_payload:", string(content), c.Request().ContentLength)
	//json数据封装到user对象中
	err := json.Unmarshal(content, &tmp)
	if err != nil {
		return nil
	}

	return tmp
}

//取得POST数据，只能一层 返回map[string]string
func (c Controller) GetPostDataString() map[string]string {
	tmp := map[string]interface{}{}
	content := c.PostBody()
	//content, _ := ioutil.ReadAll(body)
	//fmt.Println("Post_payload:", string(content), c.Request().ContentLength)
	//json数据封装到user对象中
	err := json.Unmarshal(content, &tmp)
	if err != nil {
		return nil
	}

	tmp2 := map[string]string{}
	for key, value := range tmp {
		tmp2[key] = util.String(value)
	}
	return tmp2
}

// GetString returns the input value by key string or the default value while it's present and input is blank
func (c Controller) GetString(key string, def ...string) string {
	if v := string(c.QueryArgs().Peek(key)); v != "" {
		return v
	}
	if len(def) > 0 {
		return strings.TrimSpace(def[0])
	}
	return ""
}

// GetInt64 returns input value as int64 or the default value while it's present and input is blank.
func (c Controller) GetInt64(key string, def ...int64) int64 {
	if v := string(c.QueryArgs().Peek(key)); v != "" {
		tmp, _ := strconv.ParseInt(v, 10, 64)
		return tmp
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// GetInt64 returns input value as int64 or the default value while it's present and input is blank.
func (c Controller) GetExist(key string) bool {
	tmp := c.QueryArgs().Peek(key)
	return tmp != nil
}

// GetBool returns input value as bool or the default value while it's present and input is blank.
func (c Controller) GetBool(key string, def ...bool) bool {
	if v := string(c.QueryArgs().Peek(key)); v != "" {
		tmp, _ := strconv.ParseBool(v)
		return tmp
	}
	if len(def) > 0 {
		return def[0]
	}
	return false
}

// GetFloat returns input value as float64 or the default value while it's present and input is blank.
func (c Controller) GetFloat(key string, def ...float64) float64 {
	if v := string(c.QueryArgs().Peek(key)); v != "" {
		tmp, _ := strconv.ParseFloat(v, 64)
		return tmp
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (c Controller) RenderView(tpl string, data map[string]interface{}) error {
	//判断是否登录，未登录时将data置空
 	token:=c.RequestCtx.Request.Header.Cookie("token")
 	if len(token)==0{
 		data=nil
	}
	userToken:= c.GetUserToken()
	if userToken==nil{
		util.Cookie_SetFasthttp(c.RequestCtx,"token","",-1)
		data=nil
	}
	rootpath := util.IO_GetRootPath()

	if data == nil {
		data = map[string]interface{}{}
	}

	data["gateway_package_choice"]= map[string]interface{}{
		"package":proxy.GWPACKAGE.Projecttype,"version":proxy.GWPACKAGE.Version,"projectid":proxy.GWPACKAGE.Projectid,"projecttitle":proxy.GWPACKAGE.ProjectName,"userid":proxy.GWPACKAGE.Userid,
	}

	if _, ok := data["now"]; !ok {
		data["now"] = time.Now().Unix()
	}

	date := func(value interface{}) (string, error) {
		return util.String(util.Date_FormatYmd2(value)), nil
	}
	datetime := func(value interface{}) (string, error) {
		return util.String(util.Date_FormatYmdHis(value)), nil
	}
	tmpl := template.Must(template.New("base.gohtml").Funcs(template.FuncMap{
		"date": date, "datetime": datetime,
	}).Delims("{[{","}]}").ParseFiles(rootpath+"/views/base.gohtml", rootpath+"/views/"+tpl))

	return tmpl.Execute(c.Response.BodyWriter(), data)
}

func (c Controller) RenderRedirect(alertMessage string, url string) {
	if len(url) <= 0 {
		url = "/"
	}
	html := "<script>" +
		"alert('" + alertMessage + "');" +
		"location.href='" + url + "'" +
		"</script>"
	c.HTML(html)
}

func (c Controller) RenderJSON(data map[string]interface{}) {
	if data == nil {
		data = map[string]interface{}{"errcode": -1, "msg": "数据加载失败"}
	} else {
		if data["errcode"] == nil {
			data["errcode"] = 0
		}
		if data["msg"] == nil {
			data["msg"] = ""
		}
	}

	c.JSON(data)
}

func (c Controller) RenderError(msg, url string) {
	c.RenderJSON(map[string]interface{}{"errcode": -1, "msg": msg, "url": url})
}

func (c Controller) GetUserToken() map[string]string {
	//cookie := fasthttp.Cookie{}
	cvalue:=c.RequestCtx.Request.Header.Cookie("token")
	token := util.Token_Analysis(string(cvalue))
	if token == nil {
		return nil
	} else {
		return token
	}
}

func (c Controller) HTML(html string) {
	// set some headers and status code first
	c.SetContentType("text/html")
	c.SetStatusCode(fasthttp.StatusOK)

	// then override already written body
	c.SetBody([]byte(html))
}
func (c Controller) JSON(jsonValue map[string]interface{}) {
	// set some headers and status code first
	c.SetContentType("application/json")
	c.SetStatusCode(fasthttp.StatusOK)

	// then override already written body
	tmpValue, _ := json.Marshal(jsonValue)
	c.SetBody(tmpValue)
}
