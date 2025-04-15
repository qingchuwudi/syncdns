package config

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var domainTests = []struct {
	Domain string
	Expect bool
}{
	{"", false},
	{"123", false},
	{".", false},
	{"..", false},
	{"localhost", false},
	{"localhost.localdomain", true},
	{"localhost.localdomain.local", true},
	{"127.0.0.1", true},
	{"baidu.com", true},

	{"Upper.baidu.com", true},
	{"780.baidu.com", true},
	{"中文域名", false},
	{"中文.域名", false},
}

var domainForTest = DomainConfiguration{
	"",
	"123",
	".",
	"..",
	"localhost",
	"localhost.localdomain",
	"localhost.localdomain.local",
	"127.0.0.1",
	"baidu.com",
	"Upper.baidu.com",
	"780.baidu.com",
	"中文域名",
	"中文.域名",
}

var validDomainForTest = DomainConfiguration{
	"localhost.localdomain",
	"localhost.localdomain.local",
	"baidu.com",
	"Upper.baidu.com",
	"780.baidu.com",
}

func TestDomainConfiguration_Validate(t *testing.T) {
	Convey("Test domain configuration Validate()", t, func() {
		ShouldBeNil(domainForTest.Validate())
		ShouldNotBeNil(DomainConfiguration{}.Validate())
	})
}

func TestDomainConfiguration_Filter(t *testing.T) {
	Convey("Test domain configuration Filter()", t, func() {
		ShouldBeEmpty(DomainConfiguration{}.Filter())
		ShouldEqual(domainForTest.Filter(), validDomainForTest)
	})
}

func TestIsMyDomain(t *testing.T) {
	Convey("Test domain check function", t, func() {
		for _, domain := range domainTests {
			ShouldEqual(isMyDomain(domain.Domain), domain.Expect)
		}
	})
}
