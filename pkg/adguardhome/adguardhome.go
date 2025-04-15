package adguardhome

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/qingchuwudi/syncdns/pkg/config"
	"github.com/qingchuwudi/syncdns/pkg/loger"
)

// 参考 AdGuardHome/internal/filtering/rewritehttp.go
type DnsRewriteEntry struct {
	Domain string `json:"domain"`
	Answer string `json:"answer"`
}

const (
	timeout             = 3 * time.Second
	DnsRewriteListUri   = "/control/rewrite/list"
	DnsRewriteAddUri    = "/control/rewrite/add"
	DnsRewriteDeleteUri = "/control/rewrite/delete"
	DnsRewriteUpdateUri = "/control/rewrite/update"
)

type Adh struct {
	ctx    context.Context
	client *resty.Client
}

// ------------------------------------------------------------------
//
// ------------------------------------------------------------------

var client *Adh

func Client() *Adh {
	return client
}

func InitClient(ctx context.Context) {
	if client != nil {
		return
	}
	client = NewClient(ctx)
}

// NewClient 创建新客户端
func NewClient(ctx context.Context) *Adh {
	adhCfg := config.GetConfig().Adh
	rc := resty.New()
	if adhCfg.Timeout > 0 {
		rc = rc.SetTimeout(time.Duration(adhCfg.Timeout) * time.Second)
	} else {
		rc = rc.SetTimeout(timeout)
	}
	// header
	rc.SetBasicAuth(adhCfg.Username, adhCfg.Password)
	rc.SetHeader("Referer", adhCfg.Host)
	rc.SetBaseURL(adhCfg.Host)
	//
	return &Adh{ctx: ctx, client: rc}
}

// HostRecordList 获取所有adgardhome已有数据，并转换成缓存格式的数据
func (a *Adh) HostRecordList() (records []*HostRecord, err error) {
	entries, err := a.List()
	if err != nil {
		return nil, err
	}
	if (entries == nil) || (len(entries) < 1) {
		return nil, nil
	}
	recordMap := make(map[string]*HostRecord)
	for _, entry := range entries {
		rec := recordMap[entry.Domain]
		if rec == nil {
			rec = &HostRecord{
				Host: entry.Domain,
				IPs:  make([]string, 0),
			}
			recordMap[entry.Domain] = rec
		}
		rec.IPs = append(rec.IPs, entry.Answer)
	}
	records = make([]*HostRecord, 0, len(recordMap))
	for _, record := range recordMap {
		record.RemoveDuplicate()
		records = append(records, record.DeepCopy())
	}
	return records, nil
}

// List 查看列表
func (a *Adh) List() ([]DnsRewriteEntry, error) {
	resp, err := a.newReq().Get(DnsRewriteListUri)
	if err != nil {
		loger.Error("AdGuardHome获取DNS重写数据失败", zap.String("path", DnsRewriteListUri), zap.Error(err))
		return nil, err
	}
	data := make([]DnsRewriteEntry, 0)
	if err = json.Unmarshal(resp.Body(), &data); err != nil {
		loger.Error("AdGuardHome DNS重写数据解析失败s", zap.String("path", DnsRewriteListUri), zap.Error(err))
		return nil, err
	}
	return data, nil
}

// Add 添加数据
func (a *Adh) Add(domain, answer string) error {
	resp, err := a.newReq().
		SetBody(DnsRewriteEntry{Domain: domain, Answer: answer}).
		Post(DnsRewriteAddUri)
	if err != nil {
		loger.Error("AdGuardHome添加记录失败", zap.String(domain, answer), zap.Error(err))
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("response http code : %d", resp.StatusCode())
	}
	loger.Debug("AdGuardHome添加记录成功", zap.String(domain, answer))
	return nil
}

// Delete 删除一条数据
func (a *Adh) Delete(domain, answer string) error {
	resp, err := a.newReq().
		SetBody(DnsRewriteEntry{Domain: domain, Answer: answer}).
		Post(DnsRewriteDeleteUri)
	if err != nil {
		loger.Error("AdGuardHome删除失败", zap.String(domain, answer), zap.Error(err))
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("response http code : %d", resp.StatusCode())
	}
	return nil
}

// Update 更新数据
func (a *Adh) Update(domain, answer string) error {
	resp, err := a.newReq().
		SetBody(DnsRewriteEntry{Domain: domain, Answer: answer}).
		Put(DnsRewriteUpdateUri)
	if err != nil {
		loger.Error("AdGuardHome更新失败", zap.String(domain, answer), zap.Error(err))
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("response http code : %d", resp.StatusCode())
	}
	return nil
}

func (a *Adh) newReq() *resty.Request {
	return a.client.R().SetContext(a.ctx).SetHeader("Content-Type", "application/json")
}
