package main

import (
	"example.com/m/v2/model"
	"example.com/m/v2/reportTemplateRelated"
	"example.com/m/v2/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FirmwareService struct {
	db                    *gorm.DB
	edb                   *gorm.DB
	CryptoMaterialControl map[string]string
}

type SoftwareAndCveLookup struct {
	SID1           uint                `json:"sid1"`
	SID2           uint                `json:"sid2"`
	CveLookups     []model.CveLookup22 `json:"cve_lookups"`
	CveHighLookups []model.CveLookup22 `json:"cve_high_lookups"`
	CveInLookups   []model.CveLookup22 `json:"cve_in_lookups"`
	CveLowLookups  []model.CveLookup22 `json:"cve_low_lookups"`
}
type CweCheckerAndInfo struct {
	Cwe  model.FirmwareDetectionVulnerability `json:"cwe"`
	Cwes []model.CweChecker22                 `json:"cwes"`
}
type LsConfigset struct {
	Software model.Software22 `json:"software"`
	Data     []string         `json:"data"`
}
type cnvd struct {
	CveId  string `json:"cve_id"`
	CnvdId string `json:"cnvd_id"`
}

const DSN = "root:www@admin@2020@tcp(172.16.38.11:33060)/ssp?charset=utf8&parseTime=True&loc=Local"
const DSNS = "root:www@admin@2020@tcp(172.16.38.11:33060)/firmware_audit_engine?charset=utf8&parseTime=True&loc=Local"

const DRIVER = "mysql"

func main() {
	r := gin.Default()
	a := FirmwareService{}
	var err error
	a.db, err = gorm.Open(DRIVER, DSN)
	if err != nil {
		panic(err)
	}
	a.edb, err = gorm.Open(DRIVER, DSNS)
	if err != nil {
		panic(err)
	}
	a.db.SingularTable(true)
	a.edb.SingularTable(true)

	r.GET("/report", a.NewReport)
	r.Run(":9999")
}

//����
func (this *FirmwareService) stringToKuoHao(i []string) string {
	s := "("
	for k1, v1 := range i {
		s += "'" + v1 + "'"
		if k1 < len(i)-1 {
			s += ","
		}
	}
	s += ")"
	return s
}

//new 报告
func (this *FirmwareService) NewReport(c *gin.Context) {
	reportType := "ANHEN"
	RPL := reportWords(reportType)
	id := c.Query("id")
	var f model.Firmware22
	err := this.db.Where("id = ?", id).First(&f).Error
	if err != nil {
		fmt.Println(err)
		return
	}
	var fTaks model.FirmwareTask22
	err = this.db.Where("id = ?", f.FirmwareTask22Id).First(&fTaks).Error
	if err != nil {
		fmt.Println(err)
		return
	}
	var fs model.Software22
	err = this.db.Where("id = ?", f.Software22Id).First(&fs).Error
	if err != nil {
		fmt.Println(err)
	}
	var ss []model.Software22
	err = this.db.Where("firmware22_id = ?", id).Order("software_version asc").Find(&ss).Error
	if err != nil {
		fmt.Println(err)
		return
	}
	//查询报告定制
	var fgcs []model.FirmwareGeneralConfiguration
	this.db.Where("p_key = 'report_customization'").Find(&fgcs)
	fgcsmap := map[string]interface{}{}
	for _, v := range fgcs {
		fgcsmap[v.Key] = v
	}

	//厂商名称
	var TradeNames string
	//报告Logo图片
	var logoImage string
	//报告标题
	var reportTitle string
	//报告页眉
	var reportHeader string
	//报告页脚
	var reportFooter string
	//报告水印
	var reportWatermark string
	reportWatermark = "resources/null.png"
	//报告封面背景
	var reportCoverBackground string
	// 是否存在公司简介
	var report_company_profile bool
	var introductionTitle, introduction string
	report_company_profile = fgcsmap["report_company_profile"].(model.FirmwareGeneralConfiguration).IsBool
	{
		if fgcsmap["report_company_profile"] != nil && fgcsmap["report_company_profile"].(model.FirmwareGeneralConfiguration).IsBool {
			if fgcsmap["report_introduction_title"] != nil && fgcsmap["report_introduction_title"].(model.FirmwareGeneralConfiguration).Data != "" {
				introductionTitle = fmt.Sprintf("%s", fgcsmap["report_introduction_title"].(model.FirmwareGeneralConfiguration).Data)
			}
			if fgcsmap["manufacturer_introduction"] != nil && fgcsmap["manufacturer_introduction"].(model.FirmwareGeneralConfiguration).Data != "" {
				introduction = fmt.Sprintf("%s", fgcsmap["manufacturer_introduction"].(model.FirmwareGeneralConfiguration).Data)
			}
		}
		if fgcsmap["report_cover_background"] != nil && fgcsmap["report_cover_background"].(model.FirmwareGeneralConfiguration).Data != "" {
			reportCoverBackground = fmt.Sprintf("%s", fgcsmap["report_cover_background"].(model.FirmwareGeneralConfiguration).Data)
		}
		if fgcsmap["report_watermark"] != nil && fgcsmap["report_watermark"].(model.FirmwareGeneralConfiguration).Data != "" && fgcsmap["report_watermark"].(model.FirmwareGeneralConfiguration).IsBool {
			reportWatermark = fmt.Sprintf("%s", fgcsmap["report_watermark"].(model.FirmwareGeneralConfiguration).Data)
		}
		if fgcsmap["report_footer"] != nil && fgcsmap["report_footer"].(model.FirmwareGeneralConfiguration).Data != "" {
			reportFooter = fmt.Sprintf("%s", fgcsmap["report_footer"].(model.FirmwareGeneralConfiguration).Data)
		}
		if fgcsmap["report_header"] != nil && fgcsmap["report_header"].(model.FirmwareGeneralConfiguration).Data != "" {
			reportHeader = fmt.Sprintf("%s %s", f.Name, fgcsmap["report_header"].(model.FirmwareGeneralConfiguration).Data)
		}
		if fgcsmap["report_title"] != nil && fgcsmap["report_title"].(model.FirmwareGeneralConfiguration).Data != "" {
			reportTitle = fgcsmap["report_title"].(model.FirmwareGeneralConfiguration).Data
		}
		if fgcsmap["logo"] != nil && fgcsmap["logo"].(model.FirmwareGeneralConfiguration).Data != "" {
			logoImage = fgcsmap["logo"].(model.FirmwareGeneralConfiguration).Data
		}
		if fgcsmap["trade_names"] != nil && fgcsmap["trade_names"].(model.FirmwareGeneralConfiguration).Data != "" {
			TradeNames = fgcsmap["trade_names"].(model.FirmwareGeneralConfiguration).Data
		}
	}
	{
		fmt.Println(fmt.Sprintf("厂商名称：%s", TradeNames))
		if reportType == "ANHEN" {
			logoImage = "./pic/logo2.png"
		}else {
			logoImage = "./pic/logo.png"
		}
		fmt.Println(fmt.Sprintf("Logo图片：%s", utils.AbsPath(logoImage)))
		//_, err = os.Open(utils.AbsPath(logoImage))
		//if err != nil {
		//	err := this.client.FGetObject(context.Background(), "engine-input", filepath.Base(logoImage), utils.AbsPath(logoImage), minio.GetObjectOptions{})
		//	if err != nil {
		//		fmt.Println(err)
		//	}
		//}

		fmt.Println(fmt.Sprintf("报告标题：%s", reportTitle))
		fmt.Println(fmt.Sprintf("报告页眉：%s", reportHeader))
		fmt.Println(fmt.Sprintf("报告页脚：%s", reportFooter))
		fmt.Println(fmt.Sprintf("报告水印：%s", utils.AbsPath(reportWatermark)))
		//_, err = os.Open(utils.AbsPath(reportWatermark))
		//if err != nil {
		//	err := this.client.FGetObject(context.Background(), "engine-input", filepath.Base(reportWatermark), utils.AbsPath(reportWatermark), minio.GetObjectOptions{})
		//	if err != nil {
		//		fmt.Println(err)
		//	}
		//}
		reportCoverBackground = "./pic/backendground.png"
		fmt.Println(fmt.Sprintf("报告封面背景：%s", utils.AbsPath(reportCoverBackground)))
		//_, err = os.Open(utils.AbsPath(reportCoverBackground))
		//if err != nil {
		//	err := this.client.FGetObject(context.Background(), "engine-input", filepath.Base(reportCoverBackground), utils.AbsPath(reportCoverBackground), minio.GetObjectOptions{})
		//	if err != nil {
		//		fmt.Println(err)
		//	}
		//}
		fmt.Println(fmt.Sprintf("关于：%s %s", introductionTitle, introduction))
	}

	var vs model.VulnerabilityStatistics
	ssMap := map[uint]model.Software22{}
	sSub := this.db.Model(&model.Software22{}).Select("id").Where("firmware22_id = ?", id).Order("software_version asc").SubQuery()
	softwareAndCveLookup := map[uint]SoftwareAndCveLookup{}
	cweMap := map[string]CweCheckerAndInfo{}
	cveIds := map[string]utils.Cve{}
	//cweAndCweChecker := map[string][]model.CweChecker22{}
	var (
		usersPasswords      []model.UsersPasswords22
		cryptoMaterials     []model.CryptoMaterial22
		wpahardcodeCheckers []model.WpahardcodeChecker22
		svnCheckers         []model.SvnChecker22
		vimCheckers         []model.Software22
		tmpfileChecker      []model.Software22
		bakfileChecker      []model.Software22
		sourcecodeLeakage   []model.Software22
		gitCheckers         []model.GitChecker22
		softwares           []model.Software22
		malwareScanners     []model.MalwareScanner22
		cveLookups          []model.CveLookup22
		cweCheckers         []model.CweChecker22
		//cweInfos            []model.FirmwareDetectionVulnerability
		exploitMitigations  []model.ExploitMitigations22
		configsets          []model.Configset22
		sshs                []LsConfigset
		ftps                []LsConfigset
		hardcodes           []LsConfigset
		selfstartingChecker []model.SelfstartingChecker22
		registryCheckers    []model.RegistryChecker22
		riskTotal, riskHigh, riskIn, riskLow,
		cveTotal, cveHigh, cveIn, cveLow int
		mitigationMechanismVulnerabilities int
	)
	{
		//users_passwords 明⽂⽤户密码检测
		{
			err := this.db.Where("software22_id in ? and cracked = 'true'", sSub).Find(&usersPasswords).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				fmt.Println(err)
			}
		}
		//crypto_material 证书、密钥⽂件泄露检测
		{
			err := this.db.Where("software22_id in ?", sSub).Find(&cryptoMaterials).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				fmt.Println(err)
			}
		}
		//wpahardcode_checker WPA密码硬编码检测
		{
			err := this.db.Where("software22_id in ?", sSub).Find(&wpahardcodeCheckers).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				fmt.Println(err)
			}
		}
		//svn_checker SVN信息泄露检测
		{
			err := this.db.Where("software22_id in ?", sSub).Find(&svnCheckers).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				fmt.Println(err)
			}
		}
		for _, v := range ss {
			ssMap[v.ID] = v
			//vim_checker vi/vim信息泄露检测
			if v.VimChecker != "" {
				vimCheckers = append(vimCheckers, v)
			}
			//tmpfile_checker 临时⽂件泄露检测
			if v.TmpfileChecker != "" {
				tmpfileChecker = append(tmpfileChecker, v)
			}
			//bakfile_checker 备份⽂件泄露检测
			if v.BakfileChecker != "" {
				bakfileChecker = append(bakfileChecker, v)
			}
			//sourcecode_leakage 源代码⽂件泄露检测
			if v.SourcecodeLeakage != "" {
				sourcecodeLeakage = append(sourcecodeLeakage, v)
			}
		}
		//git_checker Git信息泄露检测
		{
			err := this.db.Where("software22_id in ?", sSub).Find(&gitCheckers).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				fmt.Println(err)
			}
		}
		//软件成分
		{
			err := this.db.Where("firmware22_id = ? and software_name != '' and software_version != '--' and software_version != ''", id).Order("software_version asc").Find(&softwares).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				fmt.Println(err)
			}
		}
		//malware_scanner 恶意软件检测
		{
			err := this.db.Where("software22_id in ?", sSub).Find(&malwareScanners).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				fmt.Println(err)
			}
		}
		//软件组件漏洞检测
		{
			err := this.db.Where("software22_id in ?", sSub).Find(&cveLookups).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				fmt.Println(err)
			} else {
				var cveIdStrs []string
				for _, v := range cveLookups {
					_, ok := cveIds[v.CveId]
					if !ok {
						cveIdStrs = append(cveIdStrs, v.CveId)
						cveIds[v.CveId] = utils.Cve{}
					}
				}
				var cnnvds []utils.Cve
				if len(cveIdStrs) > 0 {
					err := this.edb.Table("cnnvd").Where("cve_id in " + this.stringToKuoHao(cveIdStrs)).Find(&cnnvds).Error
					if err != nil {
						fmt.Println(err)
					}
					for _, v := range cnnvds {
						cveIds[v.CveId] = v
					}
				}
				for _, v := range cveLookups {
					s, ok := ssMap[v.Software22Id]
					if ok {
						sc, ok := softwareAndCveLookup[s.ID]
						var sAndC SoftwareAndCveLookup
						if ok {
							sAndC = sc
							cves := sAndC.CveLookups
							cves = append(cves, v)
							sAndC.CveLookups = cves
						} else {
							sAndC = SoftwareAndCveLookup{
								CveLookups: []model.CveLookup22{v},
							}
						}
						//softwareAndCveLookup[s.ID].
						var ls []model.CveLookup22
						switch v.HazardGrade {
						case "超危", "高危":
							if sAndC.CveHighLookups != nil && len(sAndC.CveHighLookups) != 0 {
								ls = sAndC.CveHighLookups
							}
							ls = append(ls, v)
							sAndC.CveHighLookups = ls
						case "中危":
							if sAndC.CveInLookups != nil && len(sAndC.CveInLookups) != 0 {
								ls = sAndC.CveInLookups
							}
							ls = append(ls, v)
							sAndC.CveInLookups = ls
						case "低危":
							if sAndC.CveLowLookups != nil && len(sAndC.CveLowLookups) != 0 {
								ls = sAndC.CveLowLookups
							}
							ls = append(ls, v)
							sAndC.CveLowLookups = ls
						}
						softwareAndCveLookup[s.ID] = sAndC
					}
				}
			}
		}
		//cwe_checker 代码安全漏洞检测
		//{
		//	err := this.db.Where("software22_id in ? and cwe_id in ('CWE190','CWE215','CWE243','CWE248','CWE332','CWE367','CWE426','CWE457','CWE467','CWE476','CWE676','CWE782','CWE560')", sSub).Find(&cweCheckers).Error
		//	if err != nil && err != gorm.ErrRecordNotFound {
		//				fmt.Println(err)
		//	}
		//
		//	for _, v := range []string{"CWE190", "CWE215", "CWE243", "CWE248", "CWE332", "CWE367", "CWE426", "CWE457", "CWE467", "CWE476", "CWE676", "CWE782", "CWE560"} {
		//		cweAndCweChecker[v] = []model.CweChecker22{}
		//	}
		//	for _, v := range cweCheckers {
		//		cAndC, ok := cweAndCweChecker[v.CweId]
		//		if ok {
		//			cAndC = append(cAndC, v)
		//			cweAndCweChecker[v.CweId] = cAndC
		//		} else {
		//			var cAndC []model.CweChecker22
		//			cAndC = append(cAndC, v)
		//			cweAndCweChecker[v.CweId] = cAndC
		//		}
		//	}
		//	err = this.db.Where("l_id in ('CWE190','CWE215','CWE243','CWE248','CWE332','CWE367','CWE426','CWE457','CWE467','CWE476','CWE676','CWE782','CWE560')").Find(&cweInfos).Error
		//	if err != nil && err != gorm.ErrRecordNotFound {
		//				fmt.Println(err)
		//	}
		//	for _, v := range cweInfos {
		//		ls, ok := cweAndCweChecker[v.LId]
		//		if ok {
		//			cweMap[v.LId] = CweCheckerAndInfo{
		//				Cwe:  v,
		//				Cwes: ls,
		//			}
		//		}
		//	}
		//}
		//exploit_mitigations ⼆进制缓解措施检测
		{
			err := this.db.Where("software22_id in ?", sSub).Find(&exploitMitigations).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				fmt.Println(err)
			}
		}
		//配置安全检测
		{
			err := this.db.Where("software22_id in ?", sSub).Find(&configsets).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				fmt.Println(err)
			}
			for _, v := range configsets {
				if v.Ssh != "" && v.Ssh != "null" {
					var sshStrs []string
					err := jsoniter.Unmarshal([]byte(v.Ssh), &sshStrs)
					if err != nil {
						fmt.Println(err)
					} else {
						var ls LsConfigset
						ls.Data = sshStrs
						ls.Software = ssMap[v.Software22Id]
						sshs = append(sshs, ls)
					}
				}
				if v.Ftp != "" && v.Ftp != "null" {
					var ftpStrs []string
					err := jsoniter.Unmarshal([]byte(v.Ftp), &ftpStrs)
					if err != nil {
						fmt.Println(err)
					} else {
						var ls LsConfigset
						ls.Data = ftpStrs
						ls.Software = ssMap[v.Software22Id]
						ftps = append(ftps, ls)
					}
				}
				if v.Hardcode != "" && v.Hardcode != "null" {
					var hardcodeStrs []string
					err := jsoniter.Unmarshal([]byte(v.Hardcode), &hardcodeStrs)
					if err != nil {
						fmt.Println(err)
					} else {
						var ls LsConfigset
						ls.Data = hardcodeStrs
						ls.Software = ssMap[v.Software22Id]
						hardcodes = append(hardcodes, ls)
					}
				}
			}
			err = this.db.Where("software22_id in ?", sSub).Find(&selfstartingChecker).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				fmt.Println(err)
			}
		}
		//registry_checker 注册表⽂件检测
		{
			err := this.db.Where("software22_id in ?", sSub).Find(&registryCheckers).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				fmt.Println(err)
			}
		}
		//统计数量
		{
			for _, _ = range exploitMitigations {
				// 去除逻辑判断
				//if strings.Contains(v.NXStatus, "fully") &&
				//	strings.Contains(v.RelroStatus, "fully") &&
				//	strings.Contains(v.CanaryStatus, "fully") &&
				//	strings.Contains(v.FortifySourceStatus, "fully") &&
				//	strings.Contains(v.PIEStatus, "fully") {
				mitigationMechanismVulnerabilities++
				//}
			}
			riskHigh += len(cryptoMaterials) + len(malwareScanners)
			riskIn += len(usersPasswords) + len(wpahardcodeCheckers) + len(registryCheckers) +
				len(svnCheckers) + len(sourcecodeLeakage) + len(gitCheckers) + len(sshs) + len(ftps) + len(hardcodes)
			riskLow += len(vimCheckers) + len(bakfileChecker) + len(tmpfileChecker) + mitigationMechanismVulnerabilities
			riskTotal += len(usersPasswords) + len(cryptoMaterials) + len(gitCheckers) +
				len(wpahardcodeCheckers) + len(svnCheckers) + len(sourcecodeLeakage) +
				len(gitCheckers) + len(vimCheckers) + len(bakfileChecker) + len(tmpfileChecker) + len(registryCheckers) +
				len(malwareScanners) + len(sshs) + len(ftps) + len(hardcodes) + mitigationMechanismVulnerabilities

			for _, v := range softwareAndCveLookup {
				riskTotal += len(v.CveLookups)
				cveTotal += len(v.CveLookups)
				if len(v.CveHighLookups) > 0 {
					cveHigh += len(v.CveHighLookups)
					riskHigh += len(v.CveHighLookups)
				}
				if len(v.CveInLookups) > 0 {
					cveIn += len(v.CveInLookups)
					riskIn += len(v.CveInLookups)
				}
				if len(v.CveLowLookups) > 0 {
					cveLow += len(v.CveLowLookups)
					riskLow += len(v.CveLowLookups)
				}
			}
			for _, v := range cweMap {
				switch v.Cwe.RiskLevel {
				case "高":
					riskHigh += len(v.Cwes)
				case "中":
					riskIn += len(v.Cwes)
				case "低":
					riskLow += len(v.Cwes)
				}
				riskTotal += len(v.Cwes)
			}
			{
				err = this.db.Where("firmware22_id = ?", f.ID).First(&vs).Error
				if err != nil && err != gorm.ErrRecordNotFound {
					fmt.Println(err)
					return
				}
				if vs.ID == 0 {
					vs.Firmware22Id = f.ID
					if len(registryCheckers) > 0 {
						vs.RegistryChecker = true
					}
					if len(wpahardcodeCheckers) > 0 {
						vs.WpahardcodeChecker = true
					}
					if len(vimCheckers) > 0 {
						vs.VimChecker = true
					}
					if len(tmpfileChecker) > 0 {
						vs.TmpfileChecker = true
					}
					if len(svnCheckers) > 0 {
						vs.SvnChecker = true
					}
					if len(sourcecodeLeakage) > 0 {
						vs.SourcecodeLeakage = true
					}
					if len(registryCheckers) > 0 {
						vs.SelfstartingChecker = true
					}
					if len(gitCheckers) > 0 {
						vs.GitChecker = true
					}
					if len(cweCheckers) > 0 {
						vs.CweChecker = true
					}
					if len(sshs) > 0 {
						vs.SSH = true
					}
					if len(ftps) > 0 {
						vs.FTP = true
					}
					if len(cryptoMaterials) > 0 {
						vs.CryptoMaterial = true
					}
					if len(cveLookups) > 0 {
						vs.CveLookup = true
					}
					if mitigationMechanismVulnerabilities > 0 {
						vs.ExploitMitigations = true
					}
					if len(hardcodes) > 0 {
						vs.Hardcode = true
					}
					if len(usersPasswords) > 0 {
						vs.UsersPasswords = true
					}
					if len(malwareScanners) > 0 {
						vs.MalwareScanner = true
					}
					if len(bakfileChecker) > 0 {
						vs.BakfileChecker = true
					}
					err = this.db.Create(&vs).Error
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}
	var s utils.Report
	//s.InitTemplate = "初始化报告文件路径"
	s.InitTemplate = utils.AbsPath("resources/gujianbaogao/fengxi/init.docx")
	//页眉 + 水印
	//{
	//	header := &utils.Header{}
	//	header.DefaultTemplate = "header_default.ftl"
	//	header.Data.HeaderText = reportHeader
	//	var references []*utils.References
	//	references = append(references, &utils.References{
	//		Type: "watermark",
	//		Name: "watermarkRelationshipId",
	//		Property: map[string]interface{}{
	//			"imagePath": utils.AbsPath(reportWatermark),
	//			//"imagePath": "E:\\ky\\ky\\ssp\\backend\\reportj\\data\\watermark.png",
	//		},
	//	})
	//	header.Data.References = references
	//	s.Header = header
	//}
	//页脚
	//{
	//	s.Footer.DefaultTemplate = "footer_default.ftl"
	//	s.Footer.Data.FooterText = reportFooter
	//}
	//标题样式
	{
		titleStyle := &utils.TitleStyle{}
		titleStyle.LevelStyle = map[string]interface{}{
			"1": "title_level1.ftl",
			"2": "title_level2.ftl",
			"3": "title_level3.ftl",
			"4": "title_level4.ftl",
		}
		titleStyle.Data = map[string]interface{}{}
		s.TitleStyle = titleStyle
	}
	//目录
	{
		toc := &utils.Toc{}
		toc.Template = "toc.ftl"
		toc.Data = map[string]interface{}{
			"headingText": "目录",
		}

		s.Toc = toc
	}
	//封面
	{
		cover := &utils.Cover{}
		cover.Template = "cover.ftl"
		var references []utils.References
		references = append(references, utils.References{
			Type: "image",
			Name: "bg_img",
			Property: map[string]interface{}{
				"imagePath": utils.AbsPath(reportCoverBackground),
			},
		}, utils.References{
			Type: "image",
			Name: "logo_img",
			Property: map[string]interface{}{
				"imagePath": utils.AbsPath(logoImage),
			},
		})
		cover.Data = map[string]interface{}{
			//"app_name":   fmt.Sprintf("%s-%s %s", TradeNames, fTaks.DeviceName, reportTitle),
			"app_name": fmt.Sprintf("%s", fTaks.DeviceName),
			//"org_name":   "中国电子技术标准化研究",
			"org_name":   TradeNames,
			"title":      "固件安全检测报告",
			"references": references,
			"check_time": time.Now().Format("2006-01-02 15:04:05"),
		}
		s.Cover = cover
	}
	//关于
	{
		aboutImage := utils.AbsPath("resources/about.png")
		about := &utils.About{}
		about.Template = "about.ftl"
		//about.Data.Title = "中国电子技术标准化研究院固件安全分析平台"
		about.Data.Title = TradeNames + "固件安全分析报告"
		about.Data.Introduction = introduction
		about.Data.Website = introductionTitle
		about.Data.IsSHow = report_company_profile
		var references []*utils.References
		references = append(references, &utils.References{
			Type: "image",
			Name: "relationshipId",
			Property: map[string]interface{}{
				"imagePath": aboutImage,
			},
		})
		about.Data.References = references
		s.About = about
	}
	//节点
	{
		document := &utils.Document{}
		document.Level = &utils.Level{}
		document.Title = "root"
		var all []*utils.Document
		//固件安全检测总览
		var cpuStrs []string
		if fs.CpuArchitecture != "" {
			err := jsoniter.Unmarshal([]byte(fs.CpuArchitecture), &cpuStrs)
			if err != nil {
				fmt.Println(err)
			}
		}
		var fileSystemStr []string
		if fs.FileSystem != "" {
			err = jsoniter.Unmarshal([]byte(fs.FileSystem), &fileSystemStr)
			if err != nil {
				fmt.Println(err)
			}
		}
		var firmwareTypeStr []string
		if fs.FileType != "" {
			err = jsoniter.Unmarshal([]byte(fs.FileType), &firmwareTypeStr)
			fmt.Println(fmt.Sprintf(fs.FileType))
			if err != nil {
				fmt.Println(err)
			}
		}
		{
			ls := &utils.Document{}
			ls.Title = "固件安全检测总览"
			ls.Level = &utils.Level{
				ParentLevel: "",
				Format:      "1.",
				Order:       1,
			}
			var overviewFirmwareSecurityInspection []*utils.Document
			//固件检测整体情况
			{
				theOverallSituationFirmwareDetection := &utils.Document{}
				theOverallSituationFirmwareDetection.Title = "固件检测整体情况"
				theOverallSituationFirmwareDetection.Template = "section11.ftl"
				var tongJi [][]string
				tongJi = append(tongJi, []string{"发现问题总数", fmt.Sprintf("%d", riskTotal), fmt.Sprintf("%d", riskHigh),
					fmt.Sprintf("%d", riskIn), fmt.Sprintf("%d", riskLow)})
				tongJi = append(tongJi, []string{"CVE漏洞总数", fmt.Sprintf("%d", cveTotal), fmt.Sprintf("%d", cveHigh),
					fmt.Sprintf("%d", cveIn), fmt.Sprintf("%d", cveLow)})
				theOverallSituationFirmwareDetection.Data = map[string]interface{}{
					"table": tongJi,
				}
				theOverallSituationFirmwareDetection.Level = &utils.Level{
					ParentLevel: "1",
					Format:      "1.1.",
					Order:       1,
				}
				theOverallSituationFirmwareDetection.Children = []*utils.Document{}
				overviewFirmwareSecurityInspection = append(overviewFirmwareSecurityInspection, theOverallSituationFirmwareDetection)
			}
			//固件基本信息
			{
				theOverallSituationFirmwareDetection := &utils.Document{}
				theOverallSituationFirmwareDetection.Title = "固件基本信息"
				theOverallSituationFirmwareDetection.Template = "section12.ftl"
				var jiBen [][]string
				jiBen = append(jiBen, [][]string{
					{"文件名称", f.Name},
					{"设备名称", fTaks.DeviceName},
					{"厂商", fTaks.TradeNames},
					{"设备类型", fTaks.EquipmentType},
					{"文件大小", f.Size},
					{"分析时间", f.CreatedAt.Format("2006-01-02 15:04:05")},
				}...)
				theOverallSituationFirmwareDetection.Data = map[string]interface{}{
					"table": jiBen,
				}
				theOverallSituationFirmwareDetection.Level = &utils.Level{
					ParentLevel: "1",
					Format:      "1.2.",
					Order:       2,
				}
				theOverallSituationFirmwareDetection.Children = []*utils.Document{}
				overviewFirmwareSecurityInspection = append(overviewFirmwareSecurityInspection, theOverallSituationFirmwareDetection)
			}
			//固件安全检测结果总结
			{
				theOverallSituationFirmwareDetection := &utils.Document{}
				theOverallSituationFirmwareDetection.Title = "固件安全检测结果总结"
				theOverallSituationFirmwareDetection.Template = "section13.ftl"
				//"自身安全（5项）",
				theOverallSituationFirmwareDetection.Data = map[string]interface{}{
					"tblheaders": []string{RPL.SIDTitle, RPL.SoftCI+"（1项）", RPL.SoftCV, RPL.BadWare,
						RPL.ConfigSecurity},
					//"代码安全漏洞检测（1项）",
					"tbls": [][][]string{
						//{
						//	{"1", "文件hashes", "/", utils.GetStr3(fs.MD5, fs.SHA512, fs.SHA256)},
						//	{"2", "CPU架构", "/", utils.GetStr2(cpuStrs)},
						//	{"3", "操作系统", "/", utils.GetStr1(f.OperatingSystem)},
						//	{"4", "文件系统", "/", utils.GetStr2(fileSystemStr)},
						//	{"5", "固件类型", "/", utils.GetStr2(firmwareTypeStr)},
						//},
						{
							{"1", RPL.WeakPWD, "中", utils.CunZaiFengXian(len(usersPasswords))},
							{"2", RPL.KeyStore, "高", utils.CunZaiFengXian(len(cryptoMaterials))},
							{"3", RPL.WpaPWD, "中", utils.CunZaiFengXian(len(wpahardcodeCheckers))},
							{"4", RPL.SVNMessage, "中", utils.CunZaiFengXian(len(svnCheckers))},
							{"5", RPL.SourceCode, "中", utils.CunZaiFengXian(len(sourcecodeLeakage))},
							{"6", RPL.GitMessage, "中", utils.CunZaiFengXian(len(gitCheckers))},
							{"7", RPL.ViMMessage, "低", utils.CunZaiFengXian(len(vimCheckers))},
							{"8", RPL.CopyFile, "低", utils.CunZaiFengXian(len(bakfileChecker))},
							{"9", RPL.TemporaryFile, "低", utils.CunZaiFengXian(len(tmpfileChecker))},
							{"10", RPL.ResignTable, "中", utils.CunZaiFengXian(len(registryCheckers))},
						},
						{
							{"11", RPL.SoftCITwo, "提示信息", utils.CunZaiFengXian2(len(softwares))},
						},
						{
							{"12", RPL.SoftCVTwo, "/", utils.CunZaiFengXian2(cveTotal)},
						},
						{
							{"13", RPL.BadWareTwo, "高", utils.CunZaiFengXian(len(malwareScanners))},
						},
						//{
						//{"14", "整数溢出缺陷检测", "低", utils.CunZaiFengXian(len(cweMap["CWE190"].Cwes))},
						//{"15", "调试信息泄露风险", "中", utils.CunZaiFengXian(len(cweMap["CWE215"].Cwes))},
						//{"16", "chroot不安全使用漏洞", "中", utils.CunZaiFengXian(len(cweMap["CWE243"].Cwes))},
						//{"17", "未捕获异常风险", "中", utils.CunZaiFengXian(len(cweMap["CWE248"].Cwes))},
						//{"18", "伪随机数使用漏洞", "中", utils.CunZaiFengXian(len(cweMap["CWE332"].Cwes))},
						//{"19", "条件竞争漏洞", "中", utils.CunZaiFengXian(len(cweMap["CWE367"].Cwes))},
						//{"20", "不受信任的搜索路径检测", "高", utils.CunZaiFengXian(len(cweMap["CWE426"].Cwes))},
						//{"21", "使用未初始化变量风险", "中", utils.CunZaiFengXian(len(cweMap["CWE457"].Cwes))},
						//{"22", "sizeof()使用不当漏洞", "低", utils.CunZaiFengXian(len(cweMap["CWE467"].Cwes))},
						//{"23", "空指针错误引用漏洞", "中", utils.CunZaiFengXian(len(cweMap["CWE476"].Cwes))},
						//{"24", "不安全函数调用风险", "中", utils.CunZaiFengXian(len(cweMap["CWE676"].Cwes))},
						//{"25", "IOCTL调用未控制权限风险", "低", utils.CunZaiFengXian(len(cweMap["CWE782"].Cwes))},
						//{"26", "umask()参数不正确使用风险", "低", utils.CunZaiFengXian(len(cweMap["CWE560"].Cwes))},
						//{"14", "安全缓解机制检测", "低", utils.CunZaiFengXian(mitigationMechanismVulnerabilities)},
						//},
						{
							{"15", RPL.HardCode, "中", utils.CunZaiFengXian(len(hardcodes))},
							{"16", RPL.SSHSecurity, "中", utils.CunZaiFengXian(len(sshs))},
							{"17", RPL.FTPSecurity, "中", utils.CunZaiFengXian(len(ftps))},
							{"18", RPL.SelfStartSR, "提示", utils.CunZaiFengXian(len(selfstartingChecker))},
							{"19", RPL.SafetyMM, "低", utils.CunZaiFengXian(mitigationMechanismVulnerabilities)},
						},
					},
				}
				theOverallSituationFirmwareDetection.Level = &utils.Level{
					ParentLevel: "1",
					Format:      "1.3.",
					Order:       3,
				}
				theOverallSituationFirmwareDetection.Children = []*utils.Document{}
				overviewFirmwareSecurityInspection = append(overviewFirmwareSecurityInspection, theOverallSituationFirmwareDetection)
			}
			ls.Data = map[string]interface{}{}
			ls.Children = overviewFirmwareSecurityInspection
			all = append(all, ls)
		}
		//固件安全检测结果详情
		{
			ls := &utils.Document{}
			ls.Title = "固件安全检测结果详情"
			ls.Level = &utils.Level{
				ParentLevel: "",
				Format:      "2.",
				Order:       2,
			}
			ls.Children = []*utils.Document{}
			ls.Data = map[string]interface{}{}
			var detailsFirmwareSecurityTestResults []*utils.Document
			//自身安全
			{
				oneOwnSafety := &utils.Document{}
				oneOwnSafety.Title = ".固件信息"
				oneOwnSafety.Level = &utils.Level{
					ParentLevel: "2",
					Format:      "2.1",
					Order:       1,
				}
				var ls []*utils.Document
				//文件hash
				{
					hash := &utils.Document{}
					hash.Title = ".文件Hash"
					hash.Template = "section211.ftl"
					hash.Level = &utils.Level{
						ParentLevel: "2.1",
						Format:      "2.1.1",
						Order:       1,
					}
					hash.Data = map[string]interface{}{
						"table": [][]string{
							{"md5", fs.MD5},
							{"sha256", fs.SHA256},
							{"sha512", fs.SHA512},
						},
					}
					hash.Children = []*utils.Document{}
					ls = append(ls, hash)
				}
				//CPU架构
				{
					cpu := &utils.Document{}
					cpu.Title = ".可适配CPU架构"
					cpu.Template = "gjxx.ftl"
					cpu.Level = &utils.Level{
						ParentLevel: "2.1",
						Format:      "2.1.2",
						Order:       2,
					}
					var strss [][]string
					for k, v := range cpuStrs {
						strss = append(strss, []string{fmt.Sprintf("%d", k+1), v})
					}
					if len(strss) == 0 {
						strss = append(strss, []string{"1", "此固件未检测出CPU架构"})
					}
					cpu.Data = map[string]interface{}{
						"table": cpuStrs,
						"str":   "可适配CPU架构",
						"strss": strss,
					}
					cpu.Children = []*utils.Document{}
					ls = append(ls, cpu)
				}
				//操作系统
				{
					operatingSystem := &utils.Document{}
					operatingSystem.Title = ".操作系统"
					operatingSystem.Template = "gjxx.ftl"
					operatingSystem.Level = &utils.Level{
						ParentLevel: "2.1",
						Format:      "2.1.3",
						Order:       3,
					}
					var strss [][]string
					if f.OperatingSystem == "" {
						strss = append(strss, []string{"1", "此固件未检测出操作系统"})
					} else {
						strss = append(strss, []string{"1", f.OperatingSystem})
					}
					operatingSystem.Data = map[string]interface{}{
						"table": []string{f.OperatingSystem},
						"str":   "操作系统名称",
						"strss": strss,
					}
					operatingSystem.Children = []*utils.Document{}
					ls = append(ls, operatingSystem)
				}
				//文件系统
				{
					fileSystem := &utils.Document{}
					fileSystem.Title = ".文件系统"
					fileSystem.Template = "gjxx.ftl"
					fileSystem.Level = &utils.Level{
						ParentLevel: "2.1",
						Format:      "2.1.4",
						Order:       4,
					}
					fileSystem.Children = []*utils.Document{}
					var strss [][]string
					for k, v := range fileSystemStr {
						var finalStr string
						// 去重
						{
							itemList := strings.Split(v, ",")
							var strList []string
						C:
							for _, item := range itemList {
								for _, s2 := range strList {
									if s2 == item {
										continue C
									}
								}
								strList = append(strList, item)
							}
							for idx, i := range strList {
								if idx != 0 {
									finalStr += "," + i
								} else {
									finalStr += i
								}
							}
						}
						strss = append(strss, []string{fmt.Sprintf("%d", k+1), finalStr})
					}
					if len(fileSystemStr) == 0 {
						strss = append(strss, []string{"1", "此固件未检出文件系统"})
					}
					fileSystem.Data = map[string]interface{}{
						"table": fileSystemStr,
						"str":   "文件系统名称",
						"strss": strss,
					}
					ls = append(ls, fileSystem)
				}
				//固件类型
				{
					firmwareType := &utils.Document{}
					firmwareType.Title = ".固件类型"
					firmwareType.Template = "gjxx.ftl"
					firmwareType.Level = &utils.Level{
						ParentLevel: "2.1",
						Format:      "2.1.5",
						Order:       5,
					}
					firmwareType.Children = []*utils.Document{}
					var strss [][]string
					for k, v := range firmwareTypeStr {
						strss = append(strss, []string{fmt.Sprintf("%d", k+1), v})
					}
					if len(firmwareTypeStr) == 0 {
						strss = append(strss, []string{"1", "此固件未检出固件类型"})
					}
					firmwareType.Data = map[string]interface{}{
						"table": firmwareTypeStr,
						"str":   "固件类型",
						"strss": strss,
					}
					ls = append(ls, firmwareType)
				}
				oneOwnSafety.Data = map[string]interface{}{}
				oneOwnSafety.Children = ls
				detailsFirmwareSecurityTestResults = append(detailsFirmwareSecurityTestResults, oneOwnSafety)
			}
			//敏感信息检测
			{
				sensitiveInformationDetection := &utils.Document{}
				sensitiveInformationDetection.Title = "." + RPL.SIDTitleTwo
				sensitiveInformationDetection.Level = &utils.Level{
					ParentLevel: "2",
					Format:      "2.2",
					Order:       2,
				}
				var ls []*utils.Document
				//用户密码信息泄露风险
				{
					ls1 := &utils.Document{}
					ls1.Title = "." + RPL.WeakPWD
					ls1.Template = "section221.ftl"
					ls1.Level = &utils.Level{
						ParentLevel: "2.2",
						Format:      "2.2.1",
						Order:       1,
					}
					str := "经检测以下文件中存在相应风险，详细信息如下。\n"
					for kk, v := range usersPasswords {
						_, ok := ssMap[v.Software22Id]
						if !ok {
							continue
						}
						wjm := ssMap[v.Software22Id].FileName
						if wjm == "" {
							continue
						}
						//str += fmt.Sprintf("%d.文件名：%s\n", kk+1, wjm)
						str += fmt.Sprintf("%d.文件路径：%s\n", kk+1, ssMap[v.Software22Id].FilePath)
						if utils.IsStr(v.Name) {
							str += fmt.Sprintf("   泄露用户名：%s\n", v.Name)
						}
						if utils.IsStr(v.Password) {
							str += fmt.Sprintf("   密码HASH：%s （%s）\n", v.Password, utils.PoJi(fmt.Sprintf("%v", v.Cracked)))
						}
						if utils.IsStr(v.Password) {
							str += fmt.Sprintf("   密码：%s\n", v.Password)
						}
						//if utils.IsStr(v.Cracked) {
						//	str += fmt.Sprintf("   cracked：%s\n", )
						//}
					}
					if str == "" || str == "经检测以下文件中存在相应风险，详细信息如下。\n" {
						str = "未检测到弱密码泄露风险"
					}
					var datas [][]string
					datas = append(datas, []string{
						"检测固件中是否存在易被破解的弱密码信息", "中", "使用弱密码是造成用户信息泄露和群体性的网络安全攻击的重要原因。此外，使用单一的密码认证已经被证实是不安全的做法。用户为了便于记忆，常常在密码中加入一些固定规律，这使得攻击者通过撞库破解的成功率大大提升。系统管理员账号如果使用弱密码，一旦被攻击，可能会导致整个系统内的数据库信息被窃取、业务系统瘫痪等安全问题，造成用户信息的泄露和经济损失。",
						utils.CunZaiFengXian1(len(usersPasswords)), str, utils.CunZaiFengXian4(len(usersPasswords), "禁止用户使用弱密码存储密码信息"),
					})
					ls1.Data = map[string]interface{}{
						"table": datas,
					}
					ls1.Children = []*utils.Document{}
					ls = append(ls, ls1)
				}
				//证书文件和密钥泄露风险
				{
					ls1 := &utils.Document{}
					ls1.Title = "." + RPL.KeyStore
					ls1.Template = "section221.ftl"
					ls1.Level = &utils.Level{
						ParentLevel: "2.2",
						Format:      "2.2.2",
						Order:       2,
					}
					str := "经检测以下文件中存在相应风险，详细信息如下。\n"
					for kk, v := range cryptoMaterials {
						_, ok := ssMap[v.Software22Id]
						if !ok {
							continue
						}
						wjm := ssMap[v.Software22Id].FileName
						if wjm == "" {
							continue
						}
						//str += fmt.Sprintf("%d.文件名：%s\n", kk+1, wjm)
						str += fmt.Sprintf("%d.文件路径：%s\n", kk+1, ssMap[v.Software22Id].FilePath)
						if v.Key != "" {
							str += fmt.Sprintf("   文件类型：%s\n", this.CryptoMaterialControl[v.Key])
						}
						//if utils.IsStr(v.Val) {
						//	var ls []string
						//	err := jsoniter.Unmarshal([]byte(v.Val), &ls)
						//	if err != nil {
						//				fmt.Println(err)
						//	}
						//	str += fmt.Sprintf("   详细信息\n")
						//	for _, v := range ls {
						//		fg := strings.Split(v, "\n")
						//		for _, v := range fg {
						//			v = strings.ReplaceAll(v, "：", "")
						//			v = strings.ReplaceAll(v, ":", "")
						//			str += fmt.Sprintf("      %s：\n", v)
						//		}
						//	}
						//}
					}
					if str == "" || str == "经检测以下文件中存在相应风险，详细信息如下。\n" {
						str = "未检测到" + RPL.KeyStore

					}
					var datas [][]string
					datas = append(datas, []string{
						"检测固件中是否存在证书文件和密钥泄露的风险",
						"高",
						"软件开发者如果在开源固件包里中泄露了公司内部使用的私有代码签名密钥或证书文件，恶意攻击者就可以利用其对恶意软件进行签名，从而仿冒正版的软件，欺骗普通用户进行安装。使用仿冒的软件，会导致用户利益受损，影响企业名誉。",
						utils.CunZaiFengXian1(len(cryptoMaterials)),
						str,
						utils.CunZaiFengXian4(len(cryptoMaterials), "厂商在固件发布版本之前，确保删除私钥文件，避免密钥信息泄露后被黑客非法利用。"),
					})
					ls1.Data = map[string]interface{}{
						"table": datas,
					}
					ls1.Children = []*utils.Document{}
					ls = append(ls, ls1)
				}
				//WPA密码硬编码检测
				{
					ls1 := &utils.Document{}
					ls1.Title = "." + RPL.WpaPWD
					ls1.Template = "section221.ftl"
					ls1.Level = &utils.Level{
						ParentLevel: "2.2",
						Format:      "2.2.3",
						Order:       3,
					}
					str := "经检测以下文件中存在相应风险，详细信息如下。\n"
					for kk, v := range wpahardcodeCheckers {
						_, ok := ssMap[v.Software22Id]
						if !ok {
							continue
						}
						wjm := ssMap[v.Software22Id].FileName
						if wjm == "" {
							continue
						}
						//str += fmt.Sprintf("%d.文件名：%s\n", kk+1, wjm)
						str += fmt.Sprintf("%d.文件路径：%s\n", kk+1, ssMap[v.Software22Id].FilePath)
						if v.Description != "" {
							str += fmt.Sprintf("   WPA hardcode：%s\n", v.Description)
						}
						if utils.IsStr(v.RelatedStrings) {
							var ls []string
							err := jsoniter.Unmarshal([]byte(v.RelatedStrings), &ls)
							if err != nil {
								fmt.Println(err)
							}
							for _, v := range ls {
								v = strings.ReplaceAll(v, "\\x", "")
								v = strings.ReplaceAll(v, "\"", "")
								v = strings.ReplaceAll(v, "`", "")
								if v[strings.LastIndex(v, "=")+1:] == "" {
									continue
								}
								str += fmt.Sprintf("      WPA硬编码结果：%s\n", v)
							}
						}
					}
					if str == "" || str == "经检测以下文件中存在相应风险，详细信息如下。\n" {
						str = "未检测到WPA密码是否硬编码"
					}
					var datas [][]string
					datas = append(datas, []string{
						"检测固件中的WPA密码是否硬编码",
						"中",
						"WPA全名为Wi-Fi Protected Access，是一种保护无线电脑网络（Wi-Fi）安全的系统。硬编码密码是指在程序中采用硬编码方式处理密码。由于拥有代码权限的人可以查看到密码，并使用密码访问一些不具有权限限制的系统。所以硬编码密码漏洞一旦被利用，可能会使远程攻击者获取到敏感信息，或者通过访问数据库获得管理控制权，本地用户甚至可以通过读取配置文件中的硬编码用户名和密码来执行任意代码。",
						utils.CunZaiFengXian1(len(wpahardcodeCheckers)),
						str,
						utils.CunZaiFengXian4(len(wpahardcodeCheckers), "避免采用硬编码密码，建议对密码加以模糊化，并且在外部资源文件中进行处理"),
					})
					ls1.Data = map[string]interface{}{
						"table": datas,
					}
					ls1.Children = []*utils.Document{}
					ls = append(ls, ls1)
				}
				//SVN信息泄露风险
				{
					ls1 := &utils.Document{}
					ls1.Title = "." + RPL.SVNMessage
					ls1.Template = "section221.ftl"
					ls1.Level = &utils.Level{
						ParentLevel: "2.2",
						Format:      "2.2.4",
						Order:       4,
					}
					str := "经检测以下文件中存在相应风险，详细信息如下。\n"
					for kk, v := range svnCheckers {
						_, ok := ssMap[v.Software22Id]
						if !ok {
							continue
						}
						wjm := ssMap[v.Software22Id].FileName
						if wjm == "" {
							continue
						}
						//str += fmt.Sprintf("%d.文件名：%s\n", kk+1, wjm)
						str += fmt.Sprintf("%d.文件路径：%s\n", kk+1, ssMap[v.Software22Id].FilePath)
						if utils.IsStr(v.RepoUrl) {
							str += fmt.Sprintf("   远程仓库地址：%s\n", v.RepoUrl)
						}
						if utils.IsStr(v.RepoUUID) {
							str += fmt.Sprintf("   远程仓库UUID：%s\n", v.RepoUUID)
						}
						if utils.IsStr(v.Credentials) {
							var ls []map[string]string
							err := jsoniter.Unmarshal([]byte(v.Credentials), &ls)
							if err != nil {
								fmt.Println(err)
							}
							for _, v := range ls {
								_, ok1 := v["username"]
								if ok1 {
									str += fmt.Sprintf("       泄露用户名：%s\n", v["username"])
								}
								_, ok2 := v["password"]
								if ok2 {
									str += fmt.Sprintf("       泄露密码：%s\n", v["password"])
								}
							}
						}
					}
					if str == "" || str == "经检测以下文件中存在相应风险，详细信息如下。\n" {
						str = "未检测到" + RPL.SVNMessage
					}
					var datas [][]string
					datas = append(datas, []string{
						"检测固件是否存在" + RPL.SVNMessage,
						"中",
						"Subversion项目，在打包时将Subversion使用的项目相关文件一并打包，如果文件内存在源代码、代码仓库地址、开发人员邮箱等信息，会导致信息泄露风险。",
						utils.CunZaiFengXian1(len(svnCheckers)),
						str,
						utils.CunZaiFengXian4(len(svnCheckers), "Subversion项目打包时，请将名为.svn的文件夹排除掉。"),
					})
					ls1.Data = map[string]interface{}{
						"table": datas,
					}
					ls1.Children = []*utils.Document{}
					ls = append(ls, ls1)
				}
				//源代码泄露风险
				{
					ls1 := &utils.Document{}
					ls1.Title = "." + RPL.SourceCode
					ls1.Template = "section221.ftl"
					ls1.Level = &utils.Level{
						ParentLevel: "2.2",
						Format:      "2.2.5",
						Order:       5,
					}
					str := "经检测以下文件中存在相应风险，详细信息如下。\n"
					for kk, v := range sourcecodeLeakage {
						wjm := v.FileName
						if wjm == "" {
							continue
						}
						//str += fmt.Sprintf("%d.文件名：%s\n", kk+1, wjm)
						str += fmt.Sprintf("%d.文件路径：%s\n", kk+1, v.FilePath)
						if utils.IsStr(v.SourcecodeLeakage) {
							str += fmt.Sprintf("   泄露文件类型：%s\n", v.SourcecodeLeakage)
						}
					}
					if str == "" || str == "经检测以下文件中存在相应风险，详细信息如下。\n" {
						str = "未检测到" + RPL.SourceCode
					}
					var datas [][]string
					datas = append(datas, []string{
						"检测固件是否存在" + RPL.SourceCode,
						"中",
						"固件文件中如果包含开发过程中的源代码，攻击者在获取到固件文件后，会更容易的进行分析，并挖掘文件中的风险漏洞。",
						utils.CunZaiFengXian1(len(sourcecodeLeakage)),
						str,
						utils.CunZaiFengXian4(len(sourcecodeLeakage), "固件项目打包时，请将源代码文件排除掉，比如扩展名为.js、.h、.c、.java、.py、.lua等。"),
					})
					ls1.Data = map[string]interface{}{
						"table": datas,
					}
					ls1.Children = []*utils.Document{}
					ls = append(ls, ls1)
				}
				//Git信息泄露风险
				{
					ls1 := &utils.Document{}
					ls1.Title = "." + RPL.GitMessage
					ls1.Template = "section221.ftl"
					ls1.Level = &utils.Level{
						ParentLevel: "2.2",
						Format:      "2.2.6",
						Order:       6,
					}
					str := "经检测以下文件中存在相应风险，详细信息如下。\n"
					for kk, v := range gitCheckers {
						wjm := ssMap[v.Software22Id].FileName
						if wjm == "" {
							continue
						}
						//str += fmt.Sprintf("%d.文件名：%s\n", kk+1, wjm)
						str += fmt.Sprintf("%d.文件路径：%s\n", kk+1, ssMap[v.Software22Id].FilePath)
						if utils.IsStr(v.RepoUrl) {
							str += fmt.Sprintf("   远程仓库地址：%s\n", v.RepoUrl)
						}
						if utils.IsStr(v.DeveloperEmail) {
							str += fmt.Sprintf("   开发者邮箱：%s\n", v.DeveloperEmail)
						}
						if utils.IsStr(v.DeveloperName) {
							str += fmt.Sprintf("   开发者姓名：%s\n", v.DeveloperName)
						}
						if utils.IsStr(v.DeveloperCredential) {
							str += fmt.Sprintf("   开发者的⽤户名以及密钥信息：%s\n", v.DeveloperCredential)
						}
					}
					if str == "" || str == "经检测以下文件中存在相应风险，详细信息如下。\n" {
						str = "未检测到" + RPL.GitMessage
					}
					var datas [][]string
					datas = append(datas, []string{
						"检测固件是否存在Git信息泄漏风险",
						"中",
						"固件文件中包含Git信息，攻击者会轻易获取到该信息，导致Git文件、仓库地址、开发者邮箱/用户密码等信息泄露。",
						utils.CunZaiFengXian1(len(gitCheckers)),
						str,
						utils.CunZaiFengXian4(len(gitCheckers), "固件项目打包时，请将Git信息排除掉。"),
					})
					ls1.Data = map[string]interface{}{
						"table": datas,
					}
					ls1.Children = []*utils.Document{}
					ls = append(ls, ls1)
				}
					//vi/vim信息泄露风险
					{
						ls1 := &utils.Document{}
						ls1.Title = "." + RPL.ViMMessage
						ls1.Template = "section221.ftl"
						ls1.Level = &utils.Level{
							ParentLevel: "2.2",
							Format:      "2.2.7",
							Order:       7,
						}
						str := "经检测以下文件中存在相应风险，详细信息如下。\n"
						for kk, v := range vimCheckers {
							wjm := v.FileName
							if wjm == "" {
								continue
							}
							//str += fmt.Sprintf("%d.文件名：%s\n", kk+1, wjm)
							str += fmt.Sprintf("%d.文件路径：%s\n", kk+1, v.FilePath)
						}
						if str == "" || str == "经检测以下文件中存在相应风险，详细信息如下。\n" {
							str = "未检测到" + RPL.ViMMessage
						}
						var datas [][]string
						datas = append(datas, []string{
							"检测固件是否存在vi/vim信息泄漏风险",
							"低",
							"vi是所有Unix和linux系统下标准的编辑器，vim是vi的升级版。该编辑器因具有丰富的代码补完、编译及错误跳转等功能，在程序员中被广泛使用。如果开发者在开发过程中使用了vi/vim编辑器，固件中可能会包含vi/vim文件信息，若攻击者解析固件获取到vi/vim信息，会导致其中包含的敏感信息泄露。",
							utils.CunZaiFengXian1(len(vimCheckers)),
							str,
							utils.CunZaiFengXian4(len(vimCheckers), "固件项目打包时，请将vi/vim信息排除掉。"),
						})
						ls1.Data = map[string]interface{}{
							"table": datas,
						}
						ls1.Children = []*utils.Document{}
						ls = append(ls, ls1)
					}

					//备份文件泄露风险
					//{
					//	ls1 := &utils.Document{}
					//	ls1.Title = "." + RPL.CopyFile
					//	ls1.Template = "section221.ftl"
					//	ls1.Level = &utils.Level{
					//		ParentLevel: "2.2",
					//		Format:      "2.2.8",
					//		Order:       8,
					//	}
					//	str := "经检测以下文件中存在相应风险，详细信息如下。\n"
					//	for kk, v := range bakfileChecker {
					//		wjm := v.FileName
					//		if wjm == "" {
					//			continue
					//		}
					//		//str += fmt.Sprintf("%d.文件名：%s\n", kk+1, wjm)
					//		str += fmt.Sprintf("%d.文件路径：%s\n", kk+1, v.FilePath)
					//	}
					//	if str == "" || str == "经检测以下文件中存在相应风险，详细信息如下。\n" {
					//		str = "未检测到" + RPL.CopyFile
					//	}
					//	var datas [][]string
					//	datas = append(datas, []string{
					//		"检测固件是否存在" + RPL.CopyFile,
					//		"低",
					//		"临时文件一般在下载或安装软件时创建，通常创建临时文件的程序会在完成时将其删除，但有些情况下这些文件会被保留。如果攻击者提取到固件中包含的临时文件，可能导致备份文件中的敏感信息泄露。",
					//		utils.CunZaiFengXian1(len(bakfileChecker)),
					//		str,
					//		utils.CunZaiFengXian4(len(bakfileChecker), "待定"),
					//	})
					//	ls1.Data = map[string]interface{}{
					//		"table": datas,
					//	}
					//	ls1.Children = []*utils.Document{}
					//	ls = append(ls, ls1)
					//}

				//临时文件泄露风险
				//{
				//	ls1 := &utils.Document{}
				//	ls1.Title = "." + RPL.TemporaryFile
				//	ls1.Template = "section221.ftl"
				//	ls1.Level = &utils.Level{
				//		ParentLevel: "2.2",
				//		Format:      "2.2.9",
				//		Order:       9,
				//	}
				//	str := "经检测以下文件中存在相应风险，详细信息如下。\n"
				//	for kk, v := range tmpfileChecker {
				//		wjm := v.FileName
				//		if wjm == "" {
				//			continue
				//		}
				//		//str += fmt.Sprintf("%d.文件名：%s\n", kk+1, wjm)
				//		str += fmt.Sprintf("%d.文件路径：%s\n", kk+1, v.FilePath)
				//	}
				//	if str == "" || str == "经检测以下文件中存在相应风险，详细信息如下。\n" {
				//		str = "未检测到" + RPL.TemporaryFile
				//	}
				//	var datas [][]string
				//	datas = append(datas, []string{
				//		"检测固件是否存在" + RPL.TemporaryFile,
				//		"低",
				//		"临时文件一般在下载或安装软件时创建，通常创建临时文件的程序会在完成时将其删除，但有些情况下这些文件会被保留。如果攻击者提取到固件中包含的临时文件，可能导致备份文件中的敏感信息泄露。",
				//		utils.CunZaiFengXian1(len(tmpfileChecker)),
				//		str,
				//		utils.CunZaiFengXian4(len(tmpfileChecker), "固件项目打包时，请将临时文件排除掉。"),
				//	})
				//	ls1.Data = map[string]interface{}{
				//		"table": datas,
				//	}
				//	ls1.Children = []*utils.Document{}
				//	ls = append(ls, ls1)
				//}

				//注册表泄露风险
				{
					udps := []string{
						"UDP port 135", "UDP port 137", "UDP port 138", "UDP port 445",
					}
					tcps := []string{
						"TCP port 135", "TCP port 139", "UDP port 138", "TCP port 445", "TCP port 593", "TCP port 1025", "TCP port 2745", "TCP port 3127", "TCP port 6129", "TCP port 3389",
					}
					type SoftwareAndRegistryCheckers struct {
						UDPHigh []string `json:"udp_high"`
						TCPHigh []string `json:"tcp_high"`
						UDPLow  []string `json:"udp_low"`
						TCPLow  []string `json:"tcp_low"`
					}
					sarcs := map[uint]SoftwareAndRegistryCheckers{}
					for _, v := range registryCheckers {
						var lsv SoftwareAndRegistryCheckers
						_, ok := sarcs[v.Software22Id]
						if ok {
							lsv = sarcs[v.Software22Id]
						}
						if UDPHighPort := utils.DetermineWhetherTheStringHasACorrespondingPort(v.Val, udps); UDPHighPort != "" {
							var UDPHigh []string
							if lsv.UDPHigh != nil {
								UDPHigh = lsv.UDPHigh
							}
							UDPHigh = append(UDPHigh, utils.GetStrPort(UDPHighPort))
							lsv.UDPHigh = UDPHigh
						} else if TCPHighPort := utils.DetermineWhetherTheStringHasACorrespondingPort(v.Val, tcps); TCPHighPort != "" {
							var TCPHigh []string
							if lsv.TCPHigh != nil {
								TCPHigh = lsv.TCPHigh
							}
							TCPHigh = append(TCPHigh, utils.GetStrPort(TCPHighPort))
							lsv.TCPHigh = TCPHigh
						} else {
							if strings.Contains(v.Val, "UDP port") {
								var UDPLow []string
								if lsv.UDPLow != nil {
									UDPLow = lsv.UDPLow
								}
								UDPLow = append(UDPLow, v.Val)
								lsv.UDPLow = UDPLow
							} else if strings.Contains(v.Val, "TCP port") {
								var TCPLow []string
								if lsv.TCPLow != nil {
									TCPLow = lsv.TCPLow
								}
								TCPLow = append(TCPLow, v.Val)
								lsv.TCPLow = TCPLow
							}
						}
					}
					ls1 := &utils.Document{}
					ls1.Title = "." + RPL.ResignTable
					ls1.Template = "section221.ftl"
					ls1.Level = &utils.Level{
						ParentLevel: "2.2",
						Format:      "2.2.8",
						Order:       8,
					}
					str := "经检测以下文件中存在相应风险，详细信息如下。\n"
					for k, v := range sarcs {
						wjm := ssMap[k].FileName
						if wjm == "" {
							continue
						}
						//str += fmt.Sprintf("%d.文件名：%s\n", k+1, wjm)
						str += fmt.Sprintf("%d.文件路径：%s\n", k+1, ssMap[k].FilePath)
						if v.UDPHigh != nil && len(v.UDPHigh) > 0 {
							str += fmt.Sprintf("   高风险UDP端口：%v\n", v.UDPHigh)
						}
						if v.UDPLow != nil && len(v.UDPLow) > 0 {
							str += fmt.Sprintf("   低风险UDP端口：%v\n", v.UDPLow)
						}
						if v.TCPHigh != nil && len(v.TCPHigh) > 0 {
							str += fmt.Sprintf("   高风险TCP端口：%v\n", v.TCPHigh)
						}
						if v.TCPLow != nil && len(v.TCPLow) > 0 {
							str += fmt.Sprintf("   低风险TCP端口：%v\n", v.TCPLow)
						}
					}
					if str == "" || str == "经检测以下文件中存在相应风险，详细信息如下。\n" {
						str = "未检测到" + RPL.ResignTable
					}
					var datas [][]string
					datas = append(datas, []string{
						"检测固件中是否存在注册表中端口信息泄露的风险",
						"中",
						"注册表（Registry，繁体中文版Windows操作系统称之为登录档）是MicrosoftWindows中的一个重要的数据库，用于存储系统和应用程序的设置信息。端口用于区分服务，比如用于浏览网页服务的80端口，用于FTP服务的21端口等，开放不安全端口容易被恶意攻击者扫描及利用进行攻击。",
						utils.CunZaiFengXian1(len(registryCheckers)),
						str,
						utils.CunZaiFengXian4(len(registryCheckers), "厂商在固件发布版本之前，确保关闭不必要的端口。"),
					})
					ls1.Data = map[string]interface{}{
						"table": datas,
					}
					ls1.Children = []*utils.Document{}
					ls = append(ls, ls1)
				}
				sensitiveInformationDetection.Data = map[string]interface{}{}
				sensitiveInformationDetection.Children = ls
				detailsFirmwareSecurityTestResults = append(detailsFirmwareSecurityTestResults, sensitiveInformationDetection)
			}
			//软件成分识别
			{

				softwareComponentRecognition := &utils.Document{}
				softwareComponentRecognition.Title = "." + RPL.SoftCI
				softwareComponentRecognition.Level = &utils.Level{
					ParentLevel: "2",
					Format:      "2.3",
					Order:       3,
				}
				var ls []*utils.Document
				//软件成分识别
				{
					ls1 := &utils.Document{}
					ls1.Title = "." + RPL.SoftCITwo
					ls1.Level = &utils.Level{
						ParentLevel: "2.3",
						Format:      "2.3.1",
						Order:       1,
					}
					if len(softwares) > 0 {
						ls1.Template = "section231.ftl"
						var datas [][]string

						for k, v := range softwares {
							datas = append(datas, []string{fmt.Sprintf("%d", k+1), v.SoftwareName, strings.ReplaceAll(v.SoftwareVersion, "UNKNOWN", "--"), v.FilePath})
						}
						ls1.Data = map[string]interface{}{
							"table": datas,
						}
					} else {
						ls1.Template = "cfsb.ftl"
						ls1.Data = map[string]interface{}{}
					}
					ls1.Children = []*utils.Document{}
					ls = append(ls, ls1)
				}
				softwareComponentRecognition.Data = map[string]interface{}{}
				softwareComponentRecognition.Children = ls
				detailsFirmwareSecurityTestResults = append(detailsFirmwareSecurityTestResults, softwareComponentRecognition)
			}
			//软件组件漏洞检测
			{

				softwareComponentVulnerabilityDetection := &utils.Document{}
				softwareComponentVulnerabilityDetection.Title = "." + RPL.SoftCVTwo
				softwareComponentVulnerabilityDetection.Template = "section24.ftl"
				softwareComponentVulnerabilityDetection.Level = &utils.Level{
					ParentLevel: "2",
					Format:      "2.4",
					Order:       4,
				}
				var ls, ls1 []*utils.Document
				var i int
				var datas, datas1 [][]string
				//		fmt.Println(fmt.Sprintf(len(softwareAndCveLookup))
				//		fmt.Println(fmt.Sprintf("%#v",softwareAndCveLookup)
				for k, v := range softwareAndCveLookup {
					i++
					datas = append(datas, []string{
						fmt.Sprintf("%d", i), ssMap[k].SoftwareName, strings.ReplaceAll(ssMap[k].SoftwareVersion, "UNKNOWN", "--"), fmt.Sprintf("%d", len(v.CveLookups)),
						fmt.Sprintf("%d", len(v.CveHighLookups)), fmt.Sprintf("%d", len(v.CveInLookups)), fmt.Sprintf("%d", len(v.CveLowLookups)),
					})
					{
						ls1 := &utils.Document{}
						ls1.Title = fmt.Sprintf(".%s", ssMap[k].SoftwareName)
						ls1.Template = "section241.ftl"
						ls1.Level = &utils.Level{
							ParentLevel: "2.4",
							Format:      fmt.Sprintf("2.4.%d", i),
							Order:       i,
						}
						var loopholes [][]string
						var holenames []string
						lsk := 1
						for _, v := range v.CveHighLookups {
							var cnvd cnvd
							this.edb.Table("cnvd").Where("cve_id = ?", v.CveId).First(&cnvd)

							holenames = append(holenames, fmt.Sprintf("%d.%s", lsk, v.Name))
							loopholes = append(loopholes, []string{
								v.CveId, cveIds[v.CveId].CnnvdId, utils.Cnvd(cnvd.CnvdId), strings.ReplaceAll(cveIds[v.CveId].HazardGrade, "超危", "高危"), cveIds[v.CveId].LoopholeType, cveIds[v.CveId].ThreatType,
								cveIds[v.CveId].Introduction, "不安全", fmt.Sprintf("%d.文件路径：%s", 1, ssMap[v.Software22Id].FilePath),
							})
							lsk++
						}
						for _, v := range v.CveInLookups {
							var cnvd cnvd
							this.edb.Table("cnvd").Where("cve_id = ?", v.CveId).First(&cnvd)
							holenames = append(holenames, fmt.Sprintf("%d.%s", lsk, v.Name))
							loopholes = append(loopholes, []string{
								v.CveId, cveIds[v.CveId].CnnvdId, utils.Cnvd(cnvd.CnvdId), cveIds[v.CveId].HazardGrade, cveIds[v.CveId].LoopholeType, cveIds[v.CveId].ThreatType,
								cveIds[v.CveId].Introduction, "不安全", fmt.Sprintf("%d.文件路径：%s", 1, ssMap[v.Software22Id].FilePath),
							})
							lsk++
						}
						for _, v := range v.CveLowLookups {
							var cnvd cnvd
							this.edb.Table("cnvd").Where("cve_id = ?", v.CveId).First(&cnvd)
							holenames = append(holenames, fmt.Sprintf("%d.%s", lsk, v.Name))
							loopholes = append(loopholes, []string{
								v.CveId, cveIds[v.CveId].CnnvdId, utils.Cnvd(cnvd.CnvdId), cveIds[v.CveId].HazardGrade, cveIds[v.CveId].LoopholeType, cveIds[v.CveId].ThreatType,
								cveIds[v.CveId].Introduction, "不安全", fmt.Sprintf("%d.文件路径：%s", 1, ssMap[v.Software22Id].FilePath),
							})
							lsk++
						}
						ls1.Data = map[string]interface{}{
							"table": [][]string{
								{
									fmt.Sprintf("%d", i), ssMap[k].SoftwareName, strings.ReplaceAll(ssMap[k].SoftwareVersion, "UNKNOWN", "--"), fmt.Sprintf("%d", len(v.CveLookups)),
									fmt.Sprintf("%d", len(v.CveHighLookups)), fmt.Sprintf("%d", len(v.CveInLookups)), fmt.Sprintf("%d", len(v.CveLowLookups)),
								},
							},
							"loopholes": loopholes,
							"holenames": holenames,
						}
						ls1.Children = []*utils.Document{}
						ls = append(ls, ls1)
					}
				}
				if len(softwareAndCveLookup) == 0 {
					softwareComponentVulnerabilityDetection.Template = "rjzj.ftl"
				}

				var lsss, lsss1 []model.Software22
				this.db.Where("firmware22_id = ? and software_name != ''", f.ID).Order("software_version asc").Find(&lsss)
				for _, v := range lsss {
					var cl []model.CveLookup22
					this.db.Where("software22_id = ?", v.ID).
						Find(&cl)
					if len(cl) == 0 {
						continue
					}
					lsss1 = append(lsss1, v)
				}

				for k, v := range lsss1 {
					for _, v1 := range ls {
						table := v1.Data["table"].([][]string)
						if table[0][1] == v.SoftwareName && table[0][2] == strings.ReplaceAll(v.SoftwareVersion, "UNKNOWN", "--") {
							table[0][0] = fmt.Sprintf("%d", k+1)
							v1.Data["table"] = table
							v1.Level.Format = fmt.Sprintf("2.4.%d", k+1)
							ls1 = append(ls1, v1)
							break
						}
					}
					for _, v1 := range datas {
						if v1[1] == v.SoftwareName && v1[2] == strings.ReplaceAll(v.SoftwareVersion, "UNKNOWN", "--") {
							v1[0] = fmt.Sprintf("%d", k+1)
							datas1 = append(datas1, v1)
							break
						}
					}
				}
				softwareComponentVulnerabilityDetection.Data = map[string]interface{}{
					"table": datas1,
				}
				softwareComponentVulnerabilityDetection.Children = ls1
				detailsFirmwareSecurityTestResults = append(detailsFirmwareSecurityTestResults, softwareComponentVulnerabilityDetection)
			}
			//恶意软件检测
			{

				maliciousCodeDetection := &utils.Document{}
				maliciousCodeDetection.Title = "." + RPL.BadWareTwo
				maliciousCodeDetection.Level = &utils.Level{
					ParentLevel: "2",
					Format:      "2.5",
					Order:       5,
				}
				var ls []*utils.Document
				{
					ls1 := &utils.Document{}
					ls1.Title = RPL.BadWareTwo
					ls1.Template = "section251.ftl"
					ls1.Level = &utils.Level{
						ParentLevel: "2.5",
						Format:      "2.5.1",
						Order:       1,
					}
					var datas [][]string
					for _, v := range malwareScanners {
						datas = append(datas, []string{v.VirusName, v.VirusPath})
					}
					ls1.Data = map[string]interface{}{
						"table": datas,
					}
					ls1.Children = []*utils.Document{}
					ls = append(ls, ls1)
				}
				maliciousCodeDetection.Data = map[string]interface{}{}
				maliciousCodeDetection.Children = ls
				detailsFirmwareSecurityTestResults = append(detailsFirmwareSecurityTestResults, maliciousCodeDetection)
			}
			//代码安全漏洞检测
			//{
			//	codeSecurityVulnerabilityDetection := &utils.Document{}
			//	codeSecurityVulnerabilityDetection.Title = ".代码安全漏洞检测"
			//	codeSecurityVulnerabilityDetection.Level = &utils.Level{
			//		ParentLevel: "2",
			//		Format:      "2.6",
			//		Order:       6,
			//	}
			//	//var ls []*utils.Document
			//	//var i int
			//	//for _, lsv := range []string{"CWE190", "CWE215", "CWE243", "CWE248", "CWE332", "CWE367", "CWE426", "CWE457", "CWE467", "CWE476", "CWE676", "CWE782", "CWE560"} {
			//	//	v, ok := cweMap[lsv]
			//	//	if !ok {
			//	//				fmt.Println(fmt.Sprintf(lsv)
			//	//		continue
			//	//	}
			//	//	i++
			//	//	{
			//	//		ls1 := &utils.Document{}
			//	//		ls1.Title = fmt.Sprintf(".%s", this.Abaqld1[v.Cwe.LId])
			//	//		ls1.Template = "section221.ftl"
			//	//		ls1.Level = &utils.Level{
			//	//			ParentLevel: "2.6",
			//	//			Format:      fmt.Sprintf("2.6.%d", i),
			//	//			Order:       i,
			//	//		}
			//	//		str := "经检测以下文件中存在相应风险，详细信息如下。\n"
			//	//		for kk, v1 := range v.Cwes {
			//	//			_, ok := ssMap[v1.Software22Id]
			//	//			if !ok {
			//	//				continue
			//	//			}
			//	//			wjm := ssMap[v1.Software22Id].FileName
			//	//			if wjm == "" {
			//	//				continue
			//	//			}
			//	//			//str += fmt.Sprintf("%d.文件名：%s\n", kk+1, wjm)
			//	//			str += fmt.Sprintf("%d.文件路径：%s\n", kk+1, ssMap[v1.Software22Id].FilePath)
			//	//			str += fmt.Sprintf("   漏洞详情：\n")
			//	//			//var ws []string
			//	//			//err := jsoniter.Unmarshal([]byte(v1.WS), &ws)
			//	//			//if err != nil {
			//	//			//	//		fmt.Println(err)
			//	//			//}
			//	//			//for _, v2 := range ws {
			//	//			//	str += fmt.Sprintf("      %s\n", v2)
			//	//			//}
			//	//			str += fmt.Sprintf("      %s\n", v1.Warnings)
			//	//		}
			//	//		if str == "" || str == "经检测以下文件中存在相应风险，详细信息如下。\n" {
			//	//			str = this.Abaqld2[v.Cwe.LId]
			//	//		}
			//	//		var datas [][]string
			//	//		datas = append(datas, []string{
			//	//			v.Cwe.TestPurpose,
			//	//			v.Cwe.RiskLevel,
			//	//			v.Cwe.VulnerabilityDescription,
			//	//			utils.CunZaiFengXian1(len(v.Cwes)),
			//	//			str,
			//	//			utils.CunZaiFengXian4(len(v.Cwes), v.Cwe.RepairSuggestions),
			//	//		})
			//	//		ls1.Data = map[string]interface{}{
			//	//			"table": datas,
			//	//		}
			//	//		ls1.Children = []*utils.Document{}
			//	//		ls = append(ls, ls1)
			//	//	}
			//	//}
			//	//	//安全缓解机制
			//	var ls1List []*utils.Document
			//	{
			//		ls1 := &utils.Document{}
			//		ls1.Title = ".安全缓解机制检测"
			//		ls1.Template = "section221.ftl"
			//		ls1.Level = &utils.Level{
			//			ParentLevel: "2.7",
			//			Format:      fmt.Sprintf("2.7.%d", 4),
			//			Order:       4,
			//		}
			//		str := "未完全开启安全缓解机制的文件详情如下：\n"
			//		canary := "  CANARY：\n"
			//		nx := "  NX：\n"
			//		relro := "  RELRO：\n"
			//		pie := "  PIE：\n"
			//		fortifySource := "  FORTIFY_SOURCE：\n"
			//		//exploitMitigations
			//		for _, v := range exploitMitigations {
			//			// 去除逻辑判断
			//			//ls := this.ObtainVulnerabilityMitigationStatus(v.CanaryStatus); ls != "完全开启"
			//			if v.CanaryVal != "" {
			//				canary += fmt.Sprintf("    文件路径：%s", v.CanaryVal)
			//				//canary += fmt.Sprintf("      开启状态：%s", ls)
			//			}
			//			//ls := this.ObtainVulnerabilityMitigationStatus(v.NXStatus); ls != "完全开启"
			//			if v.NXVal != "" {
			//				nx += fmt.Sprintf("    文件路径：%s", v.NXVal)
			//				//nx += fmt.Sprintf("      开启状态：%s", ls)
			//			}
			//			//ls := this.ObtainVulnerabilityMitigationStatus(v.RelroStatus); ls != "完全开启"
			//			if v.RelroVal != "" {
			//				relro += fmt.Sprintf("    文件路径：%s", v.RelroVal)
			//				//relro += fmt.Sprintf("      开启状态：%s", ls)
			//			}
			//			//ls := this.ObtainVulnerabilityMitigationStatus(v.PIEStatus); ls != "完全开启"
			//			if v.PIEVal != "" {
			//				pie += fmt.Sprintf("    文件路径：%s", v.PIEVal)
			//				//pie += fmt.Sprintf("      开启状态：%s", ls)
			//			}
			//			//ls := this.ObtainVulnerabilityMitigationStatus(v.FortifySourceStatus); ls != "完全开启"
			//			if v.FortifySourceVal != "" {
			//				fortifySource += fmt.Sprintf("    文件路径：%s", v.FortifySourceVal)
			//				//fortifySource += fmt.Sprintf("      开启状态：%s", ls)
			//			}
			//		}
			//		{
			//			if canary == "  CANARY：\n" {
			//				canary = "  CANARY： 无\n"
			//			}
			//			if nx == "  NX：\n" {
			//				nx = "  NX： 无\n"
			//			}
			//			if relro == "  RELRO：\n" {
			//				relro = "  RELRO： 无\n"
			//			}
			//			if pie == "  PIE：\n" {
			//				pie = "  PIE： 无\n"
			//			}
			//			if fortifySource == "  FORTIFY_SOURCE：\n" {
			//				fortifySource = "  FORTIFY_SOURCE： 无\n"
			//			}
			//		}
			//		str += canary + nx + relro + pie + fortifySource
			//		if mitigationMechanismVulnerabilities == 0 {
			//			str = "未检测到安全缓解机制"
			//		}
			//		var datas [][]string
			//		datas = append(datas, []string{
			//			"检测固件是否充分使用了安全缓解机制",
			//			"低",
			//			"尽管很多代码安全漏洞需要在开发、测试过程中发现并解决，但操作系统仍提供了一些安全缓解措施，用以减少这些漏洞带来的安全风险，常见的安全缓解机制如下：\n1.CANARY，即栈保护。栈溢出保护是一种缓冲区溢出攻击缓解手段。当启用栈保护后，函数开始执行的时候会先往栈里插入cookie信息，当函数真正返回的时候会验证cookie信息是否合法，如果不合法就停止程序运行。\n2.NX，No-eXecute（不可执行），即堆栈不可执行。NX（DEP）的基本原理是将数据所在内存页标识为不可执行，当程序溢出成功转入shellcode时，程序会尝试在数据页面上执行指令，此时CPU就会抛出异常，而不是去执行恶意指令。\n3.RELRO（ReadonlyRelocation），即符号重定向只读。设置符号重定向表格为只读或在程序启动时就解析并绑定所有动态符号，从而减少对GOT（GlobalOffsetTable）攻击。RELRO为”PartialRELRO”，说明我们对GOT表具有写权限。\n4.PIE，即地址随机化。内存地址随机化机制建立位置独立的可执行区域，增加利用缓冲区溢出漏洞的难度。常和NX同时工作，配合使用能有效阻止攻击者在堆栈上运行恶意代码。\n5.FORTIFY，用于检查是否存在缓冲区溢出的错误，可以辅助检测但无法检测出所有的缓冲区溢出。\n",
			//			utils.CunZaiFengXian1(mitigationMechanismVulnerabilities),
			//			str,
			//			utils.CunZaiFengXian4(mitigationMechanismVulnerabilities, "建议完全开启安全缓解机制以减小安全风险。"),
			//		})
			//		ls1.Data = map[string]interface{}{
			//			"table": datas,
			//		}
			//		ls1.Children = []*utils.Document{}
			//		ls1List = append(ls1List, ls1)
			//	}
			//	codeSecurityVulnerabilityDetection.Data = map[string]interface{}{}
			//	codeSecurityVulnerabilityDetection.Children = ls1List
			//	detailsFirmwareSecurityTestResults = append(detailsFirmwareSecurityTestResults, codeSecurityVulnerabilityDetection)
			//}
			//配置安全检测
			{
				configureSecurityDetection := &utils.Document{}
				configureSecurityDetection.Title = "." + RPL.ConfigSecurityTwo
				configureSecurityDetection.Level = &utils.Level{
					ParentLevel: "2",
					Format:      "2.6",
					Order:       6,
				}
				var ls []*utils.Document
				//硬编码风险
				{
					ls1 := &utils.Document{}
					ls1.Title = "." + RPL.HardCode
					ls1.Template = "section221.ftl"
					ls1.Level = &utils.Level{
						ParentLevel: "2.6",
						Format:      "2.6.1",
						Order:       1,
					}
					str := "经检测以下文件中存在相应风险，详细信息如下。\n"
					for kk, v := range hardcodes {
						wjm := v.Software.FileName
						if wjm == "" {
							continue
						}
						//str += fmt.Sprintf("%d.文件名：%s\n", kk+1, wjm)
						str += fmt.Sprintf("%d.文件路径：%s\n", kk+1, v.Software.FilePath)
						//str += fmt.Sprintf("   详情数据：%v\n", v.Data)
					}
					if str == "" || str == "经检测以下文件中存在相应风险，详细信息如下。\n" {
						str = "未检测到" + RPL.HardCode
					}
					var datas [][]string
					datas = append(datas, []string{
						"检测固件是否存在" + RPL.HardCode,
						"中",
						"硬编码问题是指开发者将密钥硬编码在代码、数据文件中。尽管当前开发中是用的密码学算法都是公认安全且成熟的算法，但由于算法的公开性，加密安全性依赖于密钥的保密性。如果密钥固定，对于对称密码算法，可以根据该固定密钥、公开的密钥算法和加密后的密文解密出明文，加密过程形同虚设。对于非对称密码，很大一部分在开发中的应用场景是计算签名信息用于比对。利用已知固定密钥、公开算法和明文信息，可以计算出密文，从而暴露加密信息。",
						utils.CunZaiFengXian1(len(hardcodes)),
						str,
						utils.CunZaiFengXian4(len(hardcodes), "避免使用固定的密钥串，建议通过动态计算获得密钥。"),
					})
					ls1.Data = map[string]interface{}{
						"table": datas,
					}
					ls1.Children = []*utils.Document{}
					ls = append(ls, ls1)
				}
				//SSH安全风险
				{
					ls1 := &utils.Document{}
					ls1.Title = "." + RPL.SSHSecurity
					ls1.Template = "section221.ftl"
					ls1.Level = &utils.Level{
						ParentLevel: "2.6",
						Format:      "2.6.2",
						Order:       2,
					}
					str := "经检测以下文件中存在相应风险，详细信息如下。\n"
					for kk, v := range sshs {
						wjm := v.Software.FileName
						if wjm == "" {
							continue
						}
						//str += fmt.Sprintf("%d.文件名：%s\n", kk+1, wjm)
						str += fmt.Sprintf("%d.文件路径：%s\n", kk+1, v.Software.FilePath)
						str += fmt.Sprintf("   详情数据：%v\n", v.Data)
					}
					if str == "" || str == "经检测以下文件中存在相应风险，详细信息如下。\n" {
						str = "未检测到" + RPL.SSHSecurity
					}
					var datas [][]string
					datas = append(datas, []string{
						"检测固件是否存在" + RPL.SSHSecurity,
						"中",
						"sshd_config是sshd的配置文件。在IoT设备中，如果开启SSH服务，意味着可以通过SSH协议登录远程服务器，在配置文件中PermitRootLogin如果设为yes，表明允许root用户以任何认证方式登录。当root用户存在弱密码隐患时，攻击者可借助这一问题登录设备，获得root权限。",
						utils.CunZaiFengXian1(len(sshs)),
						str,
						utils.CunZaiFengXian4(len(sshs), "\"建议不要将配置文件中的PermitRootLogin设为yes，建议设置为以下三种方式：PermitRootLogin=without-password 表示除密码外方式登录\\nPermitRootLogin=forced-commands-only仅允许使用密钥仅允许已授权的命令\\nPermitRootLogin=no禁止ssh登录\\n\""),
					})
					ls1.Data = map[string]interface{}{
						"table": datas,
					}
					ls1.Children = []*utils.Document{}
					ls = append(ls, ls1)
				}
				//FTP安全风险
				{
					ls1 := &utils.Document{}
					ls1.Title = "." + RPL.FTPSecurity
					ls1.Template = "section221.ftl"
					ls1.Level = &utils.Level{
						ParentLevel: "2.6",
						Format:      "2.6.3",
						Order:       3,
					}
					str := "经检测以下文件中存在相应风险，详细信息如下。\n"
					for kk, v := range ftps {
						wjm := v.Software.FileName
						if wjm == "" {
							continue
						}
						//str += fmt.Sprintf("%d.文件名：%s\n", kk+1, wjm)
						str += fmt.Sprintf("%d.文件路径：%s\n", kk+1, v.Software.FilePath)
						str += fmt.Sprintf("   详情数据：%v\n", v.Data)
					}
					if str == "" || str == "经检测以下文件中存在相应风险，详细信息如下。\n" {
						str = "未检测到" + RPL.FTPSecurity
					}
					var datas [][]string
					datas = append(datas, []string{
						"检测固件是否存在" + RPL.FTPSecurity,
						"中",
						"和SSH配置类似，FTP一些上传、下载、访问等敏感的操作可以通过设置配置文件中的用户过滤信息来建立安全策略。由于FTP允许匿名访问，如果开启了匿名访问anonymous_enable=yes，则可能存在无需密码访问no_anon_password=yes、提供写权限write_enable=yes、上传权限anon_upload_enable=yes、删除或重命名权限anon_other_write_enable=yes等敏感行为。所以，如果没有设置任何用户过滤配置，则在没有审核用户权限的情况下就允许用户进行敏感操作，是存在安全风险的。",
						utils.CunZaiFengXian1(len(ftps)),
						str,
						utils.CunZaiFengXian4(len(ftps), "禁止匿名用户登录，启用身份认证登录。"),
					})
					ls1.Data = map[string]interface{}{
						"table": datas,
					}
					ls1.Children = []*utils.Document{}
					ls = append(ls, ls1)
				}
					//自启动服务检测
					//{
					//	//selfstartingChecker
					//	ls1 := &utils.Document{}
					//	ls1.Title = "." + RPL.SelfStartSR
					//	ls1.Template = "zqd.ftl"
					//	ls1.Level = &utils.Level{
					//		ParentLevel: "2.6",
					//		Format:      "2.6.4",
					//		Order:       4,
					//	}
					//	str := "经检测以下文件中存在相应风险，详细信息如下。\n"
					//	for kk, v := range selfstartingChecker {
					//		wjm := ssMap[v.Software22Id].FileName
					//		if wjm == "" {
					//			continue
					//		}
					//		//str += fmt.Sprintf("%d.文件名：%s\n", kk+1, wjm)
					//		str += fmt.Sprintf("%d.文件路径：%s\n", kk+1, ssMap[v.Software22Id].FilePath)
					//		str += "结果详情\n"
					//		if utils.IsStr(v.Inittab) {
					//			str += fmt.Sprintf("   自启动类型：%s\n", v.Inittab)
					//		}
					//		if utils.IsStr(v.Script) {
					//			var ls []string
					//			err := jsoniter.Unmarshal([]byte(v.Script), &ls)
					//			if err != nil {
					//				fmt.Println(err)
					//			}
					//			if len(ls) > 0 {
					//				str += fmt.Sprintf("   Script：\n")
					//				for _, v := range ls {
					//					str += fmt.Sprintf("      %s\n", v)
					//				}
					//			}
					//		}
					//		if utils.IsStr(v.Inittab) {
					//			var ls []string
					//			err := jsoniter.Unmarshal([]byte(v.Inittab), &ls)
					//			if err != nil {
					//				fmt.Println(err)
					//			}
					//			if len(ls) > 0 {
					//				str += fmt.Sprintf("   Inittab：\n")
					//				for _, v := range ls {
					//					str += fmt.Sprintf("      %s\n", v)
					//				}
					//			}
					//		}
					//		if utils.IsStr(v.ExecStart) {
					//			var ls []string
					//			err := jsoniter.Unmarshal([]byte(v.ExecStart), &ls)
					//			if err != nil {
					//				fmt.Println(err)
					//			}
					//			if len(ls) > 0 {
					//				str += fmt.Sprintf("   ExecStart：\n")
					//				for _, v := range ls {
					//					str += fmt.Sprintf("      %s\n", v)
					//				}
					//			}
					//		}
					//		if utils.IsStr(v.Exec) {
					//			var ls []string
					//			err := jsoniter.Unmarshal([]byte(v.Exec), &ls)
					//			if err != nil {
					//				fmt.Println(err)
					//			}
					//			if len(ls) > 0 {
					//				str += fmt.Sprintf("   Exec：\n")
					//				for _, v := range ls {
					//					str += fmt.Sprintf("      %s\n", v)
					//				}
					//			}
					//		}
					//		if utils.IsStr(v.Description) {
					//			var ls []string
					//			err := jsoniter.Unmarshal([]byte(v.Description), &ls)
					//			if err != nil {
					//				fmt.Println(err)
					//			}
					//			if len(ls) > 0 {
					//				str += fmt.Sprintf("   Description：\n")
					//				for _, v := range ls {
					//					str += fmt.Sprintf("      %s\n", v)
					//				}
					//			}
					//		}
					//		if utils.IsStr(v.PreStart) {
					//			var ls []string
					//			err := jsoniter.Unmarshal([]byte(v.PreStart), &ls)
					//			if err != nil {
					//				fmt.Println(err)
					//			}
					//			if len(ls) > 0 {
					//				str += fmt.Sprintf("   PreStart：\n")
					//				for _, v := range ls {
					//					str += fmt.Sprintf("      %s\n", v)
					//				}
					//			}
					//		}
					//		//str += fmt.Sprintf("   详情数据：%v",v.Data)
					//	}
					//	if str == "" || str == "经检测以下文件中存在相应风险，详细信息如下。\n" {
					//		str = "未检测到"+RPL.SelfStartSR
					//	}
					//	str = strings.ReplaceAll(str, "\t", "")
					//	var advice string
					//	if utils.CunZaiFengXian1(len(selfstartingChecker)) == "安全" {
					//		advice = "无"
					//	} else {
					//		advice = "请开发者自查验证相关"+RPL.SelfStartSR
					//	}
					//	var datas [][]string
					//	datas = append(datas, []string{
					//		"检测固件中所包含的相关"+RPL.SelfStartSR,
					//		"提示信息",
					//		//"无",
					//		utils.CunZaiFengXian1(len(selfstartingChecker)),
					//		str,
					//		advice,
					//	})
					//	fmt.Println(fmt.Sprintf("%#v", datas))
					//	ls1.Data = map[string]interface{}{
					//		"table": datas,
					//	}
					//	ls1.Children = []*utils.Document{}
					//	ls = append(ls, ls1)
					//}
				//安全缓解机制
				{
					ls1 := &utils.Document{}
					ls1.Title = "." + RPL.SafetyMM
					ls1.Template = "section221.ftl"
					ls1.Level = &utils.Level{
						ParentLevel: "2.6",
						Format:      fmt.Sprintf("2.6.%d", 4),
						Order:       4,
					}
					str := "未完全开启安全缓解机制的文件详情如下：\n"
					canary := "  CANARY：\n"
					nx := "  NX：\n"
					relro := "  RELRO：\n"
					pie := "  PIE：\n"
					fortifySource := "  FORTIFY_SOURCE：\n"
					//exploitMitigations
					for _, v := range exploitMitigations {
						// 去除逻辑判断
						//ls := this.ObtainVulnerabilityMitigationStatus(v.CanaryStatus); ls != "完全开启"
						if v.CanaryVal != "" {
							canary += fmt.Sprintf("    文件路径：%s", v.CanaryVal)
							//canary += fmt.Sprintf("      开启状态：%s", ls)
						}
						//ls := this.ObtainVulnerabilityMitigationStatus(v.NXStatus); ls != "完全开启"
						if v.NXVal != "" {
							nx += fmt.Sprintf("    文件路径：%s", v.NXVal)
							//nx += fmt.Sprintf("      开启状态：%s", ls)
						}
						//ls := this.ObtainVulnerabilityMitigationStatus(v.RelroStatus); ls != "完全开启"
						if v.RelroVal != "" {
							relro += fmt.Sprintf("    文件路径：%s", v.RelroVal)
							//relro += fmt.Sprintf("      开启状态：%s", ls)
						}
						//ls := this.ObtainVulnerabilityMitigationStatus(v.PIEStatus); ls != "完全开启"
						if v.PIEVal != "" {
							pie += fmt.Sprintf("    文件路径：%s", v.PIEVal)
							//pie += fmt.Sprintf("      开启状态：%s", ls)
						}
						//ls := this.ObtainVulnerabilityMitigationStatus(v.FortifySourceStatus); ls != "完全开启"
						if v.FortifySourceVal != "" {
							fortifySource += fmt.Sprintf("    文件路径：%s", v.FortifySourceVal)
							//fortifySource += fmt.Sprintf("      开启状态：%s", ls)
						}
					}
					{
						if canary == "  CANARY：\n" {
							canary = "  CANARY： 无\n"
						}
						if nx == "  NX：\n" {
							nx = "  NX： 无\n"
						}
						if relro == "  RELRO：\n" {
							relro = "  RELRO： 无\n"
						}
						if pie == "  PIE：\n" {
							pie = "  PIE： 无\n"
						}
						if fortifySource == "  FORTIFY_SOURCE：\n" {
							fortifySource = "  FORTIFY_SOURCE： 无\n"
						}
					}
					switch reportType {
					case "ANHEN":
						str += canary + nx + pie
					case "ANYU":
						str += canary + relro + pie + fortifySource
					default:
						str += canary + nx + relro + pie + fortifySource
					}
					if mitigationMechanismVulnerabilities == 0 {
						str = "未检测到安全缓解机制"
					}
					var datas [][]string
					datas = append(datas, []string{
						"检测固件是否充分使用了安全缓解机制",
						"低",
						"尽管很多代码安全漏洞需要在开发、测试过程中发现并解决，但操作系统仍提供了一些安全缓解措施，用以减少这些漏洞带来的安全风险，常见的安全缓解机制如下：\n1.CANARY，即栈保护。栈溢出保护是一种缓冲区溢出攻击缓解手段。当启用栈保护后，函数开始执行的时候会先往栈里插入cookie信息，当函数真正返回的时候会验证cookie信息是否合法，如果不合法就停止程序运行。\n2.NX，No-eXecute（不可执行），即堆栈不可执行。NX（DEP）的基本原理是将数据所在内存页标识为不可执行，当程序溢出成功转入shellcode时，程序会尝试在数据页面上执行指令，此时CPU就会抛出异常，而不是去执行恶意指令。\n3.RELRO（ReadonlyRelocation），即符号重定向只读。设置符号重定向表格为只读或在程序启动时就解析并绑定所有动态符号，从而减少对GOT（GlobalOffsetTable）攻击。RELRO为”PartialRELRO”，说明我们对GOT表具有写权限。\n4.PIE，即地址随机化。内存地址随机化机制建立位置独立的可执行区域，增加利用缓冲区溢出漏洞的难度。常和NX同时工作，配合使用能有效阻止攻击者在堆栈上运行恶意代码。\n5.FORTIFY，用于检查是否存在缓冲区溢出的错误，可以辅助检测但无法检测出所有的缓冲区溢出。\n",
						utils.CunZaiFengXian1(mitigationMechanismVulnerabilities),
						str,
						utils.CunZaiFengXian4(mitigationMechanismVulnerabilities, "建议完全开启安全缓解机制以减小安全风险。"),
					})
					ls1.Data = map[string]interface{}{
						"table": datas,
					}
					ls1.Children = []*utils.Document{}
					ls = append(ls, ls1)
				}

				configureSecurityDetection.Data = map[string]interface{}{}
				configureSecurityDetection.Children = ls
				detailsFirmwareSecurityTestResults = append(detailsFirmwareSecurityTestResults, configureSecurityDetection)
			}

			ls.Data = map[string]interface{}{}
			ls.Children = detailsFirmwareSecurityTestResults
			all = append(all, ls)
		}
		document.Data = map[string]interface{}{}
		document.Children = all
		s.Document = document
	}

	lsf := fmt.Sprintf("%s-固件安全分析报告(%s).docx", f.Name, f.CreatedAt.Format("2006-01-02-15-04-05"))
	fname := f.MediaAddress + lsf
	pdf := fname[:strings.LastIndex(fname, ".docx")] + ".pdf"

	fmt.Println(fmt.Sprintf(fname))
	fmt.Println(fmt.Sprintf(pdf))
	s.OutputPath = fname
	s.PdfPath = pdf
	_ = os.MkdirAll(filepath.Dir(s.OutputPath), os.ModePerm)

	jsonByte, _ := jsoniter.Marshal(s)
	jsonFilePath := fmt.Sprintf("./media/gujian/%s.json", uuid.New())
	_ = os.MkdirAll(filepath.Dir(jsonFilePath), os.ModePerm)
	err = ioutil.WriteFile(jsonFilePath, jsonByte, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	str1, err := reportTemplateRelated.RunCommand(fmt.Sprintf("java -jar %s %s %s", "resources/report/ReportGenerator.jar",
		utils.AbsPath("resources/gujianbaogao/fengxi"), utils.AbsPath(jsonFilePath)))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fmt.Sprintf(str1))

	//err = tool.WordToPdf(fname, pdf)
	//if err != nil {
	//	fmt.Println(err)
	//}
}

//需要更改的检测项
type ReportDesignWord struct {
	SIDTitle          string `json:"sid_title"`           // 敏感信息检测标题
	SIDTitleTwo       string `json:"sid_title_two"`       // 敏感信息检测标题2
	WeakPWD           string `json:"weak_pwd"`            // 弱密码检测
	KeyStore          string `json:"key_store"`           // 证书文件和密钥泄露风险
	WpaPWD            string `json:"wpa_pwd"`             // WPA密码硬编码检测
	SVNMessage        string `json:"svn_message"`         // SVN信息泄露风险
	SourceCode        string `json:"source_code"`         // 源代码泄露风险
	GitMessage        string `json:"git_message"`         // Git信息泄露风险
	ViMMessage        string `json:"vi_m_message"`        // vi/vim信息泄露风险
	CopyFile          string `json:"copy_file"`           // 备份文件泄露风险
	TemporaryFile     string `json:"temporary_file"`      // 临时文件泄露风险
	ResignTable       string `json:"resign_table"`        // 注册表泄露风险
	SoftCI            string `json:"soft_ci"`             // 软件成分识别
	SoftCITwo         string `json:"soft_ci_two"`         // 软件成分识别2
	SoftCV            string `json:"soft_cv"`             // 软件漏洞1
	SoftCVTwo         string `json:"soft_cv_two"`         // 软件漏洞2
	BadWare           string `json:"bad_ware"`            // 恶意软件检测
	BadWareTwo        string `json:"bad_ware_two"`        // 恶意软件检测2
	ConfigSecurity    string `json:"config_security"`     // 配置安全检测
	ConfigSecurityTwo string `json:"config_security_two"` // 配置安全检测2
	HardCode          string `json:"hard_code"`           // 硬编码风险
	SSHSecurity       string `json:"ssh_security"`        // SSH安全风险
	FTPSecurity       string `json:"ftp_security"`        //FTP安全风险
	SelfStartSR       string `json:"self_start_sr"`       // 自启动服务风险
	SafetyMM          string `json:"safety_mm"`           // 安全缓解机制检测
}

func reportWords(reportType string) ReportDesignWord {
	var rdw ReportDesignWord
	switch reportType {
	case "ANHEN":
		rdw = ReportDesignWord{
			SIDTitle:          "敏感信息泄漏风险（10项）",
			SIDTitleTwo:       "敏感信息泄漏风险",
			WeakPWD:           "常用字典密码检测",
			KeyStore:          "密钥和证书文件泄露",
			WpaPWD:            "WPA硬编码风险",
			SVNMessage:        "SVN信息泄露",
			SourceCode:        "代码泄露",
			GitMessage:        "Gitlab/Github信息泄露风险",
			ViMMessage:        "",
			CopyFile:          "备份文件泄露风险",
			TemporaryFile:     "",
			ResignTable:       "系统注册表泄露风险",
			SoftCI:            "第三方组件风险感知",
			SoftCITwo:         "第三方组件分析",
			SoftCV:            "软件组件安全威胁（1项）",
			SoftCVTwo:         "软件组件安全威胁",
			BadWare:           "恶意应用识别（1项）",
			BadWareTwo:        "恶意应用识别",
			ConfigSecurity:    "固件配置安全（5项）",
			ConfigSecurityTwo: "固件配置安全",
			HardCode:          "密钥硬编码风险",
			SSHSecurity:       "SSH协议安全风险",
			FTPSecurity:       "FTP协议安全风险",
			SelfStartSR:       "自启动风险检测",
			SafetyMM:          "安全缓解机制检测",
		}
	case "ANYU":
		rdw = ReportDesignWord{
			SIDTitle:          "敏感信息泄漏检测（10项）",
			SIDTitleTwo:       "敏感信息泄漏检测",
			WeakPWD:           "弱口令检测",
			KeyStore:          "证书文件和密钥泄露漏洞",
			WpaPWD:            "密码硬编码漏洞",
			SVNMessage:        "SVN泄露漏洞",
			SourceCode:        "代码泄露漏洞",
			GitMessage:        "Git泄露漏洞",
			ViMMessage:        "vi\\vim泄露漏洞",
			CopyFile:          "",
			TemporaryFile:     "",
			ResignTable:       "系统注册表泄露漏洞",
			SoftCI:            "第三方组件安全",
			SoftCITwo:         "第三方组件识别",
			SoftCV:            "软件组件安全风险（1项）",
			SoftCVTwo:         "软件组件安全风险",
			BadWare:           "风险软件检测（1项）",
			BadWareTwo:        "风险软件检测",
			ConfigSecurity:    "配置安全风险（5项）",
			ConfigSecurityTwo: "配置安全风险",
			HardCode:          "硬编码风险",
			SSHSecurity:       "SSH安全风险",
			FTPSecurity:       "FTP安全风险",
			SelfStartSR:       "",
			SafetyMM:          "安全缓解机制检测",
		}
	default:
		rdw = ReportDesignWord{
			SIDTitle:          "敏感信息检测（10项）",
			SIDTitleTwo:       "敏感信息检测",
			WeakPWD:           "弱密码检测",
			KeyStore:          "证书文件和密钥泄露风险",
			WpaPWD:            "WPA密码硬编码检测",
			SVNMessage:        "SVN信息泄露风险",
			SourceCode:        "源代码泄露风险",
			GitMessage:        "Git信息泄露风险",
			ViMMessage:        "vi/vim信息泄露风险",
			CopyFile:          "备份文件泄露风险",
			TemporaryFile:     "临时文件泄露风险",
			ResignTable:       "注册表泄露风险",
			SoftCI:            "软件成分识别",
			SoftCITwo:         "软件成分识别",
			SoftCV:            "软件组件漏洞检测（1项）",
			SoftCVTwo:         "软件组件漏洞检测",
			BadWare:           "恶意软件检测（1项）",
			BadWareTwo:        "恶意软件检测",
			ConfigSecurity:    "配置安全检测（5项）",
			ConfigSecurityTwo: "配置安全检测",
			HardCode:          "硬编码风险",
			SSHSecurity:       "SSH安全风险",
			FTPSecurity:       "FTP安全风险",
			SelfStartSR:       "自启动服务风险检测",
			SafetyMM:          "安全缓解机制检测",
		}
	}
	return rdw
}
