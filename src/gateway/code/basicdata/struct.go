/**
开发人员：陈朝能
*/
package basicdata

type Status int64

//应用程序客户端
type Application struct {
	PKID               int64
	Title              string
	ApplicationClients []ApplicationClient
	//0：禁止访问，1：未锁定，2：锁定
	Status Status
	Path   string
}

//应用程序客户端
type ApplicationClient struct {
	PKID          int64
	ApplicationId int64
	Application   Application
	Title         string
	Client        string
	//0：禁止访问，1：未锁定，2：锁定
	Status   Status
	Agentkey string
	//启用验签
	EnableSign int64
	//验签方式
	Encryption string
	//验签盐值
	Salt string
	//启用BASE64加密
	EnableBASE64 int64
}

//微服务
type Module struct {
	PKID   int64
	Title  string
	Status Status
}

//微服务API
type ModuleAPI struct {
	PKID                 int64
	Title                string
	Status               Status
	Moduleid             int64
	URL                  string
	RequestMothed        string
	Urlparams            string
	Businessparams       string
	ReturnErrorContent   string
	ReturnTimeoutContent string
	Timeout              int64
}

//微服务集群服务器
type ModuleServer struct {
	PKID           int64
	Title          string
	Status         Status
	Moduleid       int64
	Host           string
}

//OpenAPI
type OpenAPI struct {
	PKID              int64
	Title             string
	Status            Status
	ContainAllApp     bool
	CType             string
	URL               string
	Method            string
	DataTranslateType string
	Timeout           int64

	//Applications []Application
	Nodes   map[string][]OpenAPINode
	Results map[string][]OpenAPIResult
}

//OpenAPI节点
type OpenAPINode struct {
	PKID  int64
	Title string
	Tag   string
	//JoinreQuestParams    string
	Urlparams            string
	Method               string
	DataTranslateType    string
	Businessparams       string
	IsMustReturn         bool
	ReturnErrorContent   string
	ReturnTimeoutContent string
	Timeout              int64
	SortNo               int64

	ApplicationID int64
	Module    Module
	 Moduleapi ModuleAPI
	Servers   []ModuleServer
}

//OpenAPI结果转换
type OpenAPIResult struct {
	PKID        int64
	Method      string
	Newpropname string
	Controltype string
	Oldpropname string
	Sortno      int64
}

//适配协议
type ProtocolURL struct {
	PKID      int64
	Title     string
	Status    Status
	SourceURL string
	TargetURL string
	Timeout   int64

	ContainAllApp bool
	Applications  []Application
	Module        Module
	Servers       []ModuleServer
}

//GET适配协议
type ProtocolURLGET struct {
	PKID          int64
	Title         string
	Status        Status
	ContainAllApp bool
	SourceURL     string
	TargetURL     string
}

//上传配置
type UploadURL struct {
	PKID          int64
	Title         string
	Status        Status
	Valid         bool
	ContainAllApp bool
	SourceURL     string
	TargetURL     string
	Timeout       int64

	MaxFileSize int64
	FileFormat  string
	UpNum int64

	ReturnErrorContent   string
	ReturnSizeContent    string
	ReturnFormatCotnent  string
	ReturnTimeoutContent string
	Module               Module
	Servers              []ModuleServer
}
