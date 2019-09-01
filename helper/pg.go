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

package helper

import (
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/vicanso/tiny-site/config"
	"github.com/vicanso/tiny-site/log"
	"go.uber.org/zap"
)

var (
	pgClient *gorm.DB
)

func init() {
	str := config.GetPostgresConnectString()
	reg := regexp.MustCompile(`password=\S*`)
	maskStr := reg.ReplaceAllString(str, "password=***")
	logger.Info("connect to pg",
		zap.String("args", maskStr),
	)
	db, err := gorm.Open("postgres", str)
	if err != nil {
		panic(err)
	}
	db.SetLogger(log.PGLogger())
	pgClient = db
}

// PGCreate pg create
func PGCreate(data interface{}) (err error) {
	err = pgClient.Create(data).Error
	return
}

// PGGetClient pg client
func PGGetClient() *gorm.DB {
	return pgClient
}

// PGFormatSort format sort
func PGFormatSort(sort string) string {
	arr := strings.Split(sort, ",")
	newSort := []string{}
	for _, item := range arr {
		if item[0] == '-' {
			newSort = append(newSort, item[1:]+" desc")
		} else {
			newSort = append(newSort, item)
		}
	}
	return strings.Join(newSort, ",")
}
