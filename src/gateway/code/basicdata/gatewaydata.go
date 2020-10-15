/**
网关运行需要的数据，对原来数据进行加锁复制，确保无读写冲突。
*/
package basicdata

import (
	"util/deepcopy"
)


//取得网关运行需要元数据
func GetGatewayMetaData() *GatewayMetadata {
	GWMetaList.RLock()
	defer GWMetaList.RUnlock()
	gm := deepcopy.Copy(GWMetaList).(GatewayMetadata)
	return &gm
}

//取得网关运行需要配置数据
func GetGatewaySettingConfig() *Syssettingconfig {
	Config.RLock()
	defer Config.RUnlock()
	config := deepcopy.Copy(Config).(Syssettingconfig)
	return &config
}
