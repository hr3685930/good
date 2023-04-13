package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

// Logger Logger
var Logger *zap.Logger

// Log Log
type Log struct {
	filePath string
	fileName string
}

// NewLog NewLog
func NewLog(filePath, fileName string) *Log {
	return &Log{filePath, fileName}
}

// Init Init
func (z *Log) Init() (err error) {
	if err := z.createFile(z.filePath, z.fileName); err != nil {
		return err
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	atom := zap.NewAtomicLevelAt(zap.DebugLevel)

	config := zap.Config{
		Level:            atom,
		Development:      true,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{z.filePath + z.fileName},
		ErrorOutputPaths: []string{"stderr"},
	}

	Logger, err = config.Build()
	if err != nil {
		return err
	}
	return nil
}

func (z *Log) createFile(filepath, filename string) error {
	_, err := os.Stat(filepath + filename)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(filepath, 0777)
			if err != nil {
				return err
			}
			f, err := os.Create(filepath + filename)
			if err != nil {
				return err
			}
			_ = f.Close()
		}
	}
	return nil
}
