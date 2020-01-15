import axios from "axios";

import {
  FIELS_DOWNLOAD,
  FIELS_UPLOAD_UPDATE,
  FILES_DETAIL,
  FILES_ZONES_LIST_MINE,
  FILES_ZONES_ADD,
  FILES_ZONES_UPDATE,
  FILES_ZONES_LIST,
  FILES_LIST,
  FILES_UPLOAD_SAVE
} from "../urls";

// addZone 添加文件空间
export async function addZone(params) {
  const { data } = await axios.post(FILES_ZONES_ADD, params);
  return data;
}

// listMyZone 获取我的文件空间列表
export async function listMyZone() {
  const { data } = await axios.get(FILES_ZONES_LIST_MINE);
  return data.fileZones || [];
}

// listZone 文件空间列表
export async function listZone() {
  const { data } = await axios.get(FILES_ZONES_LIST);
  return data.fileZones || [];
}

// updateZone 更新文件空间
export async function updateZone(params) {
  const url = FILES_ZONES_UPDATE.replace(":fileZoneID", params.id);
  delete params.id;
  const { data } = await axios.patch(url, params);
  return data;
}

// list 获取文件列表
export async function list(params) {
  const { data } = await axios.get(FILES_LIST, {
    params
  });
  return data;
}

// saveFile 上传保存文件
export async function saveFile(params) {
  const { data } = await axios.post(FILES_UPLOAD_SAVE, params);
  return data;
}

// updateFile 更新文件信息
export async function updateFile(params) {
  const url = FIELS_UPLOAD_UPDATE.replace(":fileID", params.id);
  delete params.id;
  const { data } = await axios.patch(url, params);
  return data;
}

// getByID 获取文件信息
export async function getByID(id, params) {
  const { data } = await axios.get(FILES_DETAIL.replace(":fileID", id), {
    params
  });
  return data;
}

// downloadFile 下载文件
export async function downloadFile(fileURL) {
  const { data } = await axios.get(
    FIELS_DOWNLOAD + "?file=" + encodeURIComponent(fileURL)
  );
  return data;
}
