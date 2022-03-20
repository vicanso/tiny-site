// 配置的列表展示与更新

import { NButton, NCard, FormRules } from "naive-ui";
import { css } from "@linaria/core";
import { defineComponent, PropType, ref } from "vue";
import { padding } from "../constants/style";
import ExConfigEditor, { getDefaultFormItems } from "./ExConfigEditor";
import ExConfigTable from "./ExConfigTable";
import { FormItem } from "./ExForm";
import { Mode } from "../states/common";

const addButtonClass = css`
  width: 100%;
  margin-top: ${2 * padding}px;
`;

export default defineComponent({
  name: "ExConfigEditorList",
  props: {
    listTitle: {
      type: String,
      required: true,
    },
    editorTitle: {
      type: String,
      required: true,
    },
    editorDescription: {
      type: String,
      required: true,
    },
    category: {
      type: String,
      required: true,
    },
    extraFormItems: {
      type: Array as PropType<FormItem[]>,
      default: () => [],
    },
    rules: {
      type: Object as PropType<FormRules>,
      default: null,
    },
  },
  setup() {
    const mode = ref(Mode.List);
    const updatedID = ref(0);
    const toggle = (value: Mode) => {
      mode.value = value;
    };
    return {
      updatedID,
      toggle,
      mode,
    };
  },
  render() {
    const {
      listTitle,
      editorTitle,
      category,
      editorDescription,
      extraFormItems,
      rules,
    } = this.$props;
    const { mode, toggle, updatedID } = this;
    if (mode === Mode.List) {
      return (
        <NCard title={listTitle}>
          <ExConfigTable
            category={category}
            onUpdate={(id: number) => {
              this.updatedID = id;
              toggle(Mode.Update);
            }}
          />
          <NButton
            size="large"
            class={addButtonClass}
            onClick={() => {
              this.updatedID = 0;
              toggle(Mode.Add);
            }}
          >
            增加配置
          </NButton>
        </NCard>
      );
    }

    const formItems = getDefaultFormItems({
      category,
    });
    extraFormItems.forEach((item) => {
      const data = Object.assign({}, item);
      formItems.push(data as FormItem);
    });

    return (
      <ExConfigEditor
        title={editorTitle}
        description={editorDescription}
        id={updatedID}
        formItems={formItems}
        onSubmitDone={() => {
          toggle(Mode.List);
        }}
        onBack={() => {
          toggle(Mode.List);
        }}
        rules={rules}
      ></ExConfigEditor>
    );
  },
});
