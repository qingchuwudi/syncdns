/*
 * Copyright (c) 2021 qingchuwudi
 *
 * Licensed under the Apache License, Tag 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * author bypf2009@vip.qq.com
 * create at 2021/12/10
 */

package loger

import (
	"log/slog"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/qingchuwudi/syncdns/pkg/config"
)

var sugLoger *zap.SugaredLogger

var writer *lumberjack.Logger

var (
	Debug func(msg string, keysAndValues ...interface{}) // debug日志
	Info  func(msg string, keysAndValues ...interface{}) //
	Warn  func(msg string, keysAndValues ...interface{})
	Error func(msg string, keysAndValues ...interface{})
	Fatal func(msg string, keysAndValues ...interface{})
	Panic func(msg string, keysAndValues ...interface{})

	Debugf func(template string, args ...interface{}) // debug日志
	Infof  func(template string, args ...interface{}) //
	Warnf  func(template string, args ...interface{})
	Errorf func(template string, args ...interface{})
	Fatalf func(template string, args ...interface{})
	Panicf func(template string, args ...interface{})
)

// 根据配置初始化日志组件
func InitLogger(cfg *config.LogConfiguration) {
	if cfg == nil {
		panic("loger config is nil")
	}
	// 日志配置
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "logger",
		CallerKey:     "file",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		// EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 增加组件配置
	writer = getWriter(cfg)
	var writes = []zapcore.WriteSyncer{zapcore.AddSync(writer)}
	// 同时在控制台上也输出
	writes = append(writes, zapcore.AddSync(os.Stdout))
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	// atomicLevel.SetLevel(zap.DebugLevel)
	atomicLevel.SetLevel(getLevel(cfg.Level))

	// 配置生效
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		// 日志格式默认是Json格式，转为普通格式的日志
		// zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writes...),
		atomicLevel,
	)

	var zapLoger *zap.Logger
	if cfg.Develop {
		// 开启开发模式，堆栈跟踪(可以看到文件名、代码行数)
		// 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数
		caller := zap.AddCaller()
		// 开启文件及行号
		development := zap.Development()
		// 构造日志
		zapLoger = zap.New(core, caller, development)
	} else {
		zapLoger = zap.New(core)
	}

	zapLoger.Info("服务启动，日志记录器启动成功")

	// 赋值
	sugLoger = zapLoger.Sugar()
	Debug = sugLoger.Debugw // debug日志
	Info = sugLoger.Infow   //
	Warn = sugLoger.Warnw
	Error = sugLoger.Errorw
	Fatal = sugLoger.Fatalw
	Panic = sugLoger.Panicw

	// 赋值
	Debugf = sugLoger.Debugf // debug日志
	Infof = sugLoger.Infof   //
	Warnf = sugLoger.Warnf
	Errorf = sugLoger.Errorf
	Fatalf = sugLoger.Fatalf
	Panicf = sugLoger.Panicf
}

// 日志等级
func getLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

// 获取日志输出文件
func getWriter(cfg *config.LogConfiguration) *lumberjack.Logger {
	fullfile := filepath.Join(cfg.Path, cfg.FileName)
	return &lumberjack.Logger{
		Filename:   fullfile,        // 日志文件路径
		MaxSize:    cfg.Size,        // 每个日志文件保存的大小 单位:M
		MaxAge:     cfg.Age,         // 文件最多保存多少天
		MaxBackups: cfg.BackupCount, // 日志文件最多保存多少个备份
		LocalTime:  cfg.LocalTime,   // 使用本地时间记录
		Compress:   cfg.Compress,    // 是否压缩
	}
}

// NewSlogLoger initializes logger with slog.
func NewSlogLoger() *slog.Logger {
	lvl := slog.LevelInfo
	switch config.GetConfig().Log.Level {
	case "debug":
		lvl = slog.LevelDebug
	case "info":
		lvl = slog.LevelInfo
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	}
	lgrOpts := &slog.HandlerOptions{Level: lvl}

	logger := slog.New(slog.NewJSONHandler(writer, lgrOpts))
	return logger
}
