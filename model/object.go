package model

import "time"

//基础类
type Base struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

//固件
type Firmware22 struct {
	Base
	TaskUID             string `json:"task_uid" gorm:"comment:'状态机任务id'"`
	FirmwareTask22Id    uint   `gorm:"index" json:"firmware_task22_id"`
	Status              string `gorm:"index" json:"status" gorm:"comment:'状态'"`
	RequestFilePath     string `json:"request_file_path" gorm:"comment:'提交文件路径'"`
	CompressedFilePath  string `json:"compressed_file_path" gorm:"comment:'压缩包文件路径'"`
	MediaAddress        string `json:"media_address" gorm:"comment:'media文件路径前缀'"`
	Size                string `json:"size" gorm:"comment:'固件大小'"`
	OperatingSystem     string `json:"operating_system" gorm:"comment:'操作系统'"`
	WordFile            string `json:"word_file" gorm:"comment:'word文件地址'"`
	PdfFile             string `json:"pdf_file" gorm:"comment:'pdf文件地址'"`
	Name                string `json:"name" gorm:"comment:'固件名称'"`
	Software22Id        uint   `json:"software22_id"`
	RequestId           string `gorm:"index" json:"request_id"`     //请求id
	RequestIdInt        int    `gorm:"index" json:"request_id_int"` //请求id
	TemplateUid         string `gorm:"index" json:"template_uid" gorm:"comment:'模板UID'"`
	UserNameId          uint   `gorm:"index" json:"user_name_id"`
	Ability             string `json:"abilityability" gorm:"type:longtext"`
	Total               int    `json:"total"`
	Finished            int    `json:"finished"`
	LoopholesTotal      int    `json:"loopholes_total"`
	TeamPosition        uint   `gorm:"index" json:"team_position"` //在队列中的位置
	CreatedBy           uint   `gorm:"index" json:"created_by"`
	ErrorInfo           string `gorm:"index" json:"error_info"` //错误信息
	Flag                string `gorm:"index" json:"flag"`
	Recording           int    `json:"recording"`
	DecompiledRequestId string `json:"decompiled_request_id"`
	DecompiledFile      string `json:"decompiled_file"`
	ActualTotal         int    `json:"actual_total"` //真实的任务总数
	DecompilationState  string `json:"decompilation_state"`
	TimeDifference      string `json:"time_difference"`
}

//固件任务
type FirmwareTask22 struct {
	Base
	MissionName   string `gorm:"index" json:"mission_name" gorm:"comment:'任务名称'"`
	EquipmentType string `gorm:"index" json:"equipment_type" gorm:"comment:'设备类型'"`
	TradeNames    string `gorm:"index" json:"trade_names" gorm:"comment:'厂商名称'"`
	DeviceName    string `gorm:"index" json:"device_name" gorm:"comment:'设备名称'"`
	//DeviceModel    string `gorm:"index" json:"device_model" gorm:"comment:'设备型号'"`
	TemplateName   string `gorm:"index" json:"template_name" gorm:"comment:'模板名称'"`
	TemplateUid    string `gorm:"index" json:"template_uid" gorm:"comment:'模板UID'"`
	Status         string `gorm:"index" json:"status" gorm:"comment:'状态'"`
	ClientName     string `json:"user" gorm:"size:50;index:user_idx"`             //客户名称
	User           string `json:"client_name" gorm:"size:50;index:user_idx"`      //客户名称
	UserAccount    string `json:"user_account" gorm:"size:50;index:account_idx" ` //账号
	UserRealName   string `json:"user_real_name"`                                 //用户名
	ClientKey      string `gorm:"index" json:"client_key"`
	TaskUID        string `gorm:"index" json:"task_uid" gorm:"comment:'状态机任务id'"`
	LoopholesTotal int    `json:"loopholes_total"`
	ErrorInfo      string `json:"error_info"` //错误信息
}

// 软件成分
type Software22 struct {
	Base
	Firmware22Id      uint   `gorm:"index" json:"firmware22_id"`
	MD5               string `json:"md_5" gorm:"comment:'MD5'"`
	SHA256            string `json:"sha_256" gorm:"comment:'SHA256';size:64"`
	SHA512            string `json:"sha_512" gorm:"comment:'SHA512';size:128"`
	FilePath          string `json:"file_path" gorm:"comment:'文件路径';type:longtext"`
	FileName          string `json:"file_name" gorm:"comment:'文件名称'"`
	FileType          string `json:"file_type" gorm:"comment:'文件类型';type:longtext"`
	FileSystem        string `json:"file_system" gorm:"comment:'文件系统';type:longtext"`
	CpuArchitecture   string `json:"cpu_architecture" gorm:"comment:'CPU架构';type:longtext"`
	SoftwareName      string `gorm:"index" json:"software_name" gorm:"comment:'软件名'"`
	SoftwareVersion   string `gorm:"index" json:"software_version" gorm:"comment:'软件版本'"`
	VimChecker        string `gorm:"index" json:"vim_checker" gorm:"comment:'是否是vim文件';"`
	TmpfileChecker    string `gorm:"index" json:"tmpfile_checker" gorm:"comment:'是否是临时文件';"`
	BakfileChecker    string `gorm:"index" json:"bakfile_checker" gorm:"comment:'是否是备份文件';"`
	SourcecodeLeakage string `gorm:"index" json:"sourcecode_leakage"` //泄露源码的类型
	Disabled          bool   `json:"disabled"`                        //是不是已经反编译了
	UniqueSting       string `json:"unique_sting" gorm:"comment:'唯一索引';unique_index:fileId"`
}

//固件相关通用配置
type FirmwareGeneralConfiguration struct {
	Base
	Key             string `json:"key"`                               //关键词
	PKey            string `json:"p_key"`                             //上级关键词
	Data            string `json:"data"`                              //数据
	Type            string `json:"type"`                              //类型
	IsBool          bool   `json:"is_bool"`                           //
	Remarks         string `json:"remarks"`                           //备注
	BackupData      string `json:"backup_data"  gorm:"type:longtext"` //次要数据
	CallbackTrigger string `json:"callback_trigger"`                  //触发器
}

//漏洞发生统计
type VulnerabilityStatistics struct {
	Base
	Firmware22Id        uint `gorm:"index" json:"firmware22_id"`
	CweChecker          bool `gorm:"index" json:"cwe_checker"`
	SSH                 bool `gorm:"index" json:"ssh"`
	FTP                 bool `gorm:"index" json:"ftp"`
	CryptoMaterial      bool `gorm:"index" json:"crypto_material"`
	CveLookup           bool `gorm:"index" json:"cve_lookup"`
	ExploitMitigations  bool `gorm:"index" json:"exploit_mitigations"`
	Hardcode            bool `gorm:"index" json:"hardcode"`
	UsersPasswords      bool `gorm:"index" json:"users_passwords"`
	MalwareScanner      bool `gorm:"index" json:"malware_scanner"`
	BakfileChecker      bool `gorm:"index" json:"bakfile_checker"`
	GitChecker          bool `gorm:"index" json:"git_checker"`
	SelfstartingChecker bool `gorm:"index" json:"selfstarting_checker"`
	SourcecodeLeakage   bool `gorm:"index" json:"sourcecode_leakage"`
	SvnChecker          bool `gorm:"index" json:"svn_checker"`
	TmpfileChecker      bool `gorm:"index" json:"tmpfile_checker"`
	VimChecker          bool `gorm:"index" json:"vim_checker"`
	WpahardcodeChecker  bool `gorm:"index" json:"wpahardcode_checker"`
	RegistryChecker     bool `gorm:"index" json:"registry_checker"`
}
//CVE漏洞检测
type CveLookup22 struct {
	Base
	Software22Id uint   `gorm:"index" json:"software22_id"`
	CveId        string `gorm:"index" json:"cve_id"`
	Name         string `gorm:"index" json:"name"`
	CnnvdId      string `json:"cnnvd_id"`
	CnvdId       string `json:"cnvd_id"`
	HazardGrade  string `json:"hazard_grade"`
	LoopholeType string `json:"loophole_type"`
	Introduction string `json:"introduction" gorm:"type:longtext"`
	Notice       string `json:"notice" gorm:"type:longtext"`
}

//漏洞
type FirmwareDetectionVulnerability struct {
	Base
	LId                      string `gorm:"index" json:"l_id"`                              //漏洞编号
	TestItems                string `json:"test_items"`                                     //检测项目
	TestPurpose              string `json:"test_purpose"`                                   //检测目的
	RiskLevel                string `json:"risk_level"`                                     //风险等级
	VulnerabilityDescription string `json:"vulnerability_description" gorm:"type:longtext"` //漏洞描述
	RepairSuggestions        string `json:"repair_suggestions" gorm:"type:longtext"`        //修复建议
}

//cwe漏洞
type CweChecker22 struct {
	Base
	Software22Id  uint   `gorm:"index" json:"software22_id" gorm:"index"`
	CweId         string `gorm:"index" json:"cwe_id"`
	PluginVersion string `json:"plugin_version"`
	Warnings      string `json:"warnings" gorm:"type:longtext"`
	WS            string `json:"ws" gorm:"type:longtext"`
}
//明文用户密码检测
type UsersPasswords22 struct {
	Base
	Software22Id uint   `gorm:"index" json:"software22_id"`
	Name         string `json:"name"`
	PasswordHash string `json:"password_hash"`
	Password     string `json:"password"`
	Cracked      string `json:"cracked"`
}
//证书、密钥文件泄露检测
type CryptoMaterial22 struct {
	Base
	Software22Id uint   `gorm:"index" json:"software22_id"`
	Key          string `json:"key"`
	Val          string `json:"val" gorm:"type:longtext"`
	Path         string `json:"path" gorm:"type:longtext"`
}
//WPA密码硬编码检测
type WpahardcodeChecker22 struct {
	Base
	Software22Id   uint   `gorm:"index" json:"software22_id"`
	Description    string `json:"description"`
	RelatedStrings string `json:"related_strings" gorm:"type:longtext"`
}
type SvnChecker22 struct {
	Base
	Software22Id uint   `gorm:"index" json:"software22_id"`
	RepoUrl      string `json:"repo_url"`                         //远程仓库地址字符串
	RepoUUID     string `json:"repo_uuid"`                        //远程仓库的UUID
	Credentials  string `json:"credentials" gorm:"type:longtext"` //开发者的⽤户名和密钥信息
}
type GitChecker22 struct {
	Base
	Software22Id        uint   `gorm:"index" json:"software22_id"`
	RepoUrl             string `json:"repo_url"`             //远程仓库地址
	DeveloperEmail      string `json:"developer_email"`      //开发者邮箱
	DeveloperName       string `json:"developer_name"`       //开发者姓名
	DeveloperCredential string `json:"developer_credential"` //开发者的⽤户名以及密钥信息
}
//恶意软件
type MalwareScanner22 struct {
	Base
	Software22Id uint   `gorm:"index" json:"software22_id"`
	VirusName    string `json:"virus_name"`
	VirusPath    string `json:"virus_path"`
}
//⼆进制缓解措施检测
type ExploitMitigations22 struct {
	Base
	Software22Id        uint   `gorm:"index" json:"software22_id"`
	NXStatus            string `gorm:"index" json:"nx_status"`
	NXVal               string `json:"nx_val" gorm:"type:longtext"`
	RelroStatus         string `gorm:"index" json:"relro_status"`
	RelroVal            string `json:"relro_val" gorm:"type:longtext"`
	CanaryStatus        string `gorm:"index" json:"canary_status"`
	CanaryVal           string `json:"canary_val" gorm:"type:longtext"`
	FortifySourceStatus string `gorm:"index" json:"fortify_source_status"`
	FortifySourceVal    string `json:"fortify_source_val" gorm:"type:longtext"`
	PIEStatus           string `gorm:"index" json:"pie_status"`
	PIEVal              string `json:"pie_val" gorm:"type:longtext"`
}
//不安全的配置检测
type Configset22 struct {
	Base
	Software22Id uint   `gorm:"index" json:"software22_id"`
	Ssh          string `gorm:"index" json:"ssh" gorm:"type:longtext"`
	Ftp          string `gorm:"index" json:"ftp" gorm:"type:longtext"`
	Hardcode     string `gorm:"index" json:"hardcode" gorm:"type:longtext"`
}
//自启动服务检测
type SelfstartingChecker22 struct {
	Base
	Software22Id uint   `gorm:"index" json:"software22_id"`
	InitType     string `json:"init_type"`
	Script       string `json:"script" gorm:"type:longtext"`
	Inittab      string `json:"inittab" gorm:"type:longtext"`
	ExecStart    string `json:"exec_start" gorm:"type:longtext"`
	Exec         string `json:"exec" gorm:"type:longtext"`
	Description  string `json:"description" gorm:"type:longtext"`
	PreStart     string `json:"pre_start"`
}
//注册表文件检测
type RegistryChecker22 struct {
	Base
	Software22Id uint   `gorm:"index" json:"software22_id"`
	Key          string `json:"key"`
	Val          string `json:"val" gorm:"type:longtext"`
}