package logger

import "os"

const (
	PermFileMode   os.FileMode = 0666
	TimeFormat     string      = "2006-01-02T15:04:05.000Z0700"
	MaxStackLength int         = 40
)

type LogSystem string

const (
	Service  = LogSystem("SERVICE")
	Database = LogSystem("DATABASE")
)

type logTypes string

const (
	AppLog = logTypes("AppLog")
)

type LogLevel string

const (
	ALL   = LogLevel("ALL")
	TRACE = LogLevel("TRACE")
	DEBUG = LogLevel("DEBUG")
	INFO  = LogLevel("INFO")
	WARN  = LogLevel("WARN")
	ERROR = LogLevel("ERROR")
	FATAL = LogLevel("FATAL")
	OFF   = LogLevel("OFF")
)

func (level LogLevel) Integer() int {
	var levelLog int

	switch level {
	case ALL:
		levelLog = 7
	case TRACE:
		levelLog = 6
	case DEBUG:
		levelLog = 5
	case INFO:
		levelLog = 4
	case WARN:
		levelLog = 3
	case ERROR:
		levelLog = 2
	case FATAL:
		levelLog = 1
	case OFF:
		levelLog = 0
	default:
		levelLog = 0
	}

	return levelLog
}

type monitorTypes string

func (monitorType monitorTypes) isMonitorType() monitorTypes {
	return monitorType
}

type LogMonitorType interface {
	isMonitorType() monitorTypes
}

type PatternLogger struct {
	ApplicationName string
	ProductName     string
	SourceSystem    LogSystem
	TargetSystem    LogSystem
	SetLogger       SetLogger
	SourceHostName  string
	Level           LogLevel
}

type SetLogger struct {
	IsJson    bool
	WriteFile bool
	Path      string
	FileName  string
}

type logMessageBean struct {
	Timestamp       string      `json:"@timestamp"`
	ApplicationName string      `json:"@suffix,omitempty"`
	LogType         string      `json:"logType,omitempty"`
	CorrelationID   string      `json:"correlationID,omitempty"`
	Level           LogLevel    `json:"level,omitempty"`
	Message         string      `json:"message,omitempty"`
	StackTrace      interface{} `json:"stackTrace,omitempty"`
	SourceSystem    LogSystem   `json:"sourceSystem,omitempty"`
	SourceHostName  string      `json:"sourceHostName,omitempty"`
	TargetSystem    LogSystem   `json:"targetSystem,omitempty"`
	ElapsedTime     int64       `json:"elapsedTime,omitempty"`
	ResponseCode    string      `json:"responseCode,omitempty"`
}
type callerInfoBean struct {
	ClassName  string `json:"className"`
	MethodName string `json:"methodName"`
	FileName   string `json:"fileName"`
}

type stackTrace struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Function string `json:"function"`
}
