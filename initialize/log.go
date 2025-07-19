package initialize

import (
	"gochat/global"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Reset  = "\033[0m"
)

func InitLogger() *zap.Logger {
	writeSyncer := getLogWrite(global.Config.Zap.Filename, global.Config.Zap.MaxSize, global.Config.Zap.MaxBackups, global.Config.Zap.MaxAge)
	if global.Config.Zap.IsConsolePrint {
		writeSyncer = zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(zapcore.Lock(os.Stdout)))
	}
	encoder := getEncoder()

	var loglevel zapcore.Level
	if err := loglevel.UnmarshalText([]byte(global.Config.Zap.Level)); err != nil {
		log.Fatal("无法解析日志等级:% v ", err)
	}
	core := zapcore.NewCore(encoder, writeSyncer, loglevel)
	logger := zap.New(core, zap.AddCaller())
	return logger
}

func coloredLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch level {
	case zapcore.DebugLevel:
		enc.AppendString(Green + level.String() + Reset)
	case zapcore.InfoLevel:
		enc.AppendString(Blue + level.String() + Reset)
	case zapcore.WarnLevel:
		enc.AppendString(Yellow + level.String() + Reset)
	case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		enc.AppendString(Red + level.String() + Reset)
	default:
		enc.AppendString(level.String())
	}
}

func getLogWrite(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = coloredLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
