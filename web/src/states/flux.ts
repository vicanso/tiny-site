import request from "../helpers/request";
import {
  FLUXES_TRACKERS,
  FLUXES_HTTP_ERRORS,
  FLUXES_TAG_VALUES,
  FLUXES_REQUESTS,
  FLUXES_FIND_ONE,
} from "../constants/url";
import { DeepReadonly, reactive, readonly } from "vue";
import { formatDate } from "../helpers/util";
import { IList } from "./interface";

export const measurementUserTracker = "userTracker";
export const measurementHttpRequest = "httpRequest";
export const measurementHttpError = "httpError";

function sortByTime(
  item1: {
    _time: string;
  },
  item2: {
    _time: string;
  }
) {
  if (item1._time === item2._time) {
    return 0;
  }
  if (item1._time > item2._time) {
    return -1;
  }
  return 1;
}

// 用户行为轨迹
interface UserTracker {
  [key: string]: unknown;
  _time: string;
  key: string;
  createdAt: string;
  account: string;
  action: string;
  hostname: string;
  ip: string;
  result: string;
  resultDesc: string;
  sid: string;
  tid: string;
  form: string;
  query: string;
  params: string;
  error: string;
}

interface UserTrackers extends IList<UserTracker> {
  flux: string;
}
const userTrackers: UserTrackers = reactive({
  processing: false,
  items: [],
  count: -1,
  flux: "",
});

// 用户行为轨迹类型
const userTrackerActions: IList<string> = reactive({
  processing: false,
  items: [],
  count: -1,
});

function fillUserTrackerInfo(data: UserTracker) {
  if (data.error) {
    const reg = /, message=([\s\S]*)/;
    const result = reg.exec(data.error);
    if (result && result.length === 2) {
      data.error = `${result[1]}, ${data.error.replace(result[0], "")}`;
    }
  }
  if (data.result === "0") {
    data.resultDesc = "成功";
  } else {
    data.resultDesc = "失败";
  }
  data.key = data._time;
  data.createdAt = formatDate(data._time);
}

// HTTP出错类型
const httpErrorCategories: IList<string> = reactive({
  processing: false,
  items: [],
  count: -1,
});

// HTTPError 客户端HTTP请求出错记录
interface HTTPError {
  _time: string;
  createdAt: string;
  key: string;
  account: string;
  category: string;
  error: string;
  exception: boolean;
  hostname: string;
  ip: string;
  method: string;
  route: string;
  sid: string;
  status: number;
  tid: string;
  uri: string;
}
interface HTTPErrors extends IList<HTTPError> {
  flux: string;
}
const httpErrors: HTTPErrors = reactive({
  processing: false,
  items: [],
  count: -1,
  flux: "",
});

function fillHTTPErrorInfo(data: HTTPError) {
  data.key = data._time;
  data.createdAt = formatDate(data._time);
}

// 后端HTTP请求记录
interface Request {
  _time: string;
  key: string;
  createdAt: string;
  hostname: string;
  addr: string;
  service: string;
  method: string;
  route: string;
  uri: string;
  status: number;
  reused: boolean;
  dnsUse: number;
  tcpUse: number;
  tlsUse: number;
  processingUse: number;
  use: number;
  result: string;
  errCategory: string;
  error: string;
  exception: boolean;
}
interface Requests extends IList<Request> {
  flux: string;
}
const requests: Requests = reactive({
  processing: false,
  items: [],
  count: -1,
  flux: "",
});

// request 服务名称
const requestServices: IList<string> = reactive({
  processing: false,
  items: [],
  count: -1,
});

// RequestRoutes 请求路由
const requestRoutes: IList<string> = reactive({
  processing: false,
  items: [],
  count: -1,
});

function fillRequestInfo(data: Request) {
  data.key = data._time;
  data.createdAt = formatDate(data._time);
}

// fluxListUserTracker 查询用户跟踪轨迹记录
export async function fluxListUserTracker(params: {
  account?: string;
  action?: string;
  begin: string;
  end: string;
  limit: number;
  result?: string;
}): Promise<void> {
  if (userTrackers.processing) {
    return;
  }
  try {
    userTrackers.processing = true;
    const { data } = await request.get<{
      trackers: UserTracker[];
      count: number;
      flux: string;
    }>(FLUXES_TRACKERS, {
      params,
    });
    userTrackers.items = data.trackers || [];
    userTrackers.items.sort(sortByTime);
    userTrackers.count = data.count || 0;
    userTrackers.flux = data.flux || "";
    userTrackers.items.forEach(fillUserTrackerInfo);
  } finally {
    userTrackers.processing = false;
  }
}

// fluxListUserTrackAction 查询用户轨迹action列表
export async function fluxListUserTrackAction(): Promise<void> {
  if (userTrackerActions.processing || userTrackerActions.items.length !== 0) {
    return;
  }
  try {
    userTrackerActions.processing = true;
    const url = FLUXES_TAG_VALUES.replace(
      ":measurement",
      measurementUserTracker
    ).replace(":tag", "action");
    const { data } = await request.get<{
      values: string[];
    }>(url);
    userTrackerActions.items = (data.values || []).sort();
  } finally {
    userTrackerActions.processing = false;
  }
}

// fluxListUserTrackerClear 清除tracker记录
export function fluxListUserTrackerClear(): void {
  userTrackers.items.length = 0;
  userTrackers.count = -1;
  userTrackers.flux = "";
}

// fluxListHTTPCategory 查询HTTP出错类型列表
export async function fluxListHTTPCategory(): Promise<void> {
  if (
    httpErrorCategories.processing ||
    httpErrorCategories.items.length !== 0
  ) {
    return;
  }
  try {
    httpErrorCategories.processing = true;
    const url = FLUXES_TAG_VALUES.replace(
      ":measurement",
      measurementHttpError
    ).replace(":tag", "category");
    const { data } = await request.get<{
      values: string[];
    }>(url);
    httpErrorCategories.items = (data.values || []).sort();
  } finally {
    httpErrorCategories.processing = false;
  }
}

// fluxListHTTPError 查询HTTP出错记录
export async function fluxListHTTPError(params: {
  account?: string;
  category?: string;
  begin: string;
  end: string;
  exception?: string;
  limit: number;
}): Promise<void> {
  if (httpErrors.processing) {
    return;
  }
  try {
    httpErrors.processing = true;
    const { data } = await request.get<{
      httpErrors: HTTPError[];
      count: number;
      flux: string;
    }>(FLUXES_HTTP_ERRORS, {
      params,
    });
    httpErrors.items = data.httpErrors || [];
    httpErrors.count = data.count || 0;
    httpErrors.flux = data.flux || "";
    httpErrors.items.forEach(fillHTTPErrorInfo);
    httpErrors.items.sort(sortByTime);
  } finally {
    httpErrors.processing = false;
  }
}

// fluxListHTTPErrorClear 清除http出错列表
export function fluxListHTTPErrorClear(): void {
  httpErrors.items.length = 0;
  httpErrors.count = -1;
  httpErrors.flux = "";
}

// fluxListRequest 查询后端请求记录
export async function fluxListRequest(params: {
  route?: string;
  service?: string;
  errCategory?: string;
  begin: string;
  end: string;
  exception?: string;
  limit: number;
}): Promise<void> {
  if (requests.processing) {
    return;
  }
  try {
    requests.processing = true;
    const { data } = await request.get<{
      requests: Request[];
      count: number;
      flux: string;
    }>(FLUXES_REQUESTS, {
      params,
    });
    requests.items = data.requests || [];
    requests.count = data.count || 0;
    requests.flux = data.flux || "";
    requests.items.forEach(fillRequestInfo);
    requests.items.sort(sortByTime);
  } finally {
    requests.processing = false;
  }
}

// fluxListRequestClear 清除request列表
export function fluxListRequestClear(): void {
  requests.items.length = 0;
  requests.count = -1;
  requests.flux = "";
}

// fluxListRequestService 获取request中的service列表
export async function fluxListRequestService(): Promise<void> {
  if (requestServices.processing || requestServices.items.length !== 0) {
    return;
  }
  try {
    requestServices.processing = true;
    const url = FLUXES_TAG_VALUES.replace(
      ":measurement",
      measurementHttpRequest
    ).replace(":tag", "service");
    const { data } = await request.get<{
      values: string[];
    }>(url);
    requestServices.items = data.values || [];
  } finally {
    requestServices.processing = false;
  }
}

// fluxListRequestRoute 获取request中的route列表
export async function fluxListRequestRoute(): Promise<void> {
  if (requestRoutes.processing || requestRoutes.items.length !== 0) {
    return;
  }
  try {
    requestRoutes.processing = true;
    const url = FLUXES_TAG_VALUES.replace(
      ":measurement",
      measurementHttpRequest
    ).replace(":tag", "route");
    const { data } = await request.get<{
      values: string[];
    }>(url);
    requestRoutes.items = data.values || [];
  } finally {
    requestRoutes.processing = false;
  }
}

// flux查询单条记录
export async function fluxFindOne(params: {
  measurement: string;
  time: string;
  tags: Record<string, string>;
}): Promise<Record<string, unknown>> {
  const url = FLUXES_FIND_ONE.replace(":measurement", params.measurement);
  const { data } = await request.get<Record<string, unknown>>(url, {
    params: Object.assign(
      {
        time: params.time,
      },
      params.tags
    ),
  });
  return data;
}

interface ReadonlyFluxState {
  userTrackers: DeepReadonly<UserTrackers>;
  userTrackerActions: DeepReadonly<IList<string>>;
  httpErrors: DeepReadonly<HTTPErrors>;
  httpErrorCategories: DeepReadonly<IList<string>>;
  requests: DeepReadonly<Requests>;
  requestServices: DeepReadonly<IList<string>>;
  requestRoutes: DeepReadonly<IList<string>>;
}

const state = {
  userTrackers: readonly(userTrackers),
  userTrackerActions: readonly(userTrackerActions),
  httpErrors: readonly(httpErrors),
  httpErrorCategories: readonly(httpErrorCategories),
  requests: readonly(requests),
  requestServices: readonly(requestServices),
  requestRoutes: readonly(requestRoutes),
};

export default function useFluxState(): ReadonlyFluxState {
  return state;
}
