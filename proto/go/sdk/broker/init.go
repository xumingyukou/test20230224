package broker

import "github.com/sirupsen/logrus"

// 调用库之前应该初始化次Logger，用于输出变量
var BrokerLogger *logrus.Logger

func init() {
	BrokerLogger = logrus.StandardLogger()
}
