package controller

import (
	"go.uber.org/zap"

	"github.com/qingchuwudi/syncdns/pkg/adguardhome"
	"github.com/qingchuwudi/syncdns/pkg/config"
	"github.com/qingchuwudi/syncdns/pkg/loger"
)

// 缓存的数据
type CacheData struct {
	// 域名->解析记录
	DnsRecords []*adguardhome.HostRecord
}

func NewCacheData() (*CacheData, error) {
	data := CacheData{
		DnsRecords: make([]*adguardhome.HostRecord, 0),
	}
	if err := data.LoadData(); err != nil {
		return nil, err
	}
	return &data, nil
}

// LoadData 加载待同步的域名列表，优先以配置文件为准
func (c *CacheData) LoadData() error {
	adhRecords, err := adguardhome.Client().HostRecordList()
	if err != nil {
		loger.Error("AdguardHome DNS重写数据获取失败", zap.Error(err))
		return err
	}
	// 配置文件中没有配置domain，以adguardhome为准
	if len(config.GetConfig().Domain) == 0 {
		c.DnsRecords = adhRecords
		return nil
	}

	// 否则以配置文件为准
	c.DnsRecords = make([]*adguardhome.HostRecord, 0, len(config.GetConfig().Domain))
	for _, domain := range config.GetConfig().Domain {
		has := false
		for index := range adhRecords {
			if adhRecords[index].Host == domain {
				has = true
				c.DnsRecords = append(c.DnsRecords, adhRecords[index])
				break
			}
		}
		if !has {
			c.DnsRecords = append(c.DnsRecords, &adguardhome.HostRecord{Host: domain})
		}
	}
	return nil
}
