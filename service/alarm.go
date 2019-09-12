// Copyright 2019 tree xie
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

package service

import (
	"crypto/tls"

	"github.com/vicanso/tiny-site/config"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

var (
	mailDialer *gomail.Dialer
	mailSender string
)

func init() {
	mailConfig := config.GetMailConfig()
	if mailConfig.Host != "" {
		mailSender = mailConfig.User
		mailDialer = gomail.NewDialer(mailConfig.Host, mailConfig.Port, mailConfig.User, mailConfig.Password)
		mailDialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
}

// AlarmError alarm error message
func AlarmError(message string) {
	logger.Error(message,
		zap.String("category", "alarm-error"),
	)
	if mailDialer != nil {
		m := gomail.NewMessage()
		receivers := config.GetStringSlice("alarm.receiver")
		m.SetHeader("From", mailSender)
		m.SetHeader("To", receivers...)
		m.SetHeader("Subject", "Alarm")
		m.SetBody("text/plain", message)
		// 避免发送邮件时太慢影响现有流程
		go func() {
			err := mailDialer.DialAndSend(m)
			if err != nil {
				logger.Error("send mail fail",
					zap.Error(err),
				)
			}
		}()
	}
}
