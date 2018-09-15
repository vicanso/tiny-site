import Vue from "vue";
import App from "@/App.vue";
import router from "@/router";
import store from "@/store";
import ElementUI from "element-ui";
import "@/styles/index.scss";
import "@/request-interceptors";
import { timeout } from "@/config";
import { getErrorMessage } from "@/helpers/util";

Vue.use(ElementUI);

// 注入 router 和 store
Vue.$router = router;
Vue.$store = store;

Vue.prototype.xLoading = (options = {}) => {
  let loadingInstance = ElementUI.Loading.service(
    Object.assign(
      {
        fullscreen: true,
        text: "正在处理中，请稍候..."
      },
      options
    )
  );
  let resolved = false;
  const resolve = () => {
    if (resolved) {
      return;
    }
    resolved = true;
    loadingInstance.close();
  };
  setTimeout(resolve, options.timeout || timeout);
  return resolve;
};

Vue.prototype.xError = function xError(err) {
  const message = getErrorMessage(err);
  this.$message.error(message);
};

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
