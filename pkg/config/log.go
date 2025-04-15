package config

import (
	"fmt"

	"github.com/qingchuwudi/syncdns/pkg/tools"
)

// LogConfiguration 日志配置
type LogConfiguration struct {
	Path        string `json:"path,omitempty" yaml:"path,omitempty"`               // 日志文件路径
	FileName    string `json:"filename,omitempty" yaml:"filename,omitempty"`       // 日志文件名
	Level       string `json:"level,omitempty" yaml:"level,omitempty"`             // 记录的日志等级：debug,info,warn,error
	Develop     bool   `json:"develop,omitempty" yaml:"develop,omitempty"`         // 开发者模式，开启后会输出代码文件和堆栈信息
	Size        int    `json:"size,omitempty" yaml:"size,omitempty"`               // 每个日志文件保存的大小 单位:M
	Age         int    `json:"age,omitempty" yaml:"age,omitempty"`                 // 文件最多保存多少天
	BackupCount int    `json:"backupCount,omitempty" yaml:"backupCount,omitempty"` // 日志文件最多保存多少个备份
	LocalTime   bool   `json:"localTime,omitempty" yaml:"localTime,omitempty"`     // 使用本地时间记录
	Compress    bool   `json:"compress,omitempty" yaml:"compress,omitempty"`       // 是否压缩
}

// Validate 检查日志配置有效性
func (l *LogConfiguration) Validate() error {
	if (l.Path != "") && (!tools.IsDir(l.Path)) {
		return fmt.Errorf("日志路径(%s)配置有误：路径不存在或没有权限！", l.Path)
	}
	switch l.Level {
	case "debug", "info", "warn", "error":
		break
	default:
		return fmt.Errorf("日志等级(%s)配置有误", l.Level)
	}
	l.withDefault()
	return nil
}

// 检查配置并修正
func (l *LogConfiguration) withDefault() {
	if l.FileName == "" {
		l.FileName = "github.com/qingchuwudi/syncdns.log"
	}
	if l.Level == "" {
		l.Level = "info"
	}
	if l.Size < 0 {
		l.Size = 16
	}
	if l.Age < 0 {
		l.Age = 0
	}
	if l.BackupCount < 0 {
		l.Age = 0
	}
}
