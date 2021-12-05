import { MessageApi } from "naive-ui";
import dayjs from "dayjs";

import HTTPError from "./http-error";

const oneHourMS = 3600 * 1000;
const oneDayMS = 24 * oneHourMS;

function formatError(err: Error | HTTPError | unknown): string {
  let message = "";
  if (err instanceof HTTPError) {
    message = err.message;
    if (err.category) {
      message += ` [${err.category.toUpperCase()}]`;
    }
    if (err.code) {
      message += ` [${err.code}]`;
    }
    // 如果是异常（客户端异常，如请求超时，中断等），则上报user action
    if (err.exception) {
      // const currentLocation = getCurrentLocation();
      // actionAdd({
      //   category: ERROR,
      //   route: currentLocation.name,
      //   path: currentLocation.path,
      //   result: FAIL,
      //   message,
      // });
    }
  } else if (err instanceof Error) {
    message = err.message;
  } else {
    message = (err as Error).message;
  }
  return message;
}

export function showError(
  message: MessageApi,
  err: Error | HTTPError | unknown
): void {
  message.error(formatError(err), {
    duration: 3000,
  });
}

export function showWarning(message: MessageApi, msg: string): void {
  message.warning(msg, {
    duration: 2000,
  });
}

export function toast(message: MessageApi, msg: string): void {
  message.info(msg, {
    duration: 2000,
  });
}

export function containsAny(
  data: readonly string[],
  target: string[]
): boolean {
  if (!data) {
    return false;
  }
  if (!target) {
    return true;
  }
  let exists = false;
  data.forEach((item) => {
    if (exists) {
      return;
    }
    exists = target.includes(item);
  });
  return exists;
}

// today 获取当天0点时间
export function today(): Date {
  return new Date(new Date(new Date().toLocaleDateString()).getTime());
}

// tomorrow 获取明天0点时间
export function tomorrow(): Date {
  return new Date(today().getTime() + oneDayMS);
}

// getDaysAgo 获取多少天前
export function getDaysAgo(days: number): Date {
  return new Date(today().getTime() - days * oneDayMS);
}

// getHoursAge 获取多少小时前
export function getHoursAge(hours: number): Date {
  return new Date(Date.now() - hours * oneHourMS);
}

// today 获取当天0点时间
export function yesterday(): Date {
  return getDaysAgo(1);
}

// formatDate 格式化日期
export function formatDate(str: string): string {
  return dayjs(str).format("YYYY-MM-DD HH:mm:ss");
}
// formatDateWithTZ 格式化日期（带时区）
export function formatDateWithTZ(date: Date): string {
  return dayjs(date).format("YYYY-MM-DDTHH:mm:ssZ");
}
export function formatBegin(begin: Date): string {
  return formatDateWithTZ(begin);
}
export function formatEnd(end: Date): string {
  return formatDateWithTZ(new Date(end.getTime() + 24 * 3600 * 1000 - 1));
}

interface DiffInfo {
  modifiedCount: number;
  data: Record<string, unknown>;
}
// eslint-disable-next-line
function isEqual(value: any, originalValue: any): boolean {
  // 使用json stringify对比是否相同
  return JSON.stringify(value) == JSON.stringify(originalValue);
}

// diff  对比两个object的差异
// eslint-disable-next-line
export function diff(
  current: Record<string, unknown>,
  original: Record<string, unknown>
): DiffInfo {
  const data: Record<string, unknown> = {};
  let modifiedCount = 0;
  Object.keys(current).forEach((key) => {
    const value = current[key];
    if (!isEqual(value, original[key])) {
      data[key] = value;
      modifiedCount++;
    }
  });
  return {
    modifiedCount,
    data,
  };
}

export function formatJSON(str: string): string {
  if (!str || str.length <= 2) {
    return str;
  }
  let result = str;
  const first = str[0];
  const last = str[str.length - 1];
  if ((first === "{" && last === "}") || (first === "[" && last === "]")) {
    try {
      result = JSON.stringify(JSON.parse(str), null, 2);
    } catch (err) {
      console.error(err);
    }
  }
  return result;
}

export function omitNil(
  data: Record<string, unknown>
): Record<string, unknown> {
  const result = {} as Record<string, unknown>;
  Object.keys(data).forEach((key) => {
    const value = data[key];
    if (value === null) {
      return;
    }
    result[key] = value;
  });
  return result;
}

export function durationToSeconds(d: string): number | null {
  if (!d || d.length < 2) {
    return null;
  }
  const units = ["s", "m", "h"];
  const seconds = [1, 60, 3600];
  const index = units.indexOf(d[d.length - 1]);
  if (index === -1) {
    return 0;
  }

  return Number.parseInt(d) * seconds[index];
}
