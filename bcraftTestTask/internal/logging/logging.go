package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
)

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, writer := range hook.Writer {
		_, err := writer.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

var logger *logrus.Logger

func GetLogger() *logrus.Logger {
	return logger
}

func init() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		FullTimestamp: true,
	}

	err := os.MkdirAll("./logs", 0755)
	if err != nil {
		panic(err)
	}

	log_file, err := os.OpenFile("./logs/logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		panic(err)
	}

	err_log_file, err := os.OpenFile("./logs/err_logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		panic(err)
	}

	l.SetOutput(io.Discard)

	l.AddHook(&writerHook{
		Writer:    []io.Writer{log_file},
		LogLevels: logrus.AllLevels,
	})

	l.AddHook(&writerHook{
		Writer:    []io.Writer{os.Stdout, err_log_file},
		LogLevels: []logrus.Level{logrus.FatalLevel, logrus.ErrorLevel, logrus.PanicLevel},
	})

	l.SetLevel(logrus.InfoLevel)
	logger = l
}
