package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

//定义函数--日志输出格式
func NewStdoutConsole() *zap.SugaredLogger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder  //北京时间
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder //日志级别为大写并带颜色
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), zapcore.DebugLevel) //编译器、输出目标、日志级别
	logger := zap.New(core)
	return logger.Sugar()
}
