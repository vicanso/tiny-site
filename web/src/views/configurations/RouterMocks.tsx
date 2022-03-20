import { defineComponent, onMounted } from "vue";
import ExConfigEditorList from "../../components/ExConfigEditorList";
import { ConfigCategory } from "../../states/configs";
import { FormItemTypes, FormItem } from "../../components/ExForm";
import useCommonState, { commonListRouter } from "../../states/common";
import { useMessage } from "naive-ui";
import { showError } from "../../helpers/util";
import ExLoading from "../../components/ExLoading";
import {
  getDefaultFormRules,
  newRequireRule,
} from "../../components/ExConfigEditor";

export default defineComponent({
  name: "RouterMockConfigs",
  setup() {
    const { routers } = useCommonState();
    const message = useMessage();
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
        type: FormItemTypes.Blank,
        name: "",
        key: "",
      },
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
        name: "状态码：",
        key: "data.status",
        type: FormItemTypes.InputNumber,
        placeholder: "请输入响应状态码",
      },
      {
        name: "响应数据类型：",
        type: FormItemTypes.Select,
        key: "data.contentType",
        placeholder: "请选择响应数据类型",
        options: [
          {
            label: "json",
            value: "application/json; charset=UTF-8",
          },
          {
            label: "plain",
            value: "text/plain; charset=UTF-8",
          },
          {
            label: "html",
            value: "text/html; charset=UTF-8",
          },
        ],
      },
      {
        name: "延时响应：",
        key: "data.delaySeconds",
        type: FormItemTypes.InputNumber,
        placeholder: "请输入延时响应时长(秒)",
      },
      {
        name: "完整URL：",
        key: "data.url",
        placeholder: "请输入完整的url(可选)",
      },
      {
        name: "响应数据：",
        key: "data.response",
        type: FormItemTypes.TextArea,
        span: 24,
        placeholder: "请输入响应数据",
      },
    ];
    const rules = getDefaultFormRules({
      "data.router": newRequireRule("路由不能为空"),
      "data.status": {
        required: true,
        message: "配置状态不能为空",
        trigger: "blur",
        validator(rule, value) {
          if (!value) {
            return false;
          }
          return true;
        },
      },
      "data.contentType": newRequireRule("响应类型不能为空"),
      "data.response": newRequireRule("响应数据不能为空"),
    });
    return (
      <ExConfigEditorList
        listTitle="路由Mock配置"
        editorTitle="添加/更新路由Mock配置"
        editorDescription="设置各路由的Mock响应"
        category={ConfigCategory.Router}
        extraFormItems={extraFormItems}
        rules={rules}
      />
    );
  },
});
