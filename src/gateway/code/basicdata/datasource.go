/**
网关基础数据
开发人员：陈朝能
*/
package basicdata

import (
	"util/jsondb"
	"util"
	"errors"
	"sync"
)

//元数据
type DataSource struct {
	sync.Mutex

	application       []map[string]interface{}
	applicationclient []map[string]interface{}

	module       []map[string]interface{}
	moduleserver []map[string]interface{}

	openapi []map[string]interface{}

	protocolurl            []map[string]interface{}

	protocolurlget []map[string]interface{}

	uploadlurl []map[string]interface{}
}

var DataDao DataSource

//从数据库载入全部网关数据
func LoadDataSource() {
	DataDao.Lock()
	defer DataDao.Unlock()

	err := errors.New("")
	DataDao.application, err = jsondb.Query("database/application.json", map[string]interface{}{"cstatus": 1}, map[string]string{"pkid":"desc"})
	if err != nil {
		util.LogError("database/application.json丢失或json文件无效", err)
		panic("database/application.json丢失或json文件无效")
		return
	}
	DataDao.applicationclient, err = jsondb.Query("database/applicationclient.json", map[string]interface{}{"cstatus": 1}, map[string]string{"pkid":"desc"})
	if err != nil {
		util.LogError("database/applicationclient.json丢失或json文件无效", err)
		panic("database/applicationclient.json丢失或json文件无效")
		return
	}
	DataDao.module, err = jsondb.Query("database/module.json", map[string]interface{}{"cstatus": 1}, map[string]string{"pkid":"desc"})
	if err != nil {
		util.LogError("database/module.json丢失或json文件无效", err)
		panic("database/module.json丢失或json文件无效")
		return
	}
	DataDao.moduleserver, err = jsondb.Query("database/moduleserver.json", map[string]interface{}{"cstatus": 1}, map[string]string{"pkid":"desc"})
	if err != nil {
		util.LogError("database/moduleserver.json丢失或json文件无效", err)
		panic("database/moduleserver.json丢失或json文件无效")
		return
	}

	DataDao.protocolurl, err = jsondb.Query("database/protocolurl.json", map[string]interface{}{"cstatus": 1}, map[string]string{"csort":"asc"})
	if err != nil {
		util.LogError("database/protocolurl.json丢失或json文件无效", err)
		panic("database/protocolurl.json丢失或json文件无效")
		return
	}
	DataDao.protocolurlget, err = jsondb.Query("database/protocolurlget.json", map[string]interface{}{"cstatus": 1},map[string]string{"csort":"asc"})
	if err != nil {
		util.LogError("database/protocolurlget.json丢失或json文件无效", err)
		panic("database/protocolurlget.json丢失或json文件无效")
		return
	}
	DataDao.uploadlurl, err = jsondb.Query("database/upload.json", map[string]interface{}{"cstatus": 1},map[string]string{"csort":"asc"})
	if err != nil {
		util.LogError("database/upload.json丢失或json文件无效", err)
		panic("database/upload.json丢失或json文件无效")
		return
	}

}
