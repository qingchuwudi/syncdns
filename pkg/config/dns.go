package config

import (
	"fmt"
	"net"
	"net/url"

	"github.com/ameshkov/dnsstamps"
	"github.com/mr-karan/doggo/pkg/models"
)

// ------------------------------------------------------------------
// DNS
// ------------------------------------------------------------------

type DnsConfiguration struct {
	Servers            []string `json:"servers" yaml:"servers"`                                     // dns 服务器列表
	UseIPv4            bool     `json:"use_ipv4,omitempty" yaml:"use_ipv4"`                         // 支持ipv4查询
	UseIPv6            bool     `json:"use_ipv6,omitempty" yaml:"use_ipv6"`                         // 支持ipv6查询
	InsecureSkipVerify bool     `json:"insecure_skip_verify,omitempty" yaml:"insecure_skip_verify"` // 忽略dns服务器的tls证书
	Interval           int      `json:"interval,omitempty" yaml:"interval"`                         // 域名查询周期，上一次查询 interval 秒以后开始下一次
	Timeout            int      `json:"timeout,omitempty" yaml:"timeout"`                           // 查询超时时间，默认5秒
	nameservers        []models.Nameserver
}

func (c *DnsConfiguration) Validate() error {
	if c == nil {
		return fmt.Errorf("缺少dns配置")
	}
	if (c.Servers == nil) || (len(c.Servers) < 1) {
		return fmt.Errorf("缺少dns服务器")
	}
	if (!c.UseIPv4) && (!c.UseIPv6) {
		c.UseIPv4 = true
	}
	if c.Interval == 0 {
		c.Interval = 600 // 默认600秒
	}
	if c.Timeout == 0 {
		c.Timeout = 5
	}
	return c.initNameservers()
}

func (c *DnsConfiguration) GetNameservers() []models.Nameserver {
	return c.nameservers
}

func (c *DnsConfiguration) initNameservers() error {
	c.nameservers = make([]models.Nameserver, len(c.Servers))

	for i, server := range c.Servers {
		nsv, err := formatNameserver(server)
		if err != nil {
			return fmt.Errorf("域名服务器 '%s' 加载失败: %v", server, err)
		}
		c.nameservers[i] = nsv
	}
	return nil
}

func formatNameserver(namesever string) (models.Nameserver, error) {
	// Instantiate a UDP resolver with default port as a fallback.
	ns := models.Nameserver{
		Type:    models.UDPResolver,
		Address: net.JoinHostPort(namesever, models.DefaultUDPPort),
	}
	uri, err := url.Parse(namesever)
	if err != nil {
		ip := net.ParseIP(namesever)
		if ip == nil {
			return ns, err
		}
		return ns, nil
	}
	switch uri.Scheme {
	case "sdns":
		stamp, err := dnsstamps.NewServerStampFromString(namesever)
		if err != nil {
			return ns, err
		}
		switch stamp.Proto {
		case dnsstamps.StampProtoTypeDoH:
			ns.Type = models.DOHResolver
			address := url.URL{Scheme: "https", Host: stamp.ProviderName, Path: stamp.Path}
			ns.Address = address.String()
		case dnsstamps.StampProtoTypeDNSCrypt:
			ns.Type = models.DNSCryptResolver
			ns.Address = namesever
		default:
			return ns, fmt.Errorf("unsupported protocol: %v", stamp.Proto.String())
		}

	case "https":
		ns.Type = models.DOHResolver
		ns.Address = uri.String()

	case "tls":
		ns.Type = models.DOTResolver
		if uri.Port() == "" {
			ns.Address = net.JoinHostPort(uri.Hostname(), models.DefaultTLSPort)
		} else {
			ns.Address = net.JoinHostPort(uri.Hostname(), uri.Port())
		}

	case "tcp":
		ns.Type = models.TCPResolver
		if uri.Port() == "" {
			ns.Address = net.JoinHostPort(uri.Hostname(), models.DefaultTCPPort)
		} else {
			ns.Address = net.JoinHostPort(uri.Hostname(), uri.Port())
		}

	case "udp":
		ns.Type = models.UDPResolver
		if uri.Port() == "" {
			ns.Address = net.JoinHostPort(uri.Hostname(), models.DefaultUDPPort)
		} else {
			ns.Address = net.JoinHostPort(uri.Hostname(), uri.Port())
		}
	case "quic":
		ns.Type = models.DOQResolver
		if uri.Port() == "" {
			ns.Address = net.JoinHostPort(uri.Hostname(), models.DefaultDOQPort)
		} else {
			ns.Address = net.JoinHostPort(uri.Hostname(), uri.Port())
		}
	}
	return ns, nil
}
