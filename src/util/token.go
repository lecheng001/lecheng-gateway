package util

import (
	"fmt"
	"strings"
)

//取得TOKEN信息
func Token_Get(info map[string]string) string {
	if len(info) == 0 {
		return ""
	}

	content := ""

	for key, value := range info {
		content += fmt.Sprintf("%s=%s", key, value) + "&"
	}
	md5salt := ConfigGetString("token.md5salt")
	dessalt := ConfigGetString("token.dessalt")

	content += MD5(content + md5salt)
	return Encry_DES_ECB(content, dessalt)
}

//分析TOKEN数据
func Token_Analysis(info string) map[string]string{
	defer func() {
		if r := recover(); r != nil {
			LogError("recover a error ", r)
		}
	}()

	if len(info) == 0 {
		return nil
	}

	md5salt := ConfigGetString("token.md5salt")
	dessalt := ConfigGetString("token.dessalt")

	content := Decry_DES_ECB(info, dessalt)
	if len(content) == 0 {
		return nil
	}

	//fmt.Println(content)
	tmp := strings.Split(content, "&")
	//fmt.Println(tmp)
	if len(tmp) == 1 {
		return nil
	}
	tmpinfo := ""
	result := map[string]string{}
	for i := 0; i < len(tmp)-1; i++ {
		tmpinfo += tmp[i] + "&"
		tmpitem := strings.Split(tmp[i], "=")
		if len(tmpitem) == 2 {
			result[tmpitem[0]] = tmpitem[1]
		} else {
			return nil
		}
	}
	sign_old := tmp[len(tmp)-1]
	sign_new := MD5(tmpinfo + md5salt)
	if sign_new == sign_old {
		//logs.Debug(result)
		return result
	}
	LogDebug("TOKEN解析失败或不存在！", result)
	return nil
}

func Token_Checkvalid(token string) bool {

	return true
}
