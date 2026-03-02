package logx

import (
	"go.uber.org/zap"
)

// logger 是一个全局的 zap logger 实例。
// 默认是一个无操作的 logger，以防止在 Init 调用前发生 panic。
var logger = zap.NewNop()
var httpLogger = zap.NewNop()

// Init 使用给定的配置初始化全局 logger。
// 这个函数应该在应用程序启动时被调用一次。
func InitDefault(withHttp bool, conf Config) error {
	l, _, err := New(conf)
	if err != nil {
		return err
	}
	logger = l
	httpConf := conf
	httpConf.CsvLog.Path += "_http"
	httpConf.JsonLog.Path += "_http"
	if withHttp {
		l2, _, err := New(httpConf)
		if err != nil {
			return err
		}
		httpLogger = l2
	}

	return nil
}

// Sync 刷新所有缓冲区的日志条目。
// 建议在 main 函数中使用 defer logx.Sync() 来调用它。
func Sync() error {
	_ = logger.Sync()
	_ = httpLogger.Sync()

	return nil
}
