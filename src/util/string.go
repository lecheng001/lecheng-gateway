package util

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"go/types"
	"html"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func String_GetInt64(val interface{}) int64 {
	value, err := strconv.ParseInt(fmt.Sprintf("%v", val), 10, 64)
	if err != nil {
		return 0
	}
	return value

}
func String_GetFloat(val interface{}) float64 {
	value, err := strconv.ParseFloat(fmt.Sprintf("%v", val), 64)
	if err != nil {
		return 0
	}
	return value
}
func String_GetBool(val interface{}) bool {
	value, err := strconv.ParseBool(fmt.Sprintf("%v", val))
	if err != nil {
		return false
	}
	return value

}

func String_GetInt(val interface{}) int {
	//fmt.Println(reflect.TypeOf(val))
	if value, ok := val.(string); ok {
		tmp, _ := strconv.Atoi(value)
		return tmp
	} else if value, ok := val.(int); ok {
		return value
	}else if value, ok := val.(int64); ok {
		tmp:=strconv.FormatInt(value,10)
		tmp2, _ := strconv.Atoi(tmp)
		return tmp2
	}	else {
		return 0
	}
}

//转换字符串
func String(val interface{}) string {
	switch value := val.(type) {
	case types.Nil:
		return ""
	case string:
		return strings.TrimSpace(value)
	case int:
		return strconv.Itoa(value)
	case int64:
		return strconv.FormatInt(value, 10)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case interface{}:
		return fmt.Sprintf("%s", value)
	default:
		return ""

	}

}

func String_Json2Arr(content interface{}) []string {
	if len(String(content)) > 0 {
		var tmpbasic []string
		err := json.Unmarshal([]byte(String(content)), &tmpbasic)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return tmpbasic
	}
	return []string{}
}

// HTML2str returns escaping text convert from html.
func HTML2str(html string) string {

	re, _ := regexp.Compile(`\<[\S\s]+?\>`)
	html = re.ReplaceAllStringFunc(html, strings.ToLower)

	//remove STYLE
	re, _ = regexp.Compile(`\<style[\S\s]+?\</style\>`)
	html = re.ReplaceAllString(html, "")

	//remove SCRIPT
	re, _ = regexp.Compile(`\<script[\S\s]+?\</script\>`)
	html = re.ReplaceAllString(html, "")

	re, _ = regexp.Compile(`\<[\S\s]+?\>`)
	html = re.ReplaceAllString(html, "\n")

	re, _ = regexp.Compile(`\s{2,}`)
	html = re.ReplaceAllString(html, "\n")

	html = Htmlunquote(html)
	return strings.TrimSpace(html)
}

// Htmlquote returns quoted html string.
func Htmlquote(text string) string {
	//HTML编码为实体符号
	/*
	   Encodes `text` for raw use in HTML.
	       >>> htmlquote("<'&\\">")
	       '&lt;&#39;&amp;&quot;&gt;'
	*/

	text = html.EscapeString(text)
	text = strings.NewReplacer(
		`“`, "&ldquo;",
		`”`, "&rdquo;",
		` `, "&nbsp;",
	).Replace(text)

	return strings.TrimSpace(text)
}

// Htmlunquote returns unquoted html string.
func Htmlunquote(text string) string {
	//实体符号解释为HTML
	/*
	   Decodes `text` that's HTML quoted.
	       >>> htmlunquote('&lt;&#39;&amp;&quot;&gt;')
	       '<\\'&">'
	*/

	text = html.UnescapeString(text)

	return strings.TrimSpace(text)
}

func String_isNumberic(text string) bool {
	pattern := "^[0-9]*$" //反斜杠要转义
	result, _ := regexp.MatchString(pattern, text)
	return result
}

//取得json字符串，并解析成map[string]interface{}
func String_GetJSON(s string) (map[string]interface{}, error) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}
	var js map[string]interface{}
	err := json.Unmarshal([]byte(s), &js)
	return js, err
}

//不能判断一定是，可以判断一定不是。判断方式，base64只包含特定字符;解码再转码，查验是否相等。目前貌似没有能一定判断是的方法，有的话请指正，感谢。
func String_IsBase64(str string) bool {
	pattern := "^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{4}|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)$"
	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		return false
	}
	if !(len(str)%4 == 0 && matched) {
		return false
	}
	unCodeStr, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return false
	}
	tranStr := base64.StdEncoding.EncodeToString(unCodeStr)
	//return str==base64.StdEncoding.EncodeToString(unCodeStr)
	if str == tranStr {
		return true
	}
	return false
}

//解析body_form值
func String_AnalysisURI(content string) (map[string]interface{},error){
	tmpArr := strings.Split(content, "&")
	result:= map[string]interface{}{}
	for _, value := range tmpArr {
		tmpArr2 := strings.Split(value, "=")
		if len(tmpArr2) == 2 {
			result[tmpArr2[0]],_ =url.QueryUnescape( tmpArr2[1])
		} else {
			return nil, errors.New("参数解析失败")
		}
	}
	return result,nil
}