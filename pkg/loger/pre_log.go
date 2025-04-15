package loger

import (
	"fmt"
)

// ------------------------------------------------------------------
// 透明背景的输出
// ------------------------------------------------------------------

func PreSuccess(format string, a ...interface{}) {
	fmt.Printf("\033[1;32;48m"+format+"\033[0m\n", a...)
}

func PreInfo(format string, a ...interface{}) {
	fmt.Printf("\033[1;37;48m"+format+"\033[0m\n", a...)
}

func PreError(format string, a ...interface{}) {
	fmt.Printf("\033[1;31;48m"+format+"\033[0m\n", a...)
}

// ------------------------------------------------------------------
// 带背景的输出
// ------------------------------------------------------------------

// 背景
func PreSuccessHeav(format string, a ...interface{}) {
	fmt.Printf("\033[1;32;47m"+format+"\033[0m\n", a...)
}

// 背景黄色
func PreInfoHeav(format string, a ...interface{}) {
	fmt.Printf("\033[1;37;43m"+format+"\033[0m\n", a...)
}

// 重量级错误提示，背景黑色
func PreErrorHeav(format string, a ...interface{}) {
	fmt.Printf("\033[1;31;40m"+format+"\033[0m\n", a...)
}
