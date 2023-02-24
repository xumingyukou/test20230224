package sdk

import (
	"io"
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// time="2022-03-22T19:38:52+08:00" level=info msg="test Info 日志日志日志日志日志日志日志日志日志日志日志日志" func=TickerTape/utils.Test_Logger file="/Users/mingguo/code/go/TickerTape/utils/logger_test.go:9"
type FileSplitHook struct {
}

func (hook *FileSplitHook) Fire(entry *logrus.Entry) error {
	entry.Data["appName"] = "MyAppName"
	return nil
}

func (hook *FileSplitHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func newLfsHook(logPath string, logName string, logRotationTime, logAge time.Duration) logrus.Hook {
	logFile := path.Join(logPath, logName)
	writer, err := rotatelogs.New(
		logFile+".%Y%m%d%H",
		// WithLinkName 为最新的日志建立软连接，以方便随时找到当前日志文件
		rotatelogs.WithLinkName(logFile),
		// WithRotationTime 设置日志分割的时间，这里设置为一小时分割一次
		rotatelogs.WithRotationTime(logRotationTime),
		// WithMaxAge和WithRotationCount 二者只能设置一个，
		// WithMaxAge 设置文件清理前的最长保存时间，
		// WithRotationCount 设置文件清理前最多保存的个数。
		rotatelogs.WithMaxAge(logAge),
		//rotatelogs.WithRotationCount(maxRemainCnt),
	)

	if err != nil {
		logrus.Errorf("config local file system for logger error: %v", err)
	}
	logrus.SetLevel(logrus.InfoLevel)
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.TextFormatter{DisableColors: true})
	return lfsHook
}

func NewLogger(logPath string, logName string, logRotationTime, logAge time.Duration) *logrus.Logger {
	Logger := logrus.New()
	Logger.Formatter = &logrus.TextFormatter{}
	Logger.Level = logrus.InfoLevel
	//Logger.SetReportCaller(true)
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		err = os.Mkdir(logPath, 0777)
		if err != nil {
			panic("create log file path error")
		}
	}
	Logger.SetOutput(io.Discard)
	Logger.AddHook(newLfsHook(logPath, logName, logRotationTime, logAge))

	return Logger
}

func SetLevel(logger *logrus.Logger, l logrus.Level) {
	logger.SetLevel(l)
}
