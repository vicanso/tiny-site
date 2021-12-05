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

package validate

func init() {
	// 用户账号
	AddAlias("xUserAccount", "ascii,min=2,max=10")
	// 用户密码
	AddAlias("xUserPassword", "ascii,len=44")
	// 用户名称
	AddAlias("xUserName", "min=1,max=20")
	// 用户邮箱
	AddAlias("xUserEmail", "email")
	// 用户角色
	AddAlias("xUserRole", "ascii,min=1,max=10")
	// 用户分组
	AddAlias("xUserGroup", "ascii,min=1,max=10")
	// 用户行为分类
	// TODO 是否调整为支持配置的方式
	Add("xUserActionCategory", newIsInString([]string{
		"click",
		"login",
		"register",
		"routeChange",
		"error",
	}))
	// 用户行为触发所在路由
	AddAlias("xUserActionRoute", "max=50")
}
