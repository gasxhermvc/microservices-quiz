package logger

import (
	"encoding/json"
	"fmt"
	"go/build"
	logger "log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func (log *PatternLogger) initLogger(productName string, sourceSystem LogSystem, targetSystem LogSystem, applicationName string) *PatternLogger {
	log.SetLogger.IsJson = true
	log.ProductName = productName
	log.Level = INFO
	log.SourceSystem = sourceSystem
	log.TargetSystem = targetSystem
	log.SourceHostName, _ = os.Hostname()
	log.ApplicationName = applicationName

	return log
}

func (log *PatternLogger) InitLogger(productName string, appName string, sourceSystem LogSystem, targetSystem LogSystem) *PatternLogger {
	return log.initLogger(productName, sourceSystem, targetSystem, appName)
}

func (log *PatternLogger) EnableFileLogger(path string, fileName string) {
	log.SetLogger.WriteFile = true
	log.SetLogger.Path = path
	log.SetLogger.FileName = fileName
}

func (log *PatternLogger) log(correlationID string, level LogLevel, message string, args ...interface{}) {
	messageBean := new(logMessageBean)
	messageBean.Timestamp = time.Now().Format(TimeFormat)
	messageBean.ApplicationName = log.ApplicationName
	messageBean.LogType = "App"
	messageBean.CorrelationID = correlationID
	messageBean.Level = level
	messageBean.Message = message
	messageBean.StackTrace = args[0]
	messageBean.SourceHostName = log.SourceHostName
	messageBean.SourceSystem = log.SourceSystem
	messageBean.TargetSystem = log.TargetSystem
	log.WriteLogger(messageBean)
}

func (log *PatternLogger) WriteLogger(messageBean *logMessageBean) {
	if !log.SetLogger.IsJson {
		if !log.SetLogger.WriteFile {
			_, _ = fmt.Fprintln(os.Stdout, log.logStringPattern(messageBean))
		} else {
			log.createLogFile(log.logStringPattern(messageBean))
		}
	} else {
		if !log.SetLogger.WriteFile {
			enc := json.NewEncoder(os.Stdout)
			enc.SetEscapeHTML(false)

			if "dev" == os.Args[1] {
				enc.SetIndent("", "    ")
			}
			_ = enc.Encode(messageBean)
		} else {
			marshalStruct, _ := json.Marshal(&messageBean)
			log.createLogFile(string(marshalStruct))
		}
	}
}

func (log *PatternLogger) logStringPattern(bean *logMessageBean) (logString string) {
	logString = fmt.Sprintf(`timestamp=%s|logType=%s|correlationID=%s|level=%v|message=%s|stackTrace=%v|sourceSystem=%s|sourceHostName=%s|targetSystem=%s|elapsedTime=%d|responseCode=%s`,
		bean.Timestamp, bean.LogType, bean.CorrelationID, bean.Level, bean.Message, bean.StackTrace,
		bean.SourceSystem, bean.SourceHostName, bean.TargetSystem, bean.ElapsedTime, bean.ResponseCode)
	return
}

func (log *PatternLogger) createLogFile(logging string) {
	file, _ := os.OpenFile(build.Default.GOPATH+log.SetLogger.Path+"/"+log.SetLogger.FileName+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, PermFileMode)

	logger.SetOutput(file)
	logger.SetFlags(logger.Flags() &^ (logger.Ldate | logger.Ltime))
	logger.SetFlags(0)
	logger.Println(logging)
}

func (log *PatternLogger) GetElapsedTime(start time.Time) int64 {
	var elapsedTime int64 = 0

	if !start.IsZero() {
		elapsedTime = time.Since(start).Milliseconds()
	}

	return elapsedTime
}

func (log *PatternLogger) GetCallersFrames(skipNoOfStack int) *runtime.Frames {
	stackBuf := make([]uintptr, MaxStackLength)
	length := runtime.Callers(skipNoOfStack, stackBuf[:])
	stack := stackBuf[:length]

	return runtime.CallersFrames(stack)
}

func (log *PatternLogger) GetStackTraceError(err error) string {
	frames := log.GetCallersFrames(5)
	trace := err.Error()

	for {
		frame, more := frames.Next()

		if !strings.Contains(frame.File, "runtime/") {
			trace = trace + fmt.Sprintf("\n File: %s, Line: %d, Func: %s", frame.File, frame.Line, frame.Function)
		}
		if !more {
			break
		}
	}
	return trace
}

func (log *PatternLogger) CheckArguments(args []interface{}) (message string, stackTrace string) {
	for i := 0; i < len(args); i++ {
		switch v := args[i].(type) {
		case float32, float64, complex64, complex128:
			message = message + fmt.Sprintf("%s %g", " ", v)
		case int, int8, int32, int64, uint, uint8, uint16, uint32, uint64:
			message = message + fmt.Sprintf("%s %d", " ", v)
		case bool:
			message = message + " " + strconv.FormatBool(v)
		case string:
			message = message + " " + v
		case error:
			stackTrace = stackTrace + " " + log.GetStackTraceError(v)
		case *strconv.NumError:
			stackTrace = stackTrace + " " + log.GetStackTraceError(v)
		default:
			str, _ := json.Marshal(v)
			message = message + " " + string(str)
		}
	}
	return
}

func (log *PatternLogger) Info(correlationID string, message string, args ...interface{}) {
	if len(args) > 0 {
		msg, stackTrace := log.CheckArguments(args)
		log.log(correlationID, INFO, message+msg, 0, stackTrace)
	} else {
		log.log(correlationID, INFO, message, 0)
	}
}

func (log *PatternLogger) Fatal(correlationID string, message string, args ...interface{}) {
	if len(args) > 0 {
		msg, stackTrace := log.CheckArguments(args)
		log.log(correlationID, FATAL, message+msg, 0, stackTrace)
	} else {
		log.log(correlationID, FATAL, message, 0)
	}
}

func (log *PatternLogger) Error(correlationID string, message string, args ...interface{}) {
	if len(args) > 0 {
		msg, stackTrace := log.CheckArguments(args)
		log.log(correlationID, ERROR, message+msg, 0, stackTrace)
	} else {
		log.log(correlationID, ERROR, message, 0)
	}
}

func (log *PatternLogger) Warn(correlationID string, message string, args ...interface{}) {
	if len(args) > 0 {
		msg, stackTrace := log.CheckArguments(args)
		log.log(correlationID, WARN, message+msg, 0, stackTrace)
	} else {
		log.log(correlationID, WARN, message, 0)
	}
}

func (log *PatternLogger) Debug(correlationID string, message string, args ...interface{}) {
	if len(args) > 0 {
		msg, stackTrace := log.CheckArguments(args)
		log.log(correlationID, DEBUG, message+msg, 0, stackTrace)
	} else {
		log.log(correlationID, DEBUG, message, 0)
	}
}

func (log *PatternLogger) Trace(correlationID string, message string, args ...interface{}) {
	if len(args) > 0 {
		msg, stackTrace := log.CheckArguments(args)
		log.log(correlationID, TRACE, message+msg, 0, stackTrace)
	} else {
		log.log(correlationID, TRACE, message, 0)
	}
}
