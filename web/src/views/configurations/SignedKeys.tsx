import { defineComponent } from "vue";
import ExConfigEditorList from "../../components/ExConfigEditorList";
import { ConfigCategory } from "../../states/configs";
import { getDefaultFormRules, newRequireRule } from "../../components/ExConfigEditor";

export default defineComponent({
  name: "SignedKeyConfigs",
  render() {
    const extraFormItems = [
      {
        name: "密钥：",
        key: "data",
        placeholder: "请输入签名使用的密钥，多个密钥以,分隔",
      },
    ];
    const rules = getDefaultFormRules({
      data: newRequireRule("密钥数据不能为空"),
    });
    return (
      <ExConfigEditorList
        listTitle="密钥配置"
        editorTitle="添加/更新密钥配置"
        editorDescription="配置用于生成session加密的密钥"
        category={ConfigCategory.SignedKey}
        extraFormItems={extraFormItems}
        rules={rules}
      />
    );
  },
});
