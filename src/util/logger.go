package util

import (
	"encoding/json"
	"fmt"
	"github.com/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

// error logger
var errorLogger *zap.SugaredLogger

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func getLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

func initLog() {
	fileName := "log/log.log"
	_level := ConfigGetString("log.level")
	if _level == "" {
		_level = "error"
	}
	level := getLoggerLevel(_level)
	syncWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    100, //1G  二进制左移30位
		MaxBackups: 30,  // 最多保留300个备份
		LocalTime:  true,
		Compress:   false, //达到文件大小后需要压缩成gz文件
	})
	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoder), syncWriter, zap.NewAtomicLevelAt(level))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	errorLogger = logger.Sugar()
}

/*
输出调试信息
 */
func LogDebug(args ...interface{}) {
	if ConfigGetString("log.logtype") == "console" {
		if levelMap[ConfigGetString("log.level")] > -1 {
			return
		}
		fmt.Println(time.Now().Format("2006-01-02 15:04:05.000"), args)
		return
	}

	if errorLogger == nil {
		initLog()
	}
	errorLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	if errorLogger == nil {
		initLog()
	}
	errorLogger.Debugf(template, args...)
}

func LogWrite(filename string, str_content ...interface{}) {
	path := IO_GetRootPath()
	content := map[string]interface{}{"content": str_content}
	tmp, _ := json.Marshal(content)
	tmp2 := "\r\n时间：" + time.Now().Format("2006-01-02 15:04:05") + "\r\n" + string(tmp) + "\r\n"
	tmpFileName := "log"
	if len(filename) > 0 {
		tmpFileName = filename
	}
	fd, _ := os.OpenFile(path+"/log/"+time.Now().Format("2006_01_02-")+tmpFileName+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	fd.Write([]byte(tmp2))
	fd.Close()
}

func LogInfo(args ...interface{}) {
	if ConfigGetString("log.logtype") == "console" {
		if levelMap[ConfigGetString("log.level")] > 0 {
			return
		}
		fmt.Println(time.Now().Format("2006-01-02 15:04:05.000"), args)
		return
	}

	if errorLogger == nil {
		initLog()
	}
	errorLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	if errorLogger == nil {
		initLog()
	}
	errorLogger.Infof(template, args...)
}

func LogWarn(args ...interface{}) {
	if ConfigGetString("log.logtype") == "console" {
		if levelMap[ConfigGetString("log.level")] > 1 {
			return
		}
		fmt.Println(time.Now().Format("2006-01-02 15:04:05.000"), args)
		return
	}

	if errorLogger == nil {
		initLog()
	}
	errorLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	if errorLogger == nil {
		initLog()
	}
	errorLogger.Warnf(template, args...)
}

func LogError(args ...interface{}) {
	if ConfigGetString("log.logtype") == "console" {
		if levelMap[ConfigGetString("log.level")] > 2 {
			return
		}
		fmt.Println(time.Now().Format("2006-01-02 15:04:05.000"), args)
		return
	}

	if errorLogger == nil {
		initLog()
	}
	errorLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	if errorLogger == nil {
		initLog()
	}
	errorLogger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	if errorLogger == nil {
		initLog()
	}
	errorLogger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	if errorLogger == nil {
		initLog()
	}
	errorLogger.DPanicf(template, args...)
}

func LogPanic(args ...interface{}) {
	if ConfigGetString("log.logtype") == "console" {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05.000"), args)
		return
	}

	if errorLogger == nil {
		initLog()
	}
	errorLogger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	if errorLogger == nil {
		initLog()
	}
	errorLogger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	if ConfigGetString("log.logtype") == "console" {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05.000"), args)
		return
	}

	if errorLogger == nil {
		initLog()
	}
	errorLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	if errorLogger == nil {
		initLog()
	}
	errorLogger.Fatalf(template, args...)
}
