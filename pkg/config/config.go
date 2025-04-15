package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var loadedConfig Configuration

func LoadFromFile(configFile string) error {
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, &loadedConfig)
	if err != nil {
		return err
	}
	if err = loadedConfig.Validate(); err != nil {
		return err
	}
	return nil
}

func GetConfig() *Configuration {
	return &loadedConfig
}

// Configuration 配置文件
type Configuration struct {
	// AdGuard Home
	Adh *AdGuardHomeConfiguration `json:"adguardhome" yaml:"adguardhome"`
	// dns
	Dns *DnsConfiguration `json:"dns" yaml:"dns"`
	// log
	Log *LogConfiguration `json:"log,omitempty" yaml:"log,omitempty"`
	// domains 通常不需要从配置文件中加载，每次启动在AdGuardHome中拉取。
	// 如果这里配置了域名，则与AdGuardHome重写记录进行合并。
	Domain DomainConfiguration `json:"domain,omitempty" yaml:"domain,omitempty"`
}

func (c *Configuration) Validate() error {
	if err := c.Adh.Validate(); err != nil {
		return err
	}
	if err := c.Dns.Validate(); err != nil {
		return err
	}
	if err := c.Log.Validate(); err != nil {
		return err
	}
	if err := c.Domain.Validate(); err != nil {
		return err
	}
	return nil
}
