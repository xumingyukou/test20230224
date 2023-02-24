package logger

import (
	"io/ioutil"
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"

	"github.com/sirupsen/logrus"
)

var (
	Logger *logrus.Entry
)

func init() {
	Logger = NewLogger("logs", "clients.log", time.Hour).WithFields(logrus.Fields{})
}

func LoggerInit(logPath, logName string, logLevel logrus.Level) {
	if logPath == "" {
		logPath = "logs"
	}
	if logName == "" {
		logName = "clients.log"
	}

	tmpLogger := NewLogger(logPath, logName, time.Hour)
	tmpLogger.SetLevel(logLevel)
	Logger = tmpLogger.WithFields(logrus.Fields{})
}

type FileSplitHook struct {
}

func (hook *FileSplitHook) Fire(entry *logrus.Entry) error {
	entry.Data["appName"] = "MyAppName"
	return nil
}

func (hook *FileSplitHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func newLfsHook(logPath string, logName string, logRotationTime time.Duration) logrus.Hook {
	logFile := path.Join(logPath, logName)
	writer, err := rotatelogs.New(
		logFile+".%Y%m%d%H",
		rotatelogs.WithLinkName(logFile),             // WithLinkName 为最新的日志建立软连接，以方便随时找到当前日志文件
		rotatelogs.WithRotationTime(logRotationTime), // WithRotationTime 设置日志分割的时间，这里设置为一小时分割一次
		// WithMaxAge和WithRotationCount 二者只能设置一个，
		//rotatelogs.WithMaxAge(time.Second*24),// WithMaxAge 设置文件清理前的最长保存时间，
		//rotatelogs.WithRotationCount(maxRemainCnt),// WithRotationCount 设置文件清理前最多保存的个数。
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

func NewLogger(logPath string, logName string, logRotationTime time.Duration) *logrus.Logger {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{}
	log.Level = logrus.InfoLevel
	//Logger.SetReportCaller(true)
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		err = os.Mkdir(logPath, 0777)
		if err != nil {
			panic("create log file path error")
		}
	}
	log.SetOutput(ioutil.Discard)
	log.AddHook(newLfsHook(logPath, logName, logRotationTime))

	return log
}
