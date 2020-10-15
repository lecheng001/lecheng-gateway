package global

import (
	"gateway/code/proxy"
)

//版本号写死在代码里
var GW_VERISON = "2.2"

//判断软件版本包
func CheckSoftPackage() (proxy.GwPackage, error) {
	version := GW_VERISON
	lvPackage := proxy.GwPackage{Version:version}

	//发布到GITHUB上采用下面代码
	lvPackage.Userid=1
	lvPackage.Projectid=1
	lvPackage.ProjectName="网关手动编译版"
	lvPackage.Projecttype="shequ"
	proxy.GWPACKAGE = lvPackage
	return lvPackage, nil
	
	
	//GitHub代码删除以下代码

}
