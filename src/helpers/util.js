import { isError } from "lodash-es";
import { env } from "@/config";
export function log(...args) {
  // eslint-disable-next-line
  console.info(...args);
}

export function debug(...args) {
  if (env !== "development") {
    return;
  }
  // eslint-disable-next-line
  console.debug(...args);
}

// 获取出错信息
export function getErrorMessage(err) {
  let message = err;
  if (err && err.response) {
    const { data, headers } = err.response;
    const id = headers["x-response-id"];
    if (data.code) {
      // eslint-disable-next-line
      const code = data.code.replace(`${app}-`, "");
      message = `${data.message}(${code}) [${id}]`;
    }
  }
  if (isError(message)) {
    message = message.message;
  }
  if (err.code === "ECONNABORTED") {
    message = "请求超时，请重新再试";
  }
  return message;
}

// formatDate
export function formatDate(str) {
  const date = new Date(str);
  const fill = v => {
    if (v >= 10) {
      return `${v}`;
    }
    return `0${v}`;
  };
  const month = fill(date.getMonth() + 1);
  const day = fill(date.getDate());
  const hours = fill(date.getHours());
  const mintues = fill(date.getMinutes());
  const seconds = fill(date.getSeconds());
  return `${date.getFullYear()}-${month}-${day} ${hours}:${mintues}:${seconds}`;
}
