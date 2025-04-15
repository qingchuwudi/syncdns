package mdns

import (
	"math/rand/v2"
	"sync"
	"time"

	"github.com/mr-karan/doggo/pkg/models"
	"github.com/mr-karan/doggo/pkg/resolvers"

	"github.com/qingchuwudi/syncdns/pkg/config"
	"github.com/qingchuwudi/syncdns/pkg/loger"
)

var (
	resolverLock = sync.RWMutex{}
)

var myResolvers []RslvrEntry

// InitResolver 初始化，加载所有dns解析器
func InitResolver() (err error) {
	resolverLock.Lock()
	myResolvers, err = parseResolvers(config.GetConfig().Dns)
	resolverLock.Unlock()
	return err
}

// GetResolver 随机获取一个dns解析器
func GetResolver() *RslvrEntry {
	resolverLock.RLock()
	defer resolverLock.RUnlock()
	return &myResolvers[rand.N(len(myResolvers))]
}

type RslvrEntry struct {
	models.Nameserver
	Rslvr resolvers.Resolver
}

// parseResolvers loads differently configured
// resolvers based on a list of nameserver.
func parseResolvers(cfg *config.DnsConfiguration) ([]RslvrEntry, error) {
	opts := resolvers.Options{
		UseIPv4:            cfg.UseIPv4,
		UseIPv6:            cfg.UseIPv6,
		Timeout:            time.Duration(cfg.Timeout) * time.Second,
		InsecureSkipVerify: cfg.InsecureSkipVerify,
		Logger:             loger.NewSlogLoger(),
	}
	// For each nameserver, initialise the correct resolver.
	rslvrs := make([]RslvrEntry, 0, len(cfg.GetNameservers()))
	var (
		ok    = true
		rslvr resolvers.Resolver
		err   error
	)
	for _, ns := range cfg.GetNameservers() {
		switch ns.Type {
		case models.DOHResolver:
			loger.Info("initiating DOH resolver")
			rslvr, err = resolvers.NewDOHResolver(ns.Address, opts)
			if err != nil {
				return nil, err
			}
		case models.DOTResolver:
			loger.Info("initiating DOT resolver")
			rslvr, err = resolvers.NewClassicResolver(ns.Address,
				resolvers.ClassicResolverOpts{
					UseTLS: true,
					UseTCP: true,
				}, opts)

			if err != nil {
				return nil, err
			}
		case models.TCPResolver:
			loger.Info("initiating TCP resolver")
			rslvr, err = resolvers.NewClassicResolver(ns.Address,
				resolvers.ClassicResolverOpts{
					UseTLS: false,
					UseTCP: true,
				}, opts)
			if err != nil {
				return nil, err
			}
		case models.UDPResolver:
			loger.Info("initiating UDP resolver")
			rslvr, err = resolvers.NewClassicResolver(ns.Address,
				resolvers.ClassicResolverOpts{
					UseTLS: false,
					UseTCP: false,
				}, opts)
			if err != nil {
				return nil, err
			}
		case models.DNSCryptResolver:
			loger.Info("initiating DNSCrypt resolver")
			rslvr, err = resolvers.NewDNSCryptResolver(ns.Address,
				resolvers.DNSCryptResolverOpts{
					UseTCP: false,
				}, opts)
			if err != nil {
				return nil, err
			}
		case models.DOQResolver:
			loger.Info("initiating DOQ resolver")
			rslvr, err = resolvers.NewDOQResolver(ns.Address, opts)
			if err != nil {
				return nil, err
			}
		default:
			ok = false
		}
		if ok {
			rslvrs = append(rslvrs, RslvrEntry{
				Nameserver: ns,
				Rslvr:      rslvr,
			})
		}
		ok = true
	}
	return rslvrs, nil
}
