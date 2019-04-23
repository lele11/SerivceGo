package logger

import (
	"github.com/cihub/seelog"
)

// Init 初始化日志库
func Init(fileName string, loglevel string) seelog.LoggerInterface {
	seelog.Info("Logdir:", fileName, " Loglevel:", loglevel)

	switch loglevel {
	case "prod":
		loglevel = "info"
	case "debug", "info", "warn", "error":
	default:
		seelog.Error("LogLevel unsupport! [debug/prod/info/warn/error]")
		return nil
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
	return logger
}
