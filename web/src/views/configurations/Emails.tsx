import { defineComponent } from "vue";
import ExConfigEditorList from "../../components/ExConfigEditorList";
import { ConfigCategory } from "../../states/configs";
import { FormItem, FormItemTypes } from "../../components/ExForm";
import { getDefaultFormRules, newRequireRule } from "../../components/ExConfigEditor";

export default defineComponent({
  name: "EmailConfigs",
  render() {
    const extraFormItems: FormItem[] = [
      {
        name: "邮箱地址：",
        key: "data",
        type: FormItemTypes.TextArea,
        span: 24,
        placeholder: "请输入邮箱地址，多个地址以,分隔",
      },
    ];
    const rules = getDefaultFormRules({
      "data": newRequireRule("邮箱地址不能为空"),
    });
    return (
      <ExConfigEditorList
        listTitle="邮箱地址配置"
        editorTitle="添加/更新邮箱地址配置"
        editorDescription="配置各类邮件接收列表"
        category={ConfigCategory.Email}
        extraFormItems={extraFormItems}
        rules={rules}
      />
    );
  },
});
