import { defineComponent } from "vue";
import ExConfigEditorList from "../../components/ExConfigEditorList";
import { ConfigCategory } from "../../states/configs";
import { getDefaultFormRules, newRequireRule } from "../../components/ExConfigEditor";

export default defineComponent({
  name: "BlockIPConfigs",
  render() {
    const extraFormItems = [
      {
        name: "IP地址：",
        key: "data",
        placeholder: "请输入IP地址或网段",
      },
    ];
    const rules = getDefaultFormRules({
      data: newRequireRule("IP地址或网段不能为空"),
    });
    return (
      <ExConfigEditorList
        listTitle="黑名单IP配置"
        editorTitle="添加/更新黑名单配置"
        editorDescription="用于拦截访问IP"
        category={ConfigCategory.BlockIP}
        extraFormItems={extraFormItems}
        rules={rules}
      />
    );
  },
});
