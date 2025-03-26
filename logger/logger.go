// this is a customized wrapper package on uber's zap logger

package logs

import (
	"context"
	constants "hanamark/constant"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitializeLogger() (*zap.Logger, error) {
	// using the zap's production logging snippet and modifying according to my needs
	config := zap.NewProductionEncoderConfig()

	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	path := viper.GetString("logger.filepath")
	logFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	writer := zapcore.AddSync(logFile)
	//TODO:read from the config write a switch case and set the default log level
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),                        //  writes the file into the file
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel), // stdouterr the error
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger, nil
}

func SetLoggerctx(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, constants.CONTEXT_KEY_LOGGER, l)
}

func GetLoggerctx(ctx context.Context) (l *zap.Logger) {
	val := ctx.Value(constants.CONTEXT_KEY_LOGGER)
	if val == nil {
		return nil
	}
	l, ok := val.(*zap.Logger)
	if !ok {
		return nil
	}
	return l
}
