package config

import "fmt"

type DomainConfiguration []string

func (d DomainConfiguration) Validate() error {
	if len(d) == 0 {
		fmt.Println("domain configuration is empty")
	}
	return nil
}

func (d DomainConfiguration) Filter() DomainConfiguration {
	if len(d) == 0 {
		return d
	}
	temp := make(DomainConfiguration, 0, len(d))
	for _, domain := range d {
		if isMyDomain(domain) {
			temp = append(temp, domain)
		}
	}
	return temp
}

func isMyDomain(domain string) bool {
	noDot := true
	purgeNumber := true
	for _, char := range domain {
		switch {
		case char >= 'A' && char <= 'Z':
			purgeNumber = false
		case char >= 'a' && char <= 'z':
			purgeNumber = false
		case char >= '0' && char <= '9':
		case (char == '.') && noDot:
			noDot = false
		}
	}
	return (!purgeNumber) && (!noDot)
}
