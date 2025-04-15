package controller

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/qingchuwudi/syncdns/pkg/adguardhome"
	"github.com/qingchuwudi/syncdns/pkg/loger"
	"github.com/qingchuwudi/syncdns/pkg/mdns"
	"github.com/qingchuwudi/syncdns/pkg/tools"
)

type ControllerState struct {
	ctx    context.Context
	cancel context.CancelFunc
	data   *CacheData
}

var controller *ControllerState

func NewController(ctx context.Context, cancel context.CancelFunc) *ControllerState {
	controller = &ControllerState{
		ctx:    ctx,
		cancel: cancel,
	}
	return controller
}

func (c *ControllerState) Run() {
	loger.Info("controller starting")
	defer c.cancel()
	var (
		err error
	)
	// 获取域名列表，加载到缓存
	if controller.data, err = NewCacheData(); err != nil {
		loger.Error("Controller create CacheData", zap.Error(err))
		return
	}
	lookupTimer := time.NewTimer(firstTime())
	recacheTimer := time.NewTicker((18 * time.Hour) - time.Minute)
	loger.Info("controller started")
	for {
		select {
		case <-c.ctx.Done():
			// 2. 处理退出事件
			loger.Info("controller stopped: context cancelled")
			return
		case <-recacheTimer.C:
			if data, err := NewCacheData(); err != nil {
				loger.Error("controller re-create CacheData", zap.Error(err))
			} else {
				controller.data = data
			}
		case <-lookupTimer.C:
			// 1. 定时做域名解析查询
			//    1.1 获取解析记录
			//    1.2 所有域名全部解析结束以后，对比异同，有变化，启动一个协程批量处理
			loger.Info("controller sync dns start")
			for index, record := range c.data.DnsRecords {
				// 一个个的查询并同步，分散dns和http请求的时间
				record = mdns.NSLookupOne(c.ctx, record)
				c.data.DnsRecords[index] = record
				syncOneRecordToAdGuardHome(c.ctx, record)
			}
			nextT := tools.NextHour()
			loger.Info("controller done and next sync time", zap.Time("time", time.Now().Add(nextT)))
			lookupTimer.Reset(nextT)
		}
	}
}

// syncOneRecordToAdGuardHome 新数据同步到AdGuardHome
func syncOneRecordToAdGuardHome(ctx context.Context, record *adguardhome.HostRecord) {
	if record == nil {
		return
	}
	del, add := record.IpDiff()
	if len(del) == 0 && len(add) == 0 {
		return
	}
	for _, ip := range del {
		select {
		case <-ctx.Done():
			loger.Warn("AdGuardHome delete exiting: context cancelled")
			return
		default:
			if err := adguardhome.Client().Delete(record.Host, ip); err != nil {
				loger.Warn("AdGuardHome delete fail", zap.String(record.Host, ip), zap.Error(err))
			}
		}
	}
	for _, ip := range add {
		select {
		case <-ctx.Done():
			loger.Warn("AdGuardHome add exiting: context cancelled")
			return
		default:
			if err := adguardhome.Client().Add(record.Host, ip); err != nil {
				loger.Warn("AdGuardHome add fail", zap.String(record.Host, ip), zap.Error(err))
			}
		}
	}
	loger.Info("sync record to AdGuardHome",
		zap.String("domain", record.Host),
		zap.Strings("delete", del), zap.Strings("add", add))
}

// 首次运行时要保证第二次和第一次之间间隔超过10分钟
func firstTime() time.Duration {
	h := tools.NextHour()
	if h < (10 * time.Minute) {
		return h
	}
	return 3 * time.Second
}
