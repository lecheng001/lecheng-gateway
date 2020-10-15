package util

import (
	"util/simplejson"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

var jsonData *simplejson.Json
var _mutex sync.Mutex

func JSON_RefreshConfig() (*simplejson.Json, error) {
	return getALlConfig(true)
}

func _getALlConfig(filename string, refresh bool) (*simplejson.Json, error) {
	file := IO_GetRootPath() + "/" + filename
	if _, err := os.Stat(file); err != nil {
		return &simplejson.Json{}, nil
	}
	//fmt.Println("read file:" + file)
	js, err := simplejson.NewJsonFromFile(file)
	if err != nil {
		return &simplejson.Json{}, err
	}

	return js, nil

}

//取得config配置文件的值
func JSON_Config(filename string, keys ...string) (*simplejson.Json, error) {
	js, err := _getALlConfig(filename, false)
	if err != nil {
		return &simplejson.Json{}, err
	}
	//tmparr := strings.Split(key, ".")
	for _, v := range keys {
		js = js.Get(v)
	}

	return js, nil
}

//取得config配置文件的值
func JSON_ConfigGetSliceMap(js *simplejson.Json, key string) []map[string]interface{} {
	tmparr := strings.Split(key, ".")
	for _, v := range tmparr {
		js = js.Get(v)
	}
	result := js.MustArray()
	result2 := []map[string]interface{}{}
	for _, value := range result {
		result2 = append(result2, value.(map[string]interface{}))
	}
	return result2
}

//config文件修改
func JSON_ConfigSet(filename string, key string, value interface{}) error {
	file := IO_GetRootPath() + "/" + filename
	tmpContent := []byte{}

	if len(key) > 0 {
		js, err := simplejson.NewJsonFromFile(file)
		if err != nil {
			return err
		}

		if !strings.Contains(key, ".") {
			js.Set(key, value)
		} else {
			js.SetPath(strings.Split(key, "."), value)
		}

		tmpContent, err = js.EncodePretty()
		if err != nil {
			return err
		}
	}else {
		tmpContent1, err :=json.MarshalIndent(&value, "", "  ")
		if err!=nil{
			return err
		}
		tmpContent=tmpContent1
	}
	_mutex.Lock()
	defer _mutex.Unlock()
	ioutil.WriteFile(file, tmpContent, 0666)
	return nil
}

//删除元素
func JSON_ConfigDelete(filename string, key string) error {
	_mutex.Lock()
	defer _mutex.Unlock()

	file := IO_GetRootPath() + "/" + filename
	js, err := simplejson.NewJsonFromFile(file)
	if err != nil {
		return err
	}

	//根据判断结果更新配置文件内容
	//更新config.json
	if !strings.Contains(key, ".") {
		js.Del(key)
	} else {
		tmparr := strings.Split(key, ".")
		tmpjs := js.GetPath(tmparr[:len(tmparr)-1]...)
		if tmpjs != nil {
			tmpjs.Del(tmparr[len(tmparr)-1])
		}
	}
	tmpContent, err := js.EncodePretty()
	if err != nil {
		return err
	}
	ioutil.WriteFile(file, tmpContent, 0666)

	return nil
}
