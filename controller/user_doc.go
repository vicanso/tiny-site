// +build swagger
// 用户相关接口文档

package controller

import "github.com/vicanso/tiny-site/ent"

// 用户列表响应
// swagger:response apiUserListResponse
type apiUserListResponse struct {
	// in: body
	Body *userListResp
}

// 用户信息查询参数
// swagger:parameters userList
type apiUserListParams struct {
	userListParams
}

// 用户登录Token响应
// swagger:response apiUserLoginTokenResponse
type apiUserLoginTokenResponse struct {
	// in: body
	Body *userLoginTokenResp
}

// 用户信息
// swagger:response apiUserInfoResponse
type apiUserInfoResponse struct {
	// in: body
	Body *userInfoResp
}

// 用户登录与注册参数
// swagger:parameters userRegister userLogin
type apiUserRegisterLoginParams struct {
	// in: body
	Body *userRegisterLoginParams
}

// 用户注册响应
// swagger:response apiUserRegisterResponse
type apiUserRegisterResponse struct {
	// in: body
	Body *ent.User
}
