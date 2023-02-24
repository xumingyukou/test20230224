package sdk

// 使用uber/zap+lumerjack的高性能日志组合
import (
	"os"
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// zap要动态修改level需要通过返回的AtomicLevel设置
func NewLoggerZap(logPath string, logName string, logLevel string) (*zap.SugaredLogger, zap.AtomicLevel) {
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		err = os.Mkdir(logPath, 0777)
		if err != nil {
			panic("create log file path error")
		}
	}
	atom := zap.NewAtomicLevel()
	logFile := path.Join(logPath, logName)

	writeSyncer := getLogWriter(logFile)
	encoder := getEncoder()
	l, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		panic(err)
	}
	atom.SetLevel(l)
	core := zapcore.NewCore(encoder, writeSyncer, atom)

	logger := zap.New(core)
	return logger.Sugar(), atom
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = nil
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(logPath string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100, // megabytes
		MaxBackups: 20,
		MaxAge:     30, //days
		LocalTime:  true,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func SetLevelZap(al *zap.AtomicLevel, logLevel string) error {
	l, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		return err
	}
	al.SetLevel(l)
	return nil
}
