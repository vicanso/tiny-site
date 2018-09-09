package model

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/vicanso/tiny-site/util"

	"github.com/lib/pq"
)

var (
	// errUserAccountAndPwdNil account and passowrd can not be nil
	errUserAccountAndPwdNil = &util.HTTPError{
		StatusCode: http.StatusBadRequest,
		Category:   util.ErrCategoryLogic,
		Code:       util.ErrCodeUser,
		Message:    "account and password can not be nil",
	}
	// errUserAccountExists account is exists
	errUserAccountExists = &util.HTTPError{
		StatusCode: http.StatusBadRequest,
		Category:   util.ErrCategoryLogic,
		Code:       util.ErrCodeUser,
		Message:    "account already exists",
	}
	// errUserAccountNotExists account is not exists
	errUserAccountNotExists = &util.HTTPError{
		StatusCode: http.StatusBadRequest,
		Category:   util.ErrCategoryLogic,
		Code:       util.ErrCodeUser,
		Message:    "account is not exists",
	}
)

const (
	// UserRoleSu super user
	UserRoleSu = "su"
	// UserRoleAdmin admin user
	UserRoleAdmin = "admin"
)

const (
	// UserActionAddRole add role
	UserActionAddRole = iota + 1
	// UserActionRemoveRole remove role
	UserActionRemoveRole
)

type (
	// User user model
	User struct {
		BaseModel
		Account  string         `json:"account,omitempty" gorm:"type:varchar(20);not null;unique_index:idx_users_account"`
		Password string         `json:"password,omitempty" gorm:"not null;"`
		Roles    pq.StringArray `json:"roles,omitempty" gorm:"type:text[]"`
		client   *gorm.DB
	}
	// UserLogin user login
	UserLogin struct {
		BaseModel
		Account   string `json:"account,omitempty" gorm:"index:idx_user_logins_account"`
		UserAgent string `json:"userAgent,omitempty"`
		IP        string `json:"ip,omitempty"`
		TrackID   string `json:"trackId,omitempty"`
		SessionID string `json:"sessionId,omitempty"`
	}
)

// Exists check the user exists
func (u *User) Exists() (exists bool, err error) {
	client := u.client
	if client == nil {
		client = GetClient()
	}
	result := &User{}
	err = client.Where(u).First(result).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	err = nil
	if result.ID != 0 {
		exists = true
	}
	return
}

// Save save user
func (u *User) Save() (err error) {
	if u.Account == "" || u.Password == "" {
		err = errUserAccountAndPwdNil
		return
	}
	client := u.client
	if client == nil {
		client = GetClient()
	}
	exists, err := (&User{
		Account: u.Account,
	}).Exists()
	if err != nil {
		return
	}
	if exists {
		err = errUserAccountExists
		return
	}
	err = client.Create(u).Error
	if err != nil {
		return
	}
	// the first account add `su`
	if u.ID == 1 {
		go client.Model(u).Update(User{
			Roles: []string{
				UserRoleSu,
			},
		})
	}
	return
}

// First get the first match record
func (u *User) First() (err error) {
	client := u.client
	if client == nil {
		client = GetClient()
	}
	err = client.Where(u).First(u).Error
	return
}

// UpdateRole update user role the role
func (u *User) UpdateRole(role string, action int) (err error) {
	err = u.First()
	if err != nil {
		return
	}
	if u.ID == 0 {
		err = errUserAccountNotExists
		return
	}
	roles := []string{}
	for _, v := range u.Roles {
		if v != role {
			roles = append(roles, v)
		}
	}
	if action == UserActionAddRole {
		roles = append(roles, role)
	}
	client := GetClient()
	err = client.Model(u).Update(User{
		Roles: roles,
	}).Error

	return
}

// Save save the login record
func (ul *UserLogin) Save() (err error) {
	client := GetClient()
	err = client.Create(ul).Error
	return
}
