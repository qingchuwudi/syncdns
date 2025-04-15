package config

import "fmt"

type AdGuardHomeConfiguration struct {
	Host     string `json:"host" yaml:"host"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Timeout  uint   `json:"timeout,omitempty" yaml:"timeout"`
}

func (a *AdGuardHomeConfiguration) Validate() error {
	if a == nil {
		return fmt.Errorf("adguardgome configuration need to be configured")
	}
	if a.Username == "" || a.Password == "" {
		return fmt.Errorf("adguardgome username and password need to be configured")
	}
	if a.Timeout == 0 {
		a.Timeout = 5
	}
	return nil
}
