package zlog

import (
	"fmt"

	"github.com/cihub/seelog"
)

// Init 初始化日志库
func Init(fileName string, loglevel string) {
	seelog.Info("Logdir:", fileName, " Loglevel:", loglevel)

	switch loglevel {
	case "prod":
		loglevel = "info"
	case "debug", "info", "warn", "error":
	default:
		seelog.Error("LogLevel unsupport! [debug/prod/info/warn/error]")
		return
	}

	path := `<seelog minlevel="` + loglevel + `" maxlevel="error">
				<outputs formatid="main">
					<filter levels="debug,error">
						<console />   
					</filter>
					<filter levels="info,warn,error"> 
						<buffered size="10000" flushperiod="1000">
							<rollingfile type="date" filename="` + fileName + `.infolog" datepattern="2006.01.02" fullname="true" maxrolls="168"/>  
						</buffered>
					</filter>
				</outputs>
				<formats>
					<format id="main" format="%Date(2006-01-02 15:04:05.999) [%LEV] [%File:%Line] %Msg%n"/>  
				</formats>
			</seelog>`

	defer seelog.Flush()
	logger, err := seelog.LoggerFromConfigAsString(path)
	if err != nil {
		panic(err)
	}
	seelog.ReplaceLogger(logger)

	Load(fileName)
}

func Load(fileName string) {
	path := `<seelog minlevel="debug" maxlevel="error">
				<outputs formatid="main">
					<filter levels="debug,error">
						<console />   
					</filter>
					<filter levels="info,warn,error,debug"> 
						<buffered size="10000" flushperiod="1000">
							<rollingfile type="date" filename="` + fileName + `.infolog" datepattern="2006.01.02" fullname="true" maxrolls="168"/>  
						</buffered>
					</filter>
				</outputs>
				<formats>
					<format id="main" format="%Date(2006-01-02 15:04:05.999) [%LEV] [%File:%Line] %Msg%n"/>  
				</formats>
			</seelog>`

	var e error
	ilogger, e = seelog.LoggerFromConfigAsString(path)
	ilogger.SetAdditionalStackDepth(1)
	if e != nil {
		seelog.Error(e)
	}
}

var ilogger seelog.LoggerInterface

type Logger struct {
	fmt.Stringer
}

func NewLogger(i fmt.Stringer) *Logger {
	return &Logger{i}
}

// Info 日志
func (l *Logger) Info(v ...interface{}) {
	params := []interface{}{l.String()}
	params = append(params, v...)
	ilogger.Info(params...)
}

// Infof 日志
func (l *Logger) Infof(format string, v ...interface{}) {
	ff := "%s " + format
	params := []interface{}{l.String()}
	params = append(params, v...)
	ilogger.Infof(ff, params...)
}

// Warn 日志
func (l *Logger) Warn(v ...interface{}) {
	params := []interface{}{l.String()}
	params = append(params, v...)
	ilogger.Warn(params...)
}

// Warnf 日志
func (l *Logger) Warnf(format string, v ...interface{}) {
	ff := "%s " + format
	params := []interface{}{l.String()}
	params = append(params, v...)
	ilogger.Warnf(ff, params...)
}

// Error 日志
func (l *Logger) Error(v ...interface{}) {
	params := []interface{}{l.String()}
	params = append(params, v...)
	ilogger.Error(params...)
}

// Errorf 日志
func (l *Logger) Errorf(format string, v ...interface{}) {
	ff := "%s " + format
	params := []interface{}{l.String()}
	params = append(params, v...)
	ilogger.Errorf(ff, params...)
}

// Debug 日志
func (l *Logger) Debug(v ...interface{}) {
	params := []interface{}{l.String()}
	params = append(params, v...)
	ilogger.Debug(params...)
}

// Debugf 日志
func (l *Logger) Debugf(format string, v ...interface{}) {
	ff := "%s " + format
	params := []interface{}{l.String()}
	params = append(params, v...)
	ilogger.Debugf(ff, params...)
}
