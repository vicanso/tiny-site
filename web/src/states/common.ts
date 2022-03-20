import { DeepReadonly, reactive, readonly } from "vue";
import {
  COMMONS_CAPTCHA,
  COMMONS_ROUTERS,
  COMMONS_HTTP_STATS,
} from "../constants/url";
import request from "../helpers/request";
import { settingStorage } from "../storages/local";
interface Captcha {
  data: string;
  expiredAt: string;
  id: string;
  type: string;
}

export enum Mode {
  Add = "add",
  Update = "update",
  List = "list",
}

interface Settings {
  theme: string;
  collapsed: boolean;
}

const settings: Settings = reactive({
  theme: "dark",
  collapsed: false,
});

// 路由配置
interface Router {
  method: string;
  route: string;
}
interface Routers {
  processing: boolean;
  items: Router[];
}
const routers: Routers = reactive({
  processing: false,
  items: [],
});

// http实例
interface RequestInstance {
  name: string;
  maxConcurrency: number;
  concurrency: number;
}
interface RequestInstances {
  processing: boolean;
  items: RequestInstance[];
}
const requestInstances: RequestInstances = reactive({
  processing: false,
  items: [],
});

export function commonGetEmptyCaptcha(): Captcha {
  return <Captcha>{};
}

// commonGetCaptcha 获取图形验证码
export async function commonGetCaptcha(): Promise<Captcha> {
  const { data } = await request.get(COMMONS_CAPTCHA);
  return <Captcha>data;
}

export function commonGetSettings(): void {
  const data = settingStorage.getData();
  if (data.theme) {
    settings.theme = data.theme as string;
  } else if (
    window.matchMedia &&
    window.matchMedia("(prefers-color-scheme: light)").matches
  ) {
    settings.theme = "light";
  }
  settings.collapsed = data.collapsed as boolean;
}

export async function commonUpdateSettingTheme(theme: string): Promise<void> {
  await settingStorage.set("theme", theme);
  settings.theme = theme;
}

export async function commonUpdateSettingCollapsed(
  collapsed: boolean
): Promise<void> {
  await settingStorage.set("collapsed", collapsed);
  settings.collapsed = collapsed;
}

// commonListRouter 获取路由列表
export async function commonListRouter(): Promise<void> {
  if (routers.processing || routers.items.length !== 0) {
    return;
  }
  try {
    routers.processing = true;
    const { data } = await request.get<{
      routers: Router[];
    }>(COMMONS_ROUTERS);
    routers.items = data.routers || [];
  } finally {
    routers.processing = false;
  }
}

// 获取http request实例
export async function commonListRequestInstance(): Promise<void> {
  if (requestInstances.processing || requestInstances.items.length !== 0) {
    return;
  }
  try {
    requestInstances.processing = true;
    const { data } = await request.get<{
      statsList: RequestInstance[];
    }>(COMMONS_HTTP_STATS);
    requestInstances.items = data.statsList;
  } finally {
    requestInstances.processing = false;
  }
}

interface ReadonlyCommonState {
  settings: DeepReadonly<Settings>;
  routers: DeepReadonly<Routers>;
  requestInstances: DeepReadonly<RequestInstances>;
}

const state = {
  settings: readonly(settings),
  routers: readonly(routers),
  requestInstances: readonly(requestInstances),
};

export default function useCommonState(): ReadonlyCommonState {
  return state;
}
