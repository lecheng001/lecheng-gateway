package viewserver

import (
	"encoding/json"
	"gateway/code/basicdata"
)

type HomeContrller struct {
	Controller
}

func (t HomeContrller) Index()    {
	tmp := map[string]interface{}{"error": 0}
	basicdata.GWMetaList.RLock()
	defer basicdata.GWMetaList.RUnlock()
	aa,_ :=json.Marshal( basicdata.GWMetaList)
	tmp["info"]=string(aa)
	t.RenderView( "index.gohtml",tmp)
}
