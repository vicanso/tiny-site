import { useMessage } from "naive-ui";
import { defineComponent, onMounted } from "vue";

import useCommonState, { commonListRouter } from "../../states/common";
import { showError } from "../../helpers/util";
import ExLoading from "../../components/ExLoading";
import { FormItemTypes, FormItem } from "../../components/ExForm";
import ExConfigEditorList from "../../components/ExConfigEditorList";
import { ConfigCategory } from "../../states/configs";
import { getDefaultFormRules, newRequireRule } from "../../components/ExConfigEditor";

export default defineComponent({
  name: "HTTPServerInterceptorConfigs",
  setup() {
    const message = useMessage();
    const { routers } = useCommonState();

    onMounted(async () => {
      try {
        await commonListRouter();
      } catch (err) {
        showError(message, err);
      }
    });

    return {
      routers,
    };
  },
  render() {
    const { routers } = this;
    if (routers.processing) {
      return <ExLoading />;
    }
    const extraFormItems: FormItem[] = [
      {
        name: "路由：",
        key: "data.router",
        type: FormItemTypes.Select,
        placeholder: "请选择路由",
        options: routers.items.map((item) => {
          const value = `${item.method} ${item.route}`;
          return {
            label: value,
            value,
          };
        }),
      },
      {
        name: "前置脚本：",
        key: "data.before",
        type: FormItemTypes.TextArea,
        span: 24,
        placeholder: "请输入请求处理前的相关处理脚本",
      },
      {
        name: "后置脚本：",
        key: "data.after",
        type: FormItemTypes.TextArea,
        span: 24,
        placeholder: "请输入请求处理后的相关处理脚本",
      },
    ];
    const rules = getDefaultFormRules({
      "data.router": newRequireRule("路由不能为空"),
    });
    return (
      <ExConfigEditorList
        listTitle="HTTP服务拦截配置"
        editorTitle="添加/更新HTTP服务拦截配置"
        editorDescription="设置HTTP服务各路由的拦截配置"
        category={ConfigCategory.HTTPServerInterceptor}
        extraFormItems={extraFormItems}
        rules={rules}
      />
    );
  },
});
