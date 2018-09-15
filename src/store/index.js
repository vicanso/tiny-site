import Vue from "vue";
import Vuex from "vuex";

import userStore from "@/store/user";
import fileStore from "@/store/file";

Vue.use(Vuex);

const defaultOptions = {
  state: {},
  mutations: {},
  actions: {}
};
[userStore, fileStore].forEach(item => {
  Object.assign(defaultOptions.state, item.state);
  Object.assign(defaultOptions.mutations, item.mutations);
  Object.assign(defaultOptions.actions, item.actions);
});

export default new Vuex.Store(defaultOptions);
