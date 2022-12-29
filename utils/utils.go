package utils

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Cve struct {
	CveId        string `json:"cve_id"`
	Name         string `json:"name"`
	CnnvdId      string `json:"cnnvd_id"`
	HazardGrade  string `json:"hazard_grade"`
	LoopholeType string `json:"loophole_type"`
	Introduction string `json:"introduction" gorm:"type:longtext"`
	Notice       string `json:"notice" gorm:"type:longtext"`
	ThreatType   string `json:"threat_type"`
}
type References struct {
	Type     string                 `json:"type"`
	Name     string                 `json:"name"`
	Property map[string]interface{} `json:"property"`
}
type TitleStyle struct {
	Data       map[string]interface{} `json:"data"`
	LevelStyle map[string]interface{} `json:"levelStyle"`
}
type Header struct {
	FirstTemplate   string `json:"firstTemplate"`
	DefaultTemplate string `json:"defaultTemplate"`
	EvenTemplate    string `json:"evenTemplate"`
	Data            struct {
		HeaderText string        `json:"header_text"`
		References []*References `json:"references"`
		FooterText string        `json:"footer_text"`
	} `json:"data"`
}
type Toc struct {
	Template string                 `json:"template"`
	Data     map[string]interface{} `json:"data"`
}
type Cover struct {
	Template string                 `json:"template"`
	Data     map[string]interface{} `json:"data"`
}
type About struct {
	Template string `json:"template"`
	Data     struct {
		IsSHow       bool          `json:"is_show"` // 是否显示公司介绍
		Website      string        `json:"website"`
		Title        string        `json:"title"`
		Introduction string        `json:"introduction"`
		References   []*References `json:"references"`
	} `json:"data"`
}
type Level struct {
	ParentLevel string `json:"parentLevel"`
	Format      string `json:"format"`
	Order       int    `json:"order"`
}
type Document struct {
	Title    string                 `json:"title"`
	Method   string                 `json:"method"`
	Level    *Level                 `json:"level"`
	Children []*Document            `json:"children"`
	Template string                 `json:"template"`
	Data     map[string]interface{} `json:"data"`
}

type Report struct {
	InitTemplate string      `json:"initTemplate"`
	Header       *Header     `json:"header"`
	Footer       Header      `json:"footer"`
	TitleStyle   *TitleStyle `json:"titleStyle"`
	Toc          *Toc        `json:"toc"`
	Document     *Document   `json:"document"`
	Cover        *Cover      `json:"cover"`
	About        *About      `json:"about"`
	OutputPath   string      `json:"outputPath"`
	PdfPath      string      `json:"pdfPath"`
}

func AbsPath(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println(err)
		return path
	}
	if strings.HasSuffix(path, "/") {
		absPath += "/"
	}
	return absPath
}

func CunZaiFengXian(i int) string {
	if i == 0 {
		return "安全"
	}
	return fmt.Sprintf("存在风险（发现%d处）", i)
}

func CunZaiFengXian2(i int) string {
	if i == 0 {
		return "安全"
	}
	return fmt.Sprintf("发现%d处", i)
}

func IsStr(str string) bool {
	if str != "" && str != "<nil>" && str != "null" && str != "nil" {
		return true
	}
	return false
}
func PoJi(str string) string {
	if str == "false" {
		return "破解失败"
	} else if str == "true" {
		return "破解成功"
	}
	return str
}
func CunZaiFengXian1(i int) string {
	if i == 0 {
		return "安全"
	}
	return fmt.Sprintf("存在风险（发现%d处）", i)
}
func CunZaiFengXian4(i int, s string) string {
	if i == 0 {
		return "无"
	}
	return s
}


func DetermineWhetherTheStringHasACorrespondingPort(str string, strs []string) string {
	for _, v := range strs {
		if strings.Contains(str, v) {
			return v
		}
	}
	return ""
}
func GetStrPort(str string) string {
	x := strings.LastIndex(str, " ")
	if x == -1 {
		return ""
	}
	return str[x+1:]
}
func Cnvd(str string) string {
	if str == "" {
		return "/"
	}
	return str
}