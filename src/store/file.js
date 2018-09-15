import request from "axios";
import ms from "ms";
import bytes from "bytes";

import { FILE_CATEGORIES, FILE_LIST } from "@/store/types";

import { FILES_CATEGORIES, FILES } from "@/urls";
import { debug, formatDate } from "@/helpers/util";

const fileListCategory = async ({ commit }) => {
  const res = await request.get(FILES_CATEGORIES);
  commit(FILE_CATEGORIES, res.data.categories);
};

const fileSave = async (tmp, { id, fileType, category, maxAge }) => {
  const res = await request.post(FILES, {
    file: id,
    category,
    fileType,
    maxAge
  });
  return res;
};

const fileList = async (
  { commit },
  { category, fields, order, skip, limit }
) => {
  const params = {
    category,
    fields,
    order,
    skip,
    limit
  };
  debug(params);
  const res = await request.get(FILES, {
    params
  });
  debug(res);
  commit(FILE_LIST, {
    data: res.data,
    params
  });
};

const state = {
  file: {
    categories: null,
    list: null,
    count: 0
  }
};

const actions = {
  fileListCategory,
  fileSave,
  fileList
};

const mutations = {
  // 文件分类信息
  [FILE_CATEGORIES](state, data) {
    state.file.categories = data;
  },
  [FILE_LIST](state, { data, params }) {
    if (!params.skip) {
      state.file.list = [];
    }
    data.files.forEach(function(item) {
      item.createdAt = formatDate(item.createdAt);
      item.maxAge = ms(ms(item.maxAge), { long: true });
      item.size = bytes(item.size);
      state.file.list.push(item);
    });
    if (data.count >= 0) {
      state.file.count = data.count;
    }
  }
};

export default {
  actions,
  state,
  mutations
};
