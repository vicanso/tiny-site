import { DeepReadonly, reactive, readonly } from "vue";

import request from "../helpers/request";

import { IMAGES, IMAGES_BUCKETS, IMAGES_BUCKETS_ID } from "../constants/url";
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

export interface Image {
  [key: string]: unknown;
  id: number;
  createdAt: string;
  updatedAt: string;
  bucket: string;
  name: string;
  type: string;
  size: number;
  width: number;
  height: number;
  tags: string[];
  creator: string;
  description: string;
}

const images: IList<Image> = reactive({
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

export async function bucketSearch(keyword: string): Promise<Bucket[]> {
  const { data } = await request.get<{
    buckets: Bucket[];
  }>(IMAGES_BUCKETS, {
    params: {
      keyword,
      limit: 10,
    },
  });
  return data.buckets;
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

export async function imageList(params: {
  bucket: string;
  tag: string;
  limit: number;
  offset: number;
}) {
  if (images.processing) {
    return;
  }
  try {
    images.processing = true;
    const { data } = await request.get<{
      images: Image[];
      count: number;
    }>(IMAGES, {
      params: Object.assign(
        {
          fields: [
            "id",
            "createdAt",
            "updatedAt",
            "bucket",
            "name",
            "type",
            "size",
            "width",
            "height",
            "tags",
            "creator",
            "description",
          ].join(","),
          order: "-updatedAt",
        },
        params
      ),
    });
    const count = data.count || 0;
    if (count >= 0) {
      images.count = count;
    }
    images.items = data.images || [];
  } finally {
    images.processing = false;
  }
}

interface ReadonlyImageState {
  buckets: DeepReadonly<IList<Bucket>>;
  images: DeepReadonly<IList<Image>>;
}

const state = {
  buckets: readonly(buckets),
  images: readonly(images),
};

export default function useImageState(): ReadonlyImageState {
  return state;
}
