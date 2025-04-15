package mdns

import (
	"context"

	kgdns "github.com/miekg/dns"
	"github.com/mr-karan/doggo/pkg/resolvers"
	"go.uber.org/zap"

	"github.com/qingchuwudi/syncdns/pkg/adguardhome"
	"github.com/qingchuwudi/syncdns/pkg/config"
	"github.com/qingchuwudi/syncdns/pkg/loger"
)

const answerCapMin = 1 << 2

// NSLookupOne 查询一个域名
func NSLookupOne(ctx context.Context, record *adguardhome.HostRecord) *adguardhome.HostRecord {
	loger.Debug("NSLookupOne record", zap.String("domain", record.Host))
	rslvr := GetResolver()
	if answers := doLookup(ctx, rslvr, record.Host); len(answers) > 0 {
		record.IPFromAnswer(answers)
		loger.Debug("NSLookupOne success", zap.Strings(record.Host, record.IPs))
	}
	return record
}

// doLookup 对某一个域名进行递归解析
func doLookup(ctx context.Context, rslvr *RslvrEntry, domain string) (answers []resolvers.Answer) {
	questions := []kgdns.Question{}
	if config.GetConfig().Dns.UseIPv4 {
		questions = append(questions, kgdns.Question{
			Name:   domain,
			Qtype:  kgdns.TypeA,
			Qclass: kgdns.ClassINET,
		})
	}
	if config.GetConfig().Dns.UseIPv6 {
		questions = append(questions, kgdns.Question{
			Name:   domain,
			Qtype:  kgdns.TypeAAAA,
			Qclass: kgdns.ClassINET,
		})
	}
	flags := resolvers.QueryFlags{AA: true, AD: false, CD: false, RD: true, DO: true}
	resps, err := rslvr.Rslvr.Lookup(ctx, questions, flags)
	if err != nil {
		loger.Error("dns lookup error", zap.String(domain, rslvr.Address), zap.Error(err))
		return nil
	}
	if (len(resps) == 0) || (len(resps[0].Answers) == 0) {
		return nil
	}

	answers = make([]resolvers.Answer, 0, answerCapMin)
	for _, resp := range resps {
		for _, answer := range resp.Answers {
			loger.Debug("answer", zap.Any("answer", answer))
			switch answer.Type {
			case "A":
				answers = append(answers, answer)
			case "AAAA":
				answers = append(answers, answer)
			default:
				loger.Debug("dns lookup unexpected type", zap.String(answer.Name, answer.Type))
			}
		}
	}
	return answers
}
