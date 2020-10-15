package util

import (
	"util/simplejson"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
	"sync"
)

var ALLCONFIG *simplejson.Json

func Config_RefreshConfig() (*simplejson.Json, error) {
	return getALlConfig(true)
}

func getALlConfig(refresh bool) (*simplejson.Json, error) {
	if ALLCONFIG != nil && !refresh {
		return ALLCONFIG, nil
	}

	file := IO_GetRootPath() + "/config.json"
	//fmt.Println("read file:" + file)
	js, err := simplejson.NewJsonFromFile(file)
	if err != nil {
		return &simplejson.Json{}, nil
	}
	if env, _ := js.Get("env").String(); env != "pro" {
		devfile := IO_GetRootPath() + "/config_" + env + ".json"
		if exist, _ := IO_PathExists(devfile); exist == true {
			conf, _ := ioutil.ReadFile(file)
			confdev, _ := ioutil.ReadFile(devfile)
			confmap := map[string]interface{}{}
			confdevmap := map[string]interface{}{}
			json.Unmarshal(conf, &confmap)
			json.Unmarshal(confdev, &confdevmap)

			//合并MAP
			allconfig := Merge(confmap, confdevmap)
			allconfigtmp, _ := json.Marshal(allconfig)
			//重新生成JSON
			js, err = simplejson.NewJson(allconfigtmp)

		}
	}

	mutex.Lock()
	ALLCONFIG = js
	mutex.Unlock()
	return js, err

}

//取得config配置文件的值
func Config(key string) (*simplejson.Json, error) {
	js, err := getALlConfig(false)
	if err != nil {
		return &simplejson.Json{}, err
	}
	tmparr := strings.Split(key, ".")
	for _, v := range tmparr {
		js = js.Get(v)
	}

	return js, nil
}

//取得config配置文件的值
func ConfigGetBool(key string) bool {
	js, err := getALlConfig(false)
	if err != nil {
		return false
	}
	tmparr := strings.Split(key, ".")
	for _, v := range tmparr {
		js = js.Get(v)
	}
	result, _ := js.Bool()
	return result
}

//取得config配置文件的值
func ConfigGetString(key string) string {
	js, err := getALlConfig(false)
	if err != nil {
		return ""
	}
	tmparr := strings.Split(key, ".")
	for _, v := range tmparr {
		js = js.Get(v)
	}
	result, _ := js.String()
	return result
}

//取得config配置文件的值
func ConfigGetSlice(key string) []string {
	js, err := getALlConfig(false)
	if err != nil {
		return nil
	}
	tmparr := strings.Split(key, ".")
	for _, v := range tmparr {
		js = js.Get(v)
	}
	result := js.MustStringArray()
	return result
}

//取得config配置文件的值
func ConfigGetSliceMap(key string) []map[string]interface{} {
	js, err := getALlConfig(false)
	if err != nil {
		return nil
	}
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

//取得config配置文件的值
func ConfigGetInt64(key string) int64 {
	js, err := getALlConfig(false)
	if err != nil {
		return 0
	}
	tmparr := strings.Split(key, ".")
	for _, v := range tmparr {
		js = js.Get(v)
	}
	result, _ := js.Int64()
	return result
}

//取得config配置文件的值
func ConfigGetInt(key string) int {
	js, err := getALlConfig(false)
	if err != nil {
		return 0
	}
	tmparr := strings.Split(key, ".")
	for _, v := range tmparr {
		js = js.Get(v)
	}
	result, _ := js.Int()
	return result
}

//取得config配置文件的值
func ConfigGetFloat(key string) float64 {
	js, err := getALlConfig(false)
	if err != nil {
		return 0
	}
	tmparr := strings.Split(key, ".")
	for _, v := range tmparr {
		js = js.Get(v)
	}
	result, _ := js.Float64()
	return result
}

var mutex sync.Mutex

//config文件修改
func ConfigSet(key string, value interface{}) error {
	if ALLCONFIG == nil {
		getALlConfig(false)
	}
	if ALLCONFIG == nil {
		return errors.New("配置数据读取失败")
	}
	//1：pro模式，更新config.json文件
	//2：dev模式，config_dev.json不存在 更新config.json文件
	//3：dev模式，config_dev.json存在，属性值不存在 更新config.json文件
	//4：dev模式，config_dev.json存在，属性值存在  更新config_dev.json文件
	configfile := IO_GetRootPath() + "/config.json"
	configdevfile := ""

	//是否修改dev配置文件
	//isdev := false
	ispro := true
	js, err := simplejson.NewJsonFromFile(configfile)
	if err != nil {
		return err
	}
	env, err := js.Get("env").String()
	//pro模式
	if env == "pro" {
		ispro = true
	} else if env != "pro" {
		configdevfile = IO_GetRootPath() + "/config_" + env + ".json"
		exist, _ := IO_PathExists(configdevfile)
		//文件不存在
		if exist == false {
			ispro = true
		} else {
			//文件存在
			js_dev, err2 := simplejson.NewJsonFromFile(configdevfile)
			if err2 != nil {
				ispro = true
			} else {
				if !strings.Contains(key, ".") {
					//属性值
					ispro = js_dev.Get(key).Interface() == nil
				} else if strings.Contains(key, ".") {
					//属性值
					ispro = js_dev.GetPath(strings.Split(key, ".")...).Interface() == nil
				}

			}
		}
	}

	//根据判断结果更新配置文件内容
	if ispro {
		if !strings.Contains(key, ".") {
			js.Set(key, value)
		} else {
			js.SetPath(strings.Split(key, "."), value)
		}

		tmpContent, err := js.EncodePretty()
		if err != nil {
			return err
		}
		mutex.Lock()
		ioutil.WriteFile(configfile, tmpContent, 0666)
		mutex.Unlock()
	} else {
		//dev环境config
		js_dev, _ := simplejson.NewJsonFromFile(configdevfile)
		if !strings.Contains(key, ".") {
			js_dev.Set(key, value)
		} else {
			js_dev.SetPath(strings.Split(key, "."), value)
		}
		tmpContent, err := js_dev.EncodePretty()
		if err != nil {
			return err
		}
		mutex.Lock()
		ioutil.WriteFile(configdevfile, tmpContent, 0666)
		mutex.Unlock()
	}

	getALlConfig(true)

	return nil
}

//删除元素
func ConfigDelete(key string) error {
	mutex.Lock()
	defer mutex.Unlock()

	if ALLCONFIG == nil {
		getALlConfig(false)
	}
	if ALLCONFIG == nil {
		return errors.New("配置数据读取失败")
	}
	//1：pro模式，更新config.json文件
	//2：dev模式，config_dev.json不存在 更新config.json文件
	//3：dev模式，config_dev.json存在，属性值不存在 更新config.json文件
	//4：dev模式，config_dev.json存在，属性值存在  更新config_dev.json文件
	configfile := IO_GetRootPath() + "/config.json"
	configdevfile := IO_GetRootPath() + "/config_dev.json"
	//fmt.Println("read file:" + file)
	//是否修改dev配置文件
	//isdev := false
	ispro := true
	js, err := simplejson.NewJsonFromFile(configfile)
	if err != nil {
		return err
	}
	env, err := js.Get("env").String()
	//pro模式
	if env != "dev" {
		ispro = true
	} else if env == "dev" {

		exist, _ := IO_PathExists(configdevfile)
		//文件不存在
		if exist == false {
			ispro = true
		} else {
			//文件存在
			js_dev, err2 := simplejson.NewJsonFromFile(configdevfile)
			if err2 != nil {
				ispro = true
			} else {
				if !strings.Contains(key, ".") {
					//属性值
					ispro = js_dev.Get(key).Interface() == nil
				} else if strings.Contains(key, ".") {
					//属性值
					ispro = js_dev.GetPath(strings.Split(key, ".")...).Interface() == nil
				}

			}
		}
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
	ioutil.WriteFile(configfile, tmpContent, 0666)

	if !ispro {
		js_dev, _ := simplejson.NewJsonFromFile(configdevfile)
		//更新config_dev.json
		if !strings.Contains(key, ".") {
			js_dev.Del(key)
		} else {
			tmparr := strings.Split(key, ".")
			tmpjs := js_dev.GetPath(tmparr[:len(tmparr)-1]...)
			if tmpjs != nil {
				tmpjs.Del(tmparr[len(tmparr)-1])
			}
		}
		tmpContent, err := js_dev.EncodePretty()
		if err != nil {
			return err
		}
		ioutil.WriteFile(configdevfile, tmpContent, 0666)
	}

	getALlConfig(true)

	return nil
}
