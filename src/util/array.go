package util

import (
	"sort"
)

//判断数组是否包含某个元素
func Array_Container(array []interface{}, item interface{}) bool {
	for _, value := range array {
		if value == item {
			return true
		}
	}

	return false
}
//数组排序
func Array_SortInt64(array []int64) []int64 {
	if len(array)==0 {
		return []int64{}
	}
	tmpArr1:=[]int{}
	for _, value := range array {
	tmpArr1=	append(tmpArr1,String_GetInt(value))
	}
	sort.Ints(tmpArr1)

	tmpArr2:=[]int64{}
	for _,value:=range tmpArr1 {
		tmpArr2=append(tmpArr2,String_GetInt64(value))
	}
	return tmpArr2
}


func Array_Reverse(s []interface{}) []interface{} {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
func Array_ReverseInt64(s []int64) []int64 {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
func Array_ReverseMap(s []map[string]interface{}) []map[string]interface{} {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func Array_Remove(s []map[string]interface{},index int64) []map[string]interface{} {
	if len(s)==0 {
		return  s
	}

	if index==0{
		return s[1:]
	}
	if index==int64(len(s)-1) {
		return s[0:len(s)-1]
	}
	tmp:=s[0:index-1]
	return append(tmp, s[index+1:]...)
}

func Array_2Map(data []map[string]interface{},key string) map[string]map[string]interface{} {
	result:= map[string]map[string]interface{}{}
	for _,value:=range data  {
		if val,ok:=value[key];ok{
			result[String( val)]=value
		}
	}
	return result
}
