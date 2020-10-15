package util

import (
	"encoding/json"
	"sort"
	"strconv"
)

func Map_ToInterface(value interface{}) map[string]interface{} {
	if value == nil {
		return nil
	}

	tmp := map[string]interface{}{}
	if value, ok := value.(map[string]string); ok {
		for key, val := range value {
			tmp[key] = val
		}
	}
	return tmp
}

//将map对象转换成字符串
func Map_ToString(value interface{}) string {
	if value == nil {
		return ""
	}
	
	if value,ok:=value.(map[string]interface{});ok{
		lvC,err:= json.Marshal(value)
		if err!=nil{
			return ""
		}

		return string(lvC)
	}

	return ""
}

func Map_getValues(maplist []map[string]interface{}, key, keyvalue string) map[string]interface{} {
	for _, value := range maplist {
		if String(value[key]) == keyvalue {
			return value
		}
	}

	return nil
}



type MapsSort struct {
	Key string
	MapList []map[string] interface{}
}

func (m *MapsSort) Len() int {
	return len(m.MapList)
}

func (m *MapsSort) Less(i, j int) bool {
	var ivalue float64
	var jvalue float64
	var err error
	//fmt.Println(m.Key)
	switch m.MapList[i][m.Key].(type) {
	case string:
		ivalue,err = strconv.ParseFloat(m.MapList[i][m.Key].(string),64)
		if err != nil {
			LogError("map数组排序string转float失败：%v",err)
			return true
		}
	case int:
		ivalue = float64(m.MapList[i][m.Key].(int))
	case float64:
		ivalue = m.MapList[i][m.Key].(float64)
	case int64:
		ivalue = float64(m.MapList[i][m.Key].(int64))
	}
	switch m.MapList[j][m.Key].(type) {
	case string:
		jvalue,err = strconv.ParseFloat(m.MapList[j][m.Key].(string),64)
		if err != nil {
			LogError("map数组排序string转float失败：%v",err)
			return true
		}
	case int:
		jvalue = float64(m.MapList[j][m.Key].(int))
	case float64:
		jvalue = m.MapList[j][m.Key].(float64)
	case int64:
		jvalue = float64(m.MapList[j][m.Key].(int64))
	}
	return ivalue > jvalue
}

func (m *MapsSort) Swap(i, j int) {
	m.MapList[i],m.MapList[j] = m.MapList[j],m.MapList[i]
}

func (m *MapsSort) Sort()  []map[string]interface{} {
	sort.Sort(m)
	return m.MapList
}

func (m *MapsSort) Reverse()  []map[string]interface{} {
	sort.Sort(sort.Reverse(m))
	return m.MapList
}