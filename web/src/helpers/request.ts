import axios, { AxiosRequestConfig, AxiosResponse } from "axios";
import { gzip } from "pako";

import HTTPError from "./http-error";
import { isDevelopment } from "../constants/env";
// import { httpRequests } from "../store";

const requestedAt = "X-Requested-At";
// 最小压缩长度
const compressMinLength = 10 * 1024;
const supportGzip = typeof TextEncoder !== "undefined";
const request = axios.create({
  // 默认超时为10秒
  timeout: 10 * 1000,
  transformRequest: [
    (data, header) => {
      if (!data || !header) {
        return;
      }
      header["Content-Type"] = "application/json;charset=UTF-8";
      const postData = JSON.stringify(data);
      // 如果数据较小或者不支持压缩
      if (postData.length < compressMinLength || !supportGzip) {
        return postData;
      }
      header["Content-Encoding"] = "gzip";
      return gzip(new TextEncoder().encode(postData));
    },
  ],
});

request.interceptors.request.use(
  (config) => {
    // 对请求的query部分清空值
    if (config.params) {
      Object.keys(config.params).forEach((element) => {
        // 空字符
        if (config.params[element] === "") {
          delete config.params[element];
        }
      });
    }
    if (isDevelopment()) {
      config.url = `/api${config.url}`;
    }
    if (config.headers) {
      config.headers[requestedAt] = `${Date.now()}`;
    }
    return config;
  },
  (err) => {
    return Promise.reject(err);
  }
);

// addRequestStats 添加http请求的相关记录
function addRequestStats(
  config: AxiosRequestConfig | undefined,
  res: AxiosResponse | undefined,
  he: HTTPError | undefined
): void {
  const data: Record<string, unknown> = {};
  if (config) {
    data.method = config.method;
    data.url = config.url;
    data.data = config.data;
    if (config.headers) {
      const value = config.headers[requestedAt];
      data.use = Date.now() - Number(value);
    }
  }
  if (res) {
    data.status = res.status;
  }
  if (he) {
    data.message = he.message;
  }
  // httpRequests.add(data);
}

// 设置接口最少要x ms才完成，能让客户看到loading
const minUse = 300;
const timeoutErrorCodes = ["ECONNABORTED", "ECONNREFUSED", "ECONNRESET"];
request.interceptors.response.use(
  async (res) => {
    addRequestStats(res.config, res, undefined);
    // 根据请求开始时间计算耗时，并判断是否需要延时响应
    if (res.config.headers) {
      const value = res.config.headers[requestedAt];
      if (value) {
        const use = Date.now() - Number(value);
        if (use >= 0 && use < minUse) {
          await new Promise((resolve) => setTimeout(resolve, minUse - use));
        }
      }
    }
    return res;
  },
  (err) => {
    const { response } = err;
    const he = new HTTPError(0, "请求出错");
    if (timeoutErrorCodes.includes(err.code)) {
      he.exception = true;
      he.code = err.code;
      he.category = "timeout";
      he.message = "请求超时，请稍候再试";
    } else if (response) {
      he.status = response.status;
      if (response.data && response.data.message) {
        he.message = response.data.message;
        he.code = response.data.code;
        he.category = response.data.category;
      } else {
        he.exception = true;
        he.category = "exception";
        he.message = `未知错误`;
      }
      he.extra = response.data?.extra;
    }
    addRequestStats(response?.config, response, he);
    return Promise.reject(he);
  }
);

export default request;
