import { DeepReadonly, reactive, readonly } from "vue";

import request from "../helpers/request";

import { STORAGES, STORAGES_ID } from "../constants/url";
import { IList, IStatus } from "./interface";

export enum StorageCategory {
  HTTP = "http",
  Minio = "minio",
  OSS = "oss",
  Gridfs = "gridfs",
}

export interface Storage {
  [key: string]: unknown;
  id: number;
  createdAt: string;
  updatedAt: string;
  status: IStatus;
  name: string;
  category: string;
  uri: string;
  description?: string;
}

const storages: IList<Storage> = reactive({
  processing: false,
  items: [],
  count: -1,
});

export async function storageList(): Promise<void> {
  if (storages.processing) {
    return;
  }
  try {
    storages.processing = true;
    const { data } = await request.get<{
      storages: Storage[];
    }>(STORAGES);
    storages.items = data.storages;
    storages.count = data.storages.length;
  } finally {
    storages.processing = false;
  }
}

export async function storageAdd(
  params: Record<string, unknown>
): Promise<Storage> {
  const { data } = await request.post<Storage>(STORAGES, params);
  return data;
}

export async function storageFindByID(id: number): Promise<Storage> {
  const url = STORAGES_ID.replace(":id", `${id}`);
  const { data } = await request.get<Storage>(url);
  return data;
}

export async function storageUpdateByID(
  id: number,
  updatedData: Record<string, unknown>
): Promise<Storage> {
  const url = STORAGES_ID.replace(":id", `${id}`);
  const { data } = await request.patch<Storage>(url, updatedData);
  return data;
}

interface ReadonlyStorageState {
  storages: DeepReadonly<IList<Storage>>;
}

const state = {
  storages: readonly(storages),
};

export default function useStorageState(): ReadonlyStorageState {
  return state;
}
