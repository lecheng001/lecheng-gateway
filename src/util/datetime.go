package util

import (
	"time"
)

func Date_FormatYmdHis(value interface{}) string {
	var tmp int64 = 0
	switch value.(type) {
	case int, int32, int64:
		tmp = value.(int64)
	case string:
		tmp = String_GetInt64(value)
	default:
		tmp = 0
	}
	return time.Unix(tmp, 0).Format("2006-01-02 15:04:05")
}
func Date_FormatYmd(value interface{}) int64 {
	var tmp int64 = 0
	switch value.(type) {
	case int, int32, int64:
		tmp = value.(int64)
	case string:
		tmp = String_GetInt64(value)
	default:
		tmp = 0
	}
	if tmp == 0 {
		tmp = time.Now().Unix()
	}
	return String_GetInt64(time.Unix(tmp, 0).Format("20060102"))
}
func Date_FormatYmd2(value interface{}) string {
	var tmp int64 = 0
	switch value.(type) {
	case int, int32, int64:
		tmp = value.(int64)
	case string:
		tmp = String_GetInt64(value)
	default:
		tmp = 0
	}
	return time.Unix(tmp, 0).Format("2006-01-02")
}

func Date_Formatmd(value interface{}) string {
	LogDebug(value)
	var tmp int64 = 0
	switch value.(type) {
	case int, int32, int64:
		tmp = value.(int64)
	case string:
		tmp = String_GetInt64(value)
	default:
		tmp = 0
	}
	//logs.Debug(tmp)
	return time.Unix(tmp, 0).Format("01-02")
}

func Date_Format(value interface{}, format string) string {
	//logs.Debug("FormatDate",value)
	var tmp int64 = 0
	switch value.(type) {
	case int, int32, int64:
		tmp = value.(int64)
	case string:
		tmp = String_GetInt64(value)
	default:
		tmp = 0
	}
	//logs.Debug(tmp)
	return time.Unix(tmp, 0).Format(format)
}

func Date_Str2Unix(value interface{}) int64 {
	toBeCharge := String(value)
	timeLayout := "2006-01-02 15:04:05"                             //转化所需模板
	loc, _ := time.LoadLocation("Local")                            //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在对应时区转化为time.time类型
	return theTime.Unix()
}
