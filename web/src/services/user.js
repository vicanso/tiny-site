import axios from "axios";

import {
  USERS,
  USERS_UPDATE,
  USERS_LOGOUT,
  USERS_LOGIN,
  USERS_ME,
  USERS_LOGIN_RECORDS
} from "../urls";

// logout 退出登录
export async function logout() {
  const { data } = await axios.delete(USERS_LOGOUT);
  return data;
}

// getLoginToken 获取登录token
export async function getLoginToken() {
  const { data } = await axios.get(USERS_LOGIN);
  return data;
}

// login 登录
export async function login(params) {
  const { captcha } = params;
  delete params.captcha;
  const { data } = await axios.post(USERS_LOGIN, params, {
    headers: {
      "X-Captcha": captcha
    }
  });
  return data;
}

// register 注册
export async function register(params) {
  const { captcha } = params;
  delete params.captcha;
  const { data } = await axios.post(USERS_ME, params, {
    headers: {
      "X-Captcha": captcha
    }
  });
  return data;
}

// list 列出用户列表
export async function list(params) {
  const { data } = await axios.get(USERS, {
    params
  });
  return data;
}

// updateByID 通过ID更新用户信息
export async function updateByID(id, params) {
  const url = USERS_UPDATE.replace(":id", id);
  const { data } = await axios.patch(url, params);
  return data;
}

// listLoginRecords 获取登录记录
export async function listLoginRecords(params) {
  const { data } = await axios.get(USERS_LOGIN_RECORDS, {
    params
  });
  return data;
}
