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
    if (data.message) {
      message = data.message;
    }
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

// copy copy the value
export function copy(value) {
  // 来源自：https://juejin.im/post/5a94f8eff265da4e9b593c29
  const input = document.createElement("input");
  input.setAttribute("readonly", "readonly");
  input.setAttribute("value", value);
  document.body.appendChild(input);
  input.select();
  input.setSelectionRange(0, 9999);
  document.execCommand("copy");
  document.body.removeChild(input);
}
