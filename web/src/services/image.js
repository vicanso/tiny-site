import axios from "axios";

import { IMAGES_CONFIG } from "../urls";

// getConfig 获取图片的相关配置信息
export async function getConfig() {
  const { data } = await axios.get(IMAGES_CONFIG);
  return data;
}
