import axios from "axios";

import { IMAGES_CONFIG, IMAGES_OPTIM, IMAGES_OPTIM_FROM_DATA } from "../urls";

// getConfig 获取图片的相关配置信息
export async function getConfig() {
  const { data } = await axios.get(IMAGES_CONFIG);
  return data;
}

// optim 压缩优化图片
export async function optim(file) {
  const { data } = await axios.get(IMAGES_OPTIM.replace(":file", file));
  return data;
}

// optimFromData 根据提交的图片数据压缩
export async function optimFromData(params) {
  const { data } = await axios.post(IMAGES_OPTIM_FROM_DATA, params);
  return data;
}
