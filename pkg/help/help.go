package help

import (
	"flag"
	"fmt"
	"os"

	"github.com/qingchuwudi/syncdns/pkg/loger"
)

var (
	Help bool
	Cfg  string
)

func init() {
	flag.BoolVar(&Help, "h", false, "查看帮助")
	flag.StringVar(&Cfg, "c", "config.yaml", "配置文件。例如: -c config.yaml")
	flag.Usage = Usage
}

func Usage() {
	fmt.Fprintf(os.Stderr, `dns同步工具.
Usage: github.com/qingchuwudi/syncdns [-h | -c ]

Options:
`)
	flag.PrintDefaults()
}

// ParseArgs 获取命令行参数（读取并解析）
func ParseArgs() (stop bool) {
	flag.Parse()
	if Help {
		flag.Usage()
		return true
	}
	if Cfg == "" {
		loger.PreError("请指定配置文件")
		return true
	}

	// 检查文件是否存在
	if !IsFileValid(Cfg) {
		loger.PreError("配置文件不存在，或者没有权限")
		return true
	}
	return false
}

// IsFileValid 判断所给路径文件/文件夹是否存在
func IsFileValid(file string) bool {
	stu, err := os.Stat(file) // os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		loger.PreError("文件 '%s' 不存在或没有权限。", file)
		return false
	}
	if stu.IsDir() {
		loger.PreError("参数错误，'%s' 是文件夹。", file)
		return false
	}
	return true
}
