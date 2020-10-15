package jsondb

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
	"util"
	"util/deepcopy"
)

func Query(jsonfile string, where map[string]interface{}, order map[string]string) ([]map[string]interface{}, error) {

	content, err := ioutil.ReadFile(util.IO_GetRootPath() + "/" + jsonfile)
	if err != nil {
		util.LogDebug(err)
		return nil, err
	}
	tmp := []map[string]interface{}{}
	err = json.Unmarshal(content, &tmp)
	if err != nil {
		util.LogDebug(err)
		return nil, err
	}

	if where != nil && len(where) > 0 && len(tmp) > 0 {
		hasrun := false
		tmplist := []map[string]interface{}{}
		for key, value := range where {
			if hasrun {
				tmp = deepcopy.Copy(tmplist).([]map[string]interface{})
				tmplist = []map[string]interface{}{}
			}
			hasrun = true
			for _, item := range tmp {
				tmpPara := strings.Split(key, "|")
				if len(tmpPara) == 1 {
					//where += " and " + key + "=? "
					if _, ok := value.(string); ok {
						if util.String(item[key]) == util.String(value) {
							tmplist = append(tmplist, item)
						}
					}
					if _, ok := value.(int); ok {
						if util.String_GetInt64(item[key]) == util.String_GetInt64(value) {
							tmplist = append(tmplist, item)
						}
					}

				} else if len(tmpPara) == 2 {
					//	type := "";
					switch tmpPara[1] {
					case ">=":
						if util.String_GetInt64(item[tmpPara[0]]) >= util.String_GetInt64(value) {
							tmplist = append(tmplist, item)
						}
					//like
					case "like":
						if strings.Contains(item[tmpPara[0]].(string), value.(string)) {
							tmplist = append(tmplist, item)
						}

					default:
						return nil, errors.New("未知类型数据！")
					}
				}
			}
		}
		tmp = tmplist
	}

	//排序
	for key, value := range order {
		mapsort := util.MapsSort{key, tmp}
		if value == "desc" {
			tmp = mapsort.Sort()
		} else if value == "asc" {
			tmp = mapsort.Reverse()
		}
	}
	return tmp, nil
}

func QueryRow(jsonfile string, where map[string]interface{}, order map[string]string) map[string]interface{} {

	list, err := Query(jsonfile, where, nil)
	if err != nil {
		return nil
	}
	if len(list) > 0 {
		return list[0]
	} else {
		return nil
	}
}

func GetInfo(jsonfile string, pkid int64) map[string]interface{} {
	filePath := util.IO_GetRootPath() + "/" + jsonfile
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		util.LogDebug(err)
		return nil
	}
	tmp := []map[string]interface{}{}
	err = json.Unmarshal(content, &tmp)
	if err != nil {
		util.LogDebug(err)
		return nil
	}

	for _, value := range tmp {
		if util.String_GetInt64(value["pkid"]) == pkid {
			return value
		}
	}

	return nil
}

/**
执行SQL语句
*/
func Delete(jsonfile, IdKey string, id int64) (int64, error) {
	filePath := util.IO_GetRootPath() + "/" + jsonfile
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		util.LogDebug(err)
		return 0, err
	}
	tmp := []map[string]interface{}{}
	err = json.Unmarshal(content, &tmp)
	if err != nil {
		util.LogDebug(err)
		return 0, err
	}

	for key, value := range tmp {
		if util.String_GetInt64(value[IdKey]) == id {
			tmp = util.Array_Remove(tmp, int64(key))
			break
		}
	}

	result, err := json.Marshal(tmp)
	if err != nil {
		return 0, err
	}

	err = ioutil.WriteFile(filePath, result, 777)
	if err != nil {
		return 0, err
	}
	return 1, nil
	return 0, nil
}

func InsertMap(jsonfile string, colu map[string]interface{}) (int64, error) {
	filePath := util.IO_GetRootPath() + "/" + jsonfile
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		util.LogDebug(err)
		return 0, err
	}
	tmp := []map[string]interface{}{}
	err = json.Unmarshal(content, &tmp)
	if err != nil {
		util.LogDebug(err)
		return 0, err
	}

	//取得最大ID
	maxid := int64(0)
	for _, value := range tmp {
		id := util.String_GetInt64(value["pkid"])
		if id > maxid {
			maxid = id
		}
	}
	newmaxid := maxid + 1
	colu["pkid"] = newmaxid
	tmp = append(tmp, colu)
	result, err := json.Marshal(tmp)
	if err != nil {
		return 0, err
	}

	err = ioutil.WriteFile(filePath, result, 777)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

/**
* 更新操作
* @param type pColumn 列数据，array('title'=>title,'url'=>url);
 */
func UpdateByID(jsonfile string, id int64, colu map[string]interface{}) (int64, error) {
	filePath := util.IO_GetRootPath() + "/" + jsonfile
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		util.LogDebug(err)
		return 0, err
	}
	tmp := []map[string]interface{}{}
	err = json.Unmarshal(content, &tmp)
	if err != nil {
		util.LogDebug(err)
		return 0, err
	}

	for key, value := range tmp {
		if util.String_GetInt64(value["pkid"]) == id {
			for key2, value2 := range colu {
				tmp[key][key2] = value2
			}
		}
	}

	result, err := json.Marshal(tmp)
	if err != nil {
		return 0, err
	}

	err = ioutil.WriteFile(filePath, result, 777)
	if err != nil {
		return 0, err
	}
	return 1, nil
}
func UpdateMap(jsonfile string, colu map[string]interface{}, where map[string]interface{}) (int64, error) {
	panic("未实现")
}
