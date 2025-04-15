package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/qingchuwudi/syncdns/pkg/adguardhome"
	"github.com/qingchuwudi/syncdns/pkg/config"
	"github.com/qingchuwudi/syncdns/pkg/controller"
	"github.com/qingchuwudi/syncdns/pkg/help"
	"github.com/qingchuwudi/syncdns/pkg/loger"
	"github.com/qingchuwudi/syncdns/pkg/mdns"
)

func main() {
	if help.ParseArgs() {
		help.Usage()
		return
	}

	if err := config.LoadFromFile(help.Cfg); err != nil {
		loger.PreError("配置文件加载失败：%v", err)
		return
	}

	loger.InitLogger(config.GetConfig().Log)
	if err := mdns.InitResolver(); err != nil {
		loger.PreError(err.Error())
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	adguardhome.InitClient(ctx)
	controller.NewController(ctx, stop).Run()
	<-ctx.Done()
}
