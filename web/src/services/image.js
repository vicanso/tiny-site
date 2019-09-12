import axios from "axios";

import { IMAGES_CONFIG, IMAGES_OPTIM } from "../urls";

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
