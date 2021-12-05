import localforage from "localforage";

const store = localforage.createInstance({
  name: "tiny-site",
});

class LocalStorage {
  private key: string;
  private data: Record<string, unknown>;
  constructor(key: string) {
    this.data = {};
    this.key = key;
  }
  getData(): Record<string, unknown> {
    return Object.assign({}, this.data);
  }
  // load 加载数据
  async load(): Promise<Record<string, unknown>> {
    const data = await store.getItem(this.key);
    if (!data) {
      return {};
    }
    const str = data as string;
    this.data = JSON.parse(str || "{}");
    return this.getData();
  }
  // set 设置数据
  async set(key: string, value: string | number | boolean) {
    const data = this.data;
    data[key] = value;
    await store.setItem(this.key, JSON.stringify(data));
    return this.getData();
  }
}

export const settingStorage = new LocalStorage("settings");
