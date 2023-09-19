package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"chainmaker.org/chainmaker/common/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// nolint
const (
	ModuleWeb     = "[WEB]"
	ModuleDb      = "[DB]"
	ModuleSession = "[SESSION]"
	K8s           = "[KUBERNETES]"
)

// LogConf 日志配置
type LogConf struct {
	LogLevelDefault string            `mapstructure:"log_level_default"`
	LogLevels       map[string]string `mapstructure:"log_levels"`
	FilePath        string            `mapstructure:"file_path"`
	MaxAge          int               `mapstructure:"max_age"`
	RotationTime    int               `mapstructure:"rotation_time"`
	LogInConsole    bool              `mapstructure:"log_in_console"`
	ShowColor       bool              `mapstructure:"show_color"`
}

var (
	loggers = make(map[string]*zap.SugaredLogger)
	// map[module-name]map[module-name+chainId]zap.AtomicLevel
	loggerLevels = make(map[string]map[string]zap.AtomicLevel)
	loggerMutex  sync.Mutex
	logConfig    *LogConf
)

// SetLogConfig - 设置Log配置对象
func SetLogConfig(config *LogConf) {
	logConfig = config
}

// GetLogger - 获取Logger对象
func GetLogger(name string) *zap.SugaredLogger {
	return GetLoggerByChain(name, "")
}

// GetLoggerByChain - 获取带链标识的Logger对象
func GetLoggerByChain(name, chainId string) *zap.SugaredLogger {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	var conf log.LogConfig
	var pureName string
	logHeader := name + chainId
	logger, ok := loggers[logHeader]
	if !ok {
		if logConfig == nil {
			logConfig = DefaultLogConfig()
		}
		if logConfig.LogLevelDefault == "" {
			defaultLogNode := GetDefaultLogNodeConfig()
			conf = log.LogConfig{
				Module:       "[DEFAULT]",
				ChainId:      chainId,
				LogPath:      defaultLogNode.FilePath,
				LogLevel:     log.GetLogLevel(defaultLogNode.LogLevelDefault),
				MaxAge:       defaultLogNode.MaxAge,
				RotationTime: defaultLogNode.RotationTime,
				JsonFormat:   false,
				ShowLine:     true,
				LogInConsole: defaultLogNode.LogInConsole,
				ShowColor:    defaultLogNode.ShowColor,
			}
		} else {
			pureName = strings.ToLower(strings.Trim(name, "[]"))
			value, exists := logConfig.LogLevels[pureName]
			if !exists {
				value = logConfig.LogLevelDefault
			}
			conf = log.LogConfig{
				Module:       name,
				ChainId:      chainId,
				LogPath:      logConfig.FilePath,
				LogLevel:     log.GetLogLevel(value),
				MaxAge:       logConfig.MaxAge,
				RotationTime: logConfig.RotationTime,
				JsonFormat:   false,
				ShowLine:     true,
				LogInConsole: logConfig.LogInConsole,
				ShowColor:    logConfig.ShowColor,
			}
		}
		var level zap.AtomicLevel
		logger, level = log.InitSugarLogger(&conf)
		loggers[logHeader] = logger
		if pureName != "" {
			if _, exist := loggerLevels[pureName]; !exist {
				loggerLevels[pureName] = make(map[string]zap.AtomicLevel)
			}
			loggerLevels[pureName][logHeader] = level
		}
	}
	return logger
}

// RefreshLogConfig refresh log config
func RefreshLogConfig(config *LogConf) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	// scan loggerLevels and find the level from config, if can't find level, set it to default
	for name, loggers := range loggerLevels {
		var (
			logLevel zapcore.Level
			strLevel string
			exist    bool
		)
		if strLevel, exist = config.LogLevels[name]; !exist {
			strLevel = config.LogLevelDefault
		}
		switch log.GetLogLevel(strLevel) {
		case log.LEVEL_DEBUG:
			logLevel = zap.DebugLevel
		case log.LEVEL_INFO:
			logLevel = zap.InfoLevel
		case log.LEVEL_WARN:
			logLevel = zap.WarnLevel
		case log.LEVEL_ERROR:
			logLevel = zap.ErrorLevel
		default:
			logLevel = zap.InfoLevel
		}
		for _, aLevel := range loggers {
			aLevel.SetLevel(logLevel)
		}
	}
}

// DefaultLogConfig - 获取默认Log配置
func DefaultLogConfig() *LogConf {
	defaultLogNode := GetDefaultLogNodeConfig()
	return &LogConf{
		LogLevelDefault: defaultLogNode.LogLevelDefault,
		FilePath:        defaultLogNode.FilePath,
		MaxAge:          defaultLogNode.MaxAge,
		RotationTime:    defaultLogNode.RotationTime,
		LogInConsole:    defaultLogNode.LogInConsole,
	}
}

// GetDefaultLogNodeConfig get default log node config
func GetDefaultLogNodeConfig() LogConf {
	return LogConf{
		LogLevelDefault: log.DEBUG,
		FilePath:        "../log/web.log",
		MaxAge:          log.DEFAULT_MAX_AGE,
		RotationTime:    log.DEFAULT_ROTATION_TIME,
		LogInConsole:    true,
		ShowColor:       true,
	}
}

// GetCurrentPath get current path
func GetCurrentPath() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)
}
