/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/24 12:54
 */

package core

import (
	"github.com/sirupsen/logrus"
	"github.com/warlock-backend/core/internal"
	"os"
)

func Logrus() (logger *logrus.Logger) {
	logger = logrus.New()                         //新建一个实例
	logger.SetOutput(os.Stderr)                   //设置输出类型
	logger.SetReportCaller(true)                  //开启返回函数名和行号
	logger.SetFormatter(&internal.LogFormatter{}) //设置自己定义的Formatter
	logger.SetLevel(logrus.DebugLevel)            //设置最低的Level

	return logger
}
