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
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/vicanso/elton"
	session "github.com/vicanso/elton-session"
	"github.com/vicanso/hes"
	"github.com/vicanso/tiny-site/cs"
	"github.com/vicanso/tiny-site/util"

	"go.uber.org/zap"
)

const (
	// UserAccount user account field
	UserAccount = "account"
	// UserLoginAt user login at
	UserLoginAt = "loginAt"
	// UserRoles user roles
	UserRoles = "roles"
	// UserLoginToken user login token
	UserLoginToken = "loginToken"
)

const (
	defaultUserLimit            = 10
	defaultUserLoginRecordLimit = 10
)

var (
	// admin用户角色
	adminUserRoles = []string{
		cs.UserRoleSu,
		cs.UserRoleAdmin,
	}
)

var (
	errAccountOrPasswordInvalid = hes.New("account or password is invalid")
)

type (
	// UserSession user session struct
	UserSession struct {
		se *session.Session
	}
	// User user
	User struct {
		ID        uint       `gorm:"primary_key" json:"id,omitempty"`
		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt time.Time  `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`

		Account  string         `json:"account,omitempty" gorm:"type:varchar(20);not null;unique_index:idx_users_account"`
		Password string         `json:"-" gorm:"type:varchar(128);not null;"`
		Roles    pq.StringArray `json:"roles,omitempty" gorm:"type:text[]"`
	}

	// UserLoginRecord user login
	UserLoginRecord struct {
		ID        uint       `gorm:"primary_key" json:"id,omitempty"`
		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt time.Time  `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`

		Account       string `json:"account,omitempty" gorm:"type:varchar(20);not null;index:idx_user_logins_account"`
		UserAgent     string `json:"userAgent,omitempty"`
		IP            string `json:"ip,omitempty" gorm:"type:varchar(64);not null"`
		TrackID       string `json:"trackId,omitempty" gorm:"type:varchar(64);not null"`
		SessionID     string `json:"sessionId,omitempty" gorm:"type:varchar(64);not null"`
		XForwardedFor string `json:"xForwardedFor,omitempty" gorm:"type:varchar(128)"`
		Country       string `json:"country,omitempty" gorm:"type:varchar(64)"`
		Province      string `json:"province,omitempty" gorm:"type:varchar(64)"`
		City          string `json:"city,omitempty" gorm:"type:varchar(64)"`
		ISP           string `json:"isp,omitempty" gorm:"type:varchar(64)"`
	}
	// UserTrackRecord user track record
	UserTrackRecord struct {
		ID        uint       `gorm:"primary_key" json:"id,omitempty"`
		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt time.Time  `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`
		TrackID   string     `json:"trackId,omitempty" gorm:"type:varchar(64);not null;index:idx_user_track_id"`
		UserAgent string     `json:"userAgent,omitempty"`
		IP        string     `json:"ip,omitempty" gorm:"type:varchar(64);not null"`
		Country   string     `json:"country,omitempty" gorm:"type:varchar(64)"`
		Province  string     `json:"province,omitempty" gorm:"type:varchar(64)"`
		City      string     `json:"city,omitempty" gorm:"type:varchar(64)"`
		ISP       string     `json:"isp,omitempty" gorm:"type:varchar(64)"`
	}
	// UserQueryParams user query params
	UserQueryParams struct {
		Keyword string
		Role    string
		Limit   int
	}
	// UserLoginRecordQueryParams user login record query params
	UserLoginRecordQueryParams struct {
		Begin   string
		End     string
		Account string
		Limit   int
		Offset  int
	}
	// UserSrv user service
	UserSrv struct {
	}
)

func init() {
	pgGetClient().AutoMigrate(&User{}).
		AutoMigrate(&UserLoginRecord{}).
		AutoMigrate(&UserTrackRecord{})
}

// Add add user
func (srv *UserSrv) Add(u *User) (err error) {
	err = pgCreate(u)
	// 首次创建账号，设置su权限
	if u.ID == 1 {
		pgGetClient().Model(u).Update(map[string]interface{}{
			"roles": []string{
				cs.UserRoleSu,
			},
		})
	}
	return
}

// Login user login
func (srv *UserSrv) Login(account, password, token string) (u *User, err error) {
	u = &User{}
	err = pgGetClient().Where("account = ?", account).First(u).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = errAccountOrPasswordInvalid
		}
		return
	}
	pwd := util.Sha256(u.Password + token)
	// 用于自动化测试使用
	if util.IsDevelopment() && password == "fEqNCco3Yq9h5ZUglD3CZJT4lBsfEqNCco31Yq9h5ZUB" {
		pwd = password
	}
	if pwd != password {
		err = errAccountOrPasswordInvalid
		return
	}
	return
}

// Update update user
func (srv *UserSrv) Update(user *User, attrs ...interface{}) (err error) {
	err = pgGetClient().Model(user).Update(attrs...).Error
	return
}

// AddLoginRecord add user login record
func (srv *UserSrv) AddLoginRecord(r *UserLoginRecord) (err error) {
	err = pgCreate(r)
	if r.ID != 0 {
		id := r.ID
		ip := r.IP
		go func() {
			lo, err := GetLocationByIP(ip, nil)
			if err != nil {
				logger.Error("get location by ip fail",
					zap.String("ip", ip),
					zap.Error(err),
				)
				return
			}
			pgGetClient().Model(&UserLoginRecord{
				ID: id,
			}).Update(map[string]string{
				"country":  lo.Country,
				"province": lo.Province,
				"city":     lo.City,
				"isp":      lo.ISP,
			})
		}()
	}
	return
}

// AddTrackRecord add track record
func (srv *UserSrv) AddTrackRecord(r *UserTrackRecord) (err error) {
	// TODO 后续写入influxdb，避免被攻击而产生大量的无用记录
	err = pgCreate((r))
	if r.ID != 0 {
		id := r.ID
		ip := r.IP
		go func() {
			lo, err := GetLocationByIP(ip, nil)
			if err != nil {
				logger.Error("get location by ip fail",
					zap.String("ip", ip),
					zap.Error(err),
				)
				return
			}
			pgGetClient().Model(&UserTrackRecord{
				ID: id,
			}).Update(map[string]string{
				"country":  lo.Country,
				"province": lo.Province,
				"city":     lo.City,
				"isp":      lo.ISP,
			})
		}()
	}
	return
}

// List list users
func (srv *UserSrv) List(params UserQueryParams) (result []*User, err error) {
	result = make([]*User, 0)
	db := pgGetClient()
	if params.Limit <= 0 {
		db = db.Limit(defaultUserLimit)
	} else {
		db = db.Limit(params.Limit)
	}
	if params.Role != "" {
		db = db.Where("? = ANY(roles)", params.Role)
	}
	if params.Keyword != "" {
		db = db.Where("account LIKE ?", "%"+params.Keyword+"%")
	}
	err = db.Find(&result).Error
	return
}

func newUserLoginRecordQuery(params UserLoginRecordQueryParams) *gorm.DB {
	db := pgGetClient()
	if params.Account != "" {
		db = db.Where("account = ?", params.Account)
	}
	if params.Limit <= 0 {
		db = db.Limit(defaultUserLoginRecordLimit)
	} else {
		db = db.Limit(params.Limit)
	}
	if params.Begin != "" {
		db = db.Where("created_at > ?", params.Begin)
	}
	if params.End != "" {
		db = db.Where("created_at < ?", params.End)
	}
	return db
}

// ListLoginRecord list login record
func (srv *UserSrv) ListLoginRecord(params UserLoginRecordQueryParams) (result []*UserLoginRecord, err error) {
	result = make([]*UserLoginRecord, 0)
	db := newUserLoginRecordQuery(params)
	if params.Offset > 0 {
		db = db.Offset(params.Offset)
	}
	err = db.Find(&result).Error
	return
}

// CountLoginRecord count login record
func (srv *UserSrv) CountLoginRecord(params UserLoginRecordQueryParams) (count int, err error) {
	count = 0
	db := newUserLoginRecordQuery(params)
	err = db.Model(&UserLoginRecord{}).Count(&count).Error
	return
}

// GetAccount get the account
func (u *UserSession) GetAccount() string {
	if u.se == nil {
		return ""
	}
	return u.se.GetString(UserAccount)
}

// SetAccount set the account
func (u *UserSession) SetAccount(account string) error {
	return u.se.Set(context.Background(), UserAccount, account)
}

// GetUpdatedAt get updated at
func (u *UserSession) GetUpdatedAt() string {
	return u.se.GetUpdatedAt()
}

// SetLoginAt set the login at
func (u *UserSession) SetLoginAt(date string) error {
	return u.se.Set(context.Background(), UserLoginAt, date)
}

// GetLoginAt get login at
func (u *UserSession) GetLoginAt() string {
	return u.se.GetString(UserLoginAt)
}

// SetLoginToken get user login token
func (u *UserSession) SetLoginToken(token string) error {
	return u.se.Set(context.Background(), UserLoginToken, token)
}

// GetLoginToken get user login token
func (u *UserSession) GetLoginToken() string {
	return u.se.GetString(UserLoginToken)
}

// SetRoles set user roles
func (u *UserSession) SetRoles(roles []string) error {
	return u.se.Set(context.Background(), UserRoles, roles)
}

// GetRoles get user roles
func (u *UserSession) GetRoles() []string {
	result, ok := u.se.Get(UserRoles).([]interface{})
	if !ok {
		return nil
	}
	roles := []string{}
	for _, item := range result {
		role, _ := item.(string)
		if role != "" {
			roles = append(roles, role)
		}
	}
	return roles
}

// Destroy destroy user session
func (u *UserSession) Destroy() error {
	return u.se.Destroy(context.Background())
}

// Refresh refresh user session
func (u *UserSession) Refresh() error {
	return u.se.Refresh(context.Background())
}

// ClearSessionID clear session id
func (u *UserSession) ClearSessionID() {
	u.se.ID = ""
}

// IsAdmin check user is admin
func (u *UserSession) IsAdmin() bool {
	return util.UserRoleIsValid(adminUserRoles, u.GetRoles())
}

// NewUserSession create a user session
func NewUserSession(c *elton.Context) *UserSession {
	v, _ := c.Get(session.Key)
	if v == nil {
		return nil
	}
	data, _ := c.Get(cs.UserSession)
	if data != nil {
		us, ok := data.(*UserSession)
		if ok {
			return us
		}
	}
	us := &UserSession{
		se: v.(*session.Session),
	}
	c.Set(cs.UserSession, us)

	return us
}
