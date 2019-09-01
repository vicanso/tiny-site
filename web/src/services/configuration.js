import axios from "axios";

import {
  CONFIGURATIONS_LIST,
  CONFIGURATIONS_LIST_AVAILABLE,
  CONFIGURATIONS_LIST_UNAVAILABLE,
  CONFIGURATIONS_ADD,
  CONFIGURATIONS_UPDATE,
  CONFIGURATIONS_DELETE
} from "../urls";

// list 列出所有配置
export async function list(params) {
  const { data } = await axios.get(CONFIGURATIONS_LIST, {
    params
  });
  const configs = data.configs || [];
  return configs;
}

// listAvaiable 列出当前有效配置
export async function listAvaiable(params) {
  const { data } = await axios.get(CONFIGURATIONS_LIST_AVAILABLE, {
    params
  });
  const configs = data.configs || [];
  return configs;
}

// listUnavaiable 列出当前失效配置
export async function listUnavaiable(params) {
  const { data } = await axios.get(CONFIGURATIONS_LIST_UNAVAILABLE, {
    params
  });
  const configs = data.configs || [];
  return configs;
}

// add 添加配置
export async function add(params) {
  const { data } = await axios.post(CONFIGURATIONS_ADD, params);
  return data;
}

// updateByID 更新配置
export async function updateByID(id, params) {
  const url = CONFIGURATIONS_UPDATE.replace(":id", id);
  const { data } = await axios.patch(url, params);
  return data;
}

// deleteByID 删除配置
export async function deleteByID(id) {
  const url = CONFIGURATIONS_DELETE.replace(":id", id);
  const { data } = await axios.delete(url);
  return data;
}
