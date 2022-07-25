/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/24 12:56
 */

package internal

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

//颜色
const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

type LogFormatter struct{}

// Format 实现Formatter(entry *logrus.Entry) ([]byte, error)接口
func (t *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	//根据不同的level去展示颜色
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	//自定义日期格式
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var newLog string
	//HasCaller()为true才会有调用信息
	if entry.HasCaller() {
		//自定义文件路径格式
		fileVal := filepath.Base(entry.Caller.File)
		line := entry.Caller.Line
		//自定义输出格式
		newLog = fmt.Sprintf("[%s] \x1b[%dm[%s]\x1b[0m %s:%d %s :%s\n", timestamp, levelColor,
			entry.Level, fileVal, line, entry.Caller.Function, entry.Message)

	} else {
		newLog = fmt.Sprintf("[%s] \x1b[%dm[%s]\x1b[0m  :%s\n", timestamp, levelColor, entry.Level, entry.Message)
	}
	b.WriteString(newLog)
	return b.Bytes(), nil
}
