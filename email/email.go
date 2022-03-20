// Copyright 2020 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package email

import (
	"context"
	"crypto/tls"
	"strings"
	"sync"

	"github.com/vicanso/tiny-site/config"
	"github.com/vicanso/tiny-site/log"
	"gopkg.in/gomail.v2"
)

var (
	sendingMailMutex = &sync.Mutex{}
	newMailOnce      = &sync.Once{}

	currentEmailListRMutex = &sync.RWMutex{}
	// 保存当前邮箱列表
	currentEmailList map[string][]string
)

var (
	defaultMailDialer *gomail.Dialer
)

var (
	basicInfo  = config.MustGetBasicConfig()
	mailConfig = config.MustGetMailConfig()
)

// 更新邮箱列表
func Update(data map[string]string) {
	currentEmailListRMutex.Lock()
	defer currentEmailListRMutex.Unlock()
	currentEmailList = make(map[string][]string)
	for key, item := range data {
		currentEmailList[key] = strings.Split(item, ",")
	}
}

// 根据名称获取邮件列表
func List(name string) []string {
	currentEmailListRMutex.RLock()
	defer currentEmailListRMutex.RUnlock()
	emails := currentEmailList[name]
	return emails
}

// newMailDialer 新建邮件发送dialer
func newMailDialer() *gomail.Dialer {
	newMailOnce.Do(func() {
		if mailConfig.Host == "" {
			return
		}
		d := gomail.NewDialer(mailConfig.Host, mailConfig.Port, mailConfig.User, mailConfig.Password)
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		defaultMailDialer = d
	})
	return defaultMailDialer

}

// AlarmError 发送出错警告
func Send(ctx context.Context, title, message string, receivers ...string) {
	d := newMailDialer()

	if d != nil && len(receivers) != 0 {
		m := gomail.NewMessage()
		m.SetHeader("From", mailConfig.User)
		m.SetHeader("To", receivers...)
		m.SetHeader("Subject", title)
		m.SetBody("text/plain", message)
		// 避免发送邮件时太慢影响现有流程
		go func() {
			// 一次只允许一个email发送（由于使用的邮件服务有限制）
			sendingMailMutex.Lock()
			defer sendingMailMutex.Unlock()
			err := d.DialAndSend(m)
			if err != nil {
				log.Error(ctx).
					Err(err).
					Msg("send mail fail")
			}
		}()
	}
}

// AlarmError 发送出错警告
func AlarmError(ctx context.Context, message string) {
	log.Error(ctx).
		Str("app", basicInfo.Name).
		Str("category", "alarmError").
		Msg(message)
	receivers := List("alarmReceivers")
	Send(ctx, "Alarm-"+basicInfo.Name, message, receivers...)
}
