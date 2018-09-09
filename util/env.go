package util

import "github.com/vicanso/tiny-site/config"

const (
	// envProduction production env
	envProduction = "production"
	// envTest test env
	envTest = "test"
)

// IsDevelopment 判断是否开发环境
func IsDevelopment() bool {
	return config.GetENV() == ""
}

// IsTest 判断是否测试环境
func IsTest() bool {
	return config.GetENV() == envTest
}

// IsProduction 判断是否生产环境
func IsProduction() bool {
	return config.GetENV() == envProduction
}
