package hecmcc

import (
	"encoding/csv"
	"fmt"
	"github.com/robfig/cron/v3"
	"gopkg.in/natefinch/lumberjack.v2"
	"kmesh2/utils/logger"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	DateTimeFormat        = "20060102150405"
	defaultLogsBufferSize = 200
)
const (
	LogTypeLogin   = "1"
	LogTypeOperate = "2"
)
const (
	LoginTypeLogin  = "登入"
	LoginTypeLogout = "登出"
)
const (
	ResultSuccess = "成功"
	ResultFailure = "失败"
)
const (
	OperateTypeAdd             = "1"  //增加
	OperateTypeDelete          = "2"  //删除
	OperateTypeUpdate          = "3"  //修改
	OperateTypeQuery           = "4"  //查询
	OperateTypeOther           = "5"  //其他
	OperateTypeCopy            = "6"  //复制
	OperateTypeUpload          = "7"  //上传
	OperateTypeDownload        = "8"  //下载
	OperateTypeImport          = "9"  //导入
	OperateTypeExport          = "10" //导出
	OperateTypeSessionTracking = "11" //会话跟踪
	OperateTypeLogin           = "12" //登录
	OperateTypeLogout          = "13" //登出
)
const (
	ModuleEvaluationTask      = "检测任务"
	ModuleEvaluationTemplate  = "检测模板"
	ModuleCustomerManagement  = "客户管理"
	ModuleDetectionManagement = "检测管理"
	ModuleSystemManagement    = "系统管理"
	ModuleLogManagement       = "日志管理"
	ModuleAccountMange        = "账号管理"
	ModuleRoleMange           = "角色管理"
	ModuleExportReport        = "报告导出"
	ModuleDataStatistic       = "数据统计"
)

type ApiLog struct {
	UniqueId      string    //操作日志的唯一标识
	OperateTime   time.Time //操作信息的时间，精确到秒，以YYYYMMDDHHMMSS的格式上报
	LoginAccount  string    //系统登录帐号
	ServerIPV4    string    //登录系统IPv4地址（写到配置文件中）
	ServerIPV6    string    //IPv4与IPv6至少填写一个（写到配置文件中）
	ClientIPV4    string    //登录系统的客户端IPv4地址
	ClientIPV6    string    //IPv4与IPv6至少填写一个
	LogType       string    //日志类型,1：登录日志  2：操作日志
	SystemName    string    //系统名称,所在业务系统名称，（客户后续提供具体名称，属于固定值）
	ModuleName    string    //操作模块,系统的对应模块名称，例如：日报导出 注：日志类型为操作日志时必填
	OperateType   string    //操作类型,填写操作类型编号,日志类型为操作日志时必填
	OperateDetail string    //操作内容,对操作过程的详细描述 注：日志类型为操作日志时必填
	OperateResult string    //对系统模块的操作结果（成功或者失败）,日志类型为操作日志时必填
	LoginType     string    //1：登入 2：登出,日志类型为登录日志时，必填
	LoginResult   string    //响应信息,系统对登录状态的响应结果；成功、失败,当操作类型为登录类操作或日志类型为登录日志时，该字段为必填
}

func (r *ApiLog) Format() []string {
	result := make([]string, 0, 15)
	operateDetail := r.OperateDetail
	if strings.ContainsAny(operateDetail, "\r\n") {
		operateDetail = strconv.Quote(operateDetail)
	}
	result = append(result, r.UniqueId, r.OperateTime.Format(DateTimeFormat),
		r.LoginAccount, r.ServerIPV4, r.ServerIPV6, r.ClientIPV4, r.ClientIPV6,
		r.LogType, r.SystemName, r.ModuleName, r.OperateType, operateDetail,
		r.OperateResult, r.LoginType, r.LoginResult,
	)
	return result
}

type ApiLogger struct {
	filenameNoExt    string
	rotateFileWriter *lumberjack.Logger
	recorder         *csv.Writer
	lock             sync.Mutex
	logsBuffer       []ApiLog
	fileUploader     *FtpUploader
	syncTimer        *cron.Cron
}

func NewApiLogger(logDir string, filenameNoExt string, fileUploader *FtpUploader) *ApiLogger {
	apiLogger := new(ApiLogger)
	apiLogger.filenameNoExt = filenameNoExt
	apiLogger.rotateFileWriter = &lumberjack.Logger{
		Filename:  fmt.Sprintf("%s/%s.csv", logDir, filenameNoExt),
		MaxSize:   512,
		MaxAge:    7,
		LocalTime: true,
	}
	apiLogger.recorder = csv.NewWriter(apiLogger.rotateFileWriter)
	apiLogger.recorder.Comma = '|'
	apiLogger.logsBuffer = make([]ApiLog, 0, defaultLogsBufferSize)
	apiLogger.fileUploader = fileUploader
	apiLogger.syncTimer = cron.New()
	apiLogger.startFileSync()
	return apiLogger
}

func (l *ApiLogger) Record(apiLog ApiLog) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.logsBuffer = append(l.logsBuffer, apiLog)
	if len(l.logsBuffer) < defaultLogsBufferSize {
		return nil
	}
	return l.flush()
}
func (l *ApiLogger) Flush() error {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.flush()
}
func (l *ApiLogger) startFileSync() {
	if _, err := l.syncTimer.AddFunc("0 0 * * ?", func() {
		if err := l.Flush(); err != nil {
			logger.Error(err)
		}
		if err := l.UploadFile(); err != nil {
			logger.Error(err)
		}
	}); err != nil {
		logger.Error(err)
	}
	if _, err := l.syncTimer.AddFunc(fmt.Sprintf("@every %ds", 1), func() {
		if err := l.Flush(); err != nil {
			logger.Error(err)
		}
	}); err != nil {
		logger.Error(err)
	}
	l.syncTimer.Start()
}
func (l *ApiLogger) flush() error {
	if len(l.logsBuffer) == 0 {
		return nil
	}
	var errs error
	for _, record := range l.logsBuffer {
		if err := l.recorder.Write(record.Format()); err != nil {
			if errs == nil {
				errs = err
			} else {
				errs = fmt.Errorf("%v;%w", errs, err)
			}
		}
	}
	l.logsBuffer = l.logsBuffer[0:0]
	l.recorder.Flush()
	return errs
}
func (l *ApiLogger) UploadFile() error {
	if err := l.fileUploader.Upload(nil); err != nil {
		return err
	}
	return nil
}
