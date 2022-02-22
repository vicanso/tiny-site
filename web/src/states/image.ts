import { DeepReadonly, reactive, readonly } from "vue";

import request from "../helpers/request";

import { IMAGES_BUCKETS, IMAGES_BUCKETS_ID } from "../constants/url";
import { IList } from "./interface";

export interface Bucket {
  [key: string]: unknown;
  id: number;
  createdAt: string;
  updatedAt: string;
  name: string;
  creator: string;
  owners: string[];
  description: string;
}

const buckets: IList<Bucket> = reactive({
  processing: false,
  items: [],
  count: -1,
});

export async function bucketList(params: { limit: number; offset: number }) {
  if (buckets.processing) {
    return;
  }
  try {
    buckets.processing = true;
    const { data } = await request.get<{
      count: number;
      buckets: Bucket[];
    }>(IMAGES_BUCKETS, {
      params,
    });
    const count = data.count || 0;
    if (count >= 0) {
      buckets.count = count;
    }
    buckets.items = data.buckets || [];
  } finally {
    buckets.processing = false;
  }
}

export async function bucketAdd(params: {
  name: string;
  owners: string[];
  description: string;
}): Promise<Bucket> {
  const { data } = await request.post<Bucket>(IMAGES_BUCKETS, params);
  return data;
}

export async function bucketUpdate(
  id: number,
  params: Record<string, unknown>
): Promise<void> {
  const url = IMAGES_BUCKETS_ID.replace(":id", id.toString());
  await request.patch(url, params);
}

interface ReadonlyImageState {
  buckets: DeepReadonly<IList<Bucket>>;
}

const state = {
  buckets: readonly(buckets),
};

export default function useImageState(): ReadonlyImageState {
  return state;
}
