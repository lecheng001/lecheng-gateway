package viewserver

import (
	"github.com/valyala/fasthttp"
	"strings"
)

func InitRouter(e *fasthttp.RequestCtx) {
	//fre := "/backapi/"
	//method:=string(e.Method())
	url :=strings.Trim( strings.Replace(string(e.Path()), "/html", "", 1),"/")
	cc := Controller{e}
	switch url {
	case "","home":
		if cc.IsGet() {
			HomeContrller{cc}.Index()
		}
	case "login":
		if cc.IsGet(){
			LoginController{cc}.LoginPage()
		}else if cc.IsPost(){
			LoginController{cc}.Login()
		}
	case "application/list":
		if cc.IsGet(){
			ApplicationController{cc}.Getlist()
		}
	case "application/detail":
		if cc.IsGet(){
			ApplicationController{cc}.Getinfo()
		}
	case "application/postdata":
		if cc.IsPost(){
			ApplicationController{cc}.Postdata()
		}
	case "application/dealitem":
		if cc.IsPost(){
			ApplicationController{cc}.DealStatus()
		}
	case "applicationclient/list":
		if cc.IsGet(){
			ApplicationClientController{cc}.Getlist()
		}
	case "applicationclient/detail":
		if cc.IsGet(){
			ApplicationClientController{cc}.Getinfo()
		}
	case "applicationclient/postdata":
		if cc.IsPost(){
			ApplicationClientController{cc}.Postdata()
		}
	case "applicationclient/dealitem":
		if cc.IsPost(){
			ApplicationClientController{cc}.DealStatus()
		}

	case "module/list":
		if cc.IsGet(){
			ModuleController{cc}.Getlist()
		}
	case "module/detail":
		if cc.IsGet(){
			ModuleController{cc}.Getinfo()
		}
	case "module/postdata":
		if cc.IsPost(){
			ModuleController{cc}.Postdata()
		}
	case "module/dealitem":
		if cc.IsPost(){
			ModuleController{cc}.DealStatus()
		}
	case "moduleserver/list":
		if cc.IsGet(){
			ModuleserverController{cc}.GetList()
		}
	case "moduleserver/detail":
		if cc.IsGet(){
			ModuleserverController{cc}.GetInfo()
		}
	case "moduleserver/postdata":
		if cc.IsPost(){
			ModuleserverController{cc}.PostData()
		}
	case "moduleserver/dealitem":
		if cc.IsPost(){
			ModuleserverController{cc}.DealStatus()
		}

	case "protocolurl/list":
		if cc.IsGet(){
			ProtocolUrlController{cc}.GetList()
		}
	case "protocolurl/detail":
		if cc.IsGet(){
			ProtocolUrlController{cc}.GetInfo()
		}
	case "protocolurl/postdata":
		if cc.IsPost(){
			ProtocolUrlController{cc}.PostData()
		}
	case "protocolurl/dealitem":
		if cc.IsPost(){
			ProtocolUrlController{cc}.DealStatus()
		}

	case "protocolurlget/list":
		if cc.IsGet(){
			ProtocolUrlGetController{cc}.GetList()
		}
	case "protocolurlget/detail":
		if cc.IsGet(){
			ProtocolUrlGetController{cc}.GetInfo()
		}
	case "protocolurlget/postdata":
		if cc.IsPost(){
			ProtocolUrlGetController{cc}.PostData()
		}
	case "protocolurlget/dealitem":
		if cc.IsPost(){
			ProtocolUrlGetController{cc}.DealStatus()
		}


	case "upload/list":
		if cc.IsGet(){
			UploadController{cc}.GetList()
		}
	case "upload/detail":
		if cc.IsGet(){
			UploadController{cc}.GetInfo()
		}
	case "upload/postdata":
		if cc.IsPost(){
			UploadController{cc}.PostData()
		}
	case "upload/dealitem":
		if cc.IsPost(){
			UploadController{cc}.DealStatus()
		}

	case "systemsetting":
		if cc.IsGet() {
			SystemSettingController{cc}.GetList()
		}
	case "systemsetting/write":
		if cc.IsPost() {
			SystemSettingController{cc}.ConfigWrite()
		}

	case "gatewaymanage":
		if cc.IsGet() {
			GatewayManageController{cc}.GetList()
		}
	case "gatewaymanage/refresh":
		if cc.IsPost() {
			GatewayManageController{cc}.Refresh()
		}
	case "gatewaymanage/certificate":
		if cc.IsPost() {
			GatewayManageController{cc}.Certificate()
		}
	case "gatewaymanage/refreshconfig":
		if cc.IsPost() {
			GatewayManageController{cc}.RefreshConfig()
		}
	case "gatewaymanage/changedebug":
		if cc.IsPost() {
			GatewayManageController{cc}.ChangeDebug()
		}
	case "testurl":
		if cc.IsGet() {
			TestUrlController{cc}.TestUrl()
		}
	case "testurl/list":
		if cc.IsGet() {
			TestUrlController{cc}.GetList()
		}
	case "testurl/postdata":
		if cc.IsPost() {
			TestUrlController{cc}.PostData()
		}
	case "testurl/getsign":
		if cc.IsPost() { 
			TestUrlController{cc}.RsaCode()
		}
	case "testurl/geturlcontent":
		if cc.IsPost() {
			TestUrlController{cc}.TestGetUrlContent()
		}
	case "testurl/delete":
		if cc.IsPost() {
			TestUrlController{cc}.Delete()
		}
	}

}
