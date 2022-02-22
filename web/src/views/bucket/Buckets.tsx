import { defineComponent, ref } from "vue";
import {
  FormRules,
  NButton,
  NCard,
  NPageHeader,
  NSpin,
  useMessage,
} from "naive-ui";
import { TableColumn } from "naive-ui/lib/data-table/src/interface";
import { css } from "@linaria/core";

import useImageState, {
  bucketAdd,
  bucketList,
  bucketUpdate,
} from "../../states/image";
import { Mode } from "../../states/common";
import ExTable, { newListColumn, newOPColumn } from "../../components/ExTable";
import { padding } from "../../constants/style";
import {
  newRequireRule,
  newArrayRequireRule,
} from "../../components/ExConfigEditor";
import ExForm, { FormItem, FormItemTypes } from "../../components/ExForm";
import { showError, diff, showWarning } from "../../helpers/util";

const addButtonClass = css`
  width: 100%;
  margin-top: ${2 * padding}px;
`;

function getColumns(): TableColumn[] {
  return [
    {
      title: "名称",
      key: "name",
    },
    {
      title: "创建者",
      key: "creator",
    },
    newListColumn({
      key: "owners",
      title: "拥有者",
    }),
    {
      title: "描述",
      key: "description",
    },
  ];
}

export function getFormItems(params: {
  name?: string;
  owners?: string[];
  description?: string;
}): FormItem[] {
  return [
    {
      name: "名称：",
      key: "name",
      defaultValue: params.name,
      span: 12,
      placeholder: "请输入图片bucket名称",
      disabled: !!params.name,
    },
    {
      name: "拥有者：",
      key: "owners",
      span: 12,
      defaultValue: params.owners,
      type: FormItemTypes.AccountSelect,
      placeholder: "请选择拥有者",
    },
    {
      name: "描述：",
      key: "description",
      defaultValue: params.description,
      span: 24,
      type: FormItemTypes.TextArea,
      placeholder: "请输入该存储bucket的描述",
    },
  ];
}

export default defineComponent({
  name: "BucketList",
  setup() {
    const message = useMessage();
    const { buckets } = useImageState();
    const mode = ref(Mode.List);
    const updatedID = ref(0);
    const currentBucket = ref({});
    const processing = ref(false);

    const toggle = (value: Mode) => {
      if (value === Mode.List) {
        currentBucket.value = {};
        updatedID.value = 0;
      }
      mode.value = value;
    };

    const onUpdate = async (id: number) => {
      const result = buckets.items.find((item) => item.id === id);
      if (result) {
        updatedID.value = id;
        currentBucket.value = result;
        toggle(Mode.Update);
      }
    };

    const onSubmit = async (data: Record<string, unknown>) => {
      if (processing.value) {
        return;
      }
      processing.value = true;
      try {
        if (updatedID.value != 0) {
          const diffInfo = diff(data, currentBucket.value);
          if (diffInfo.modifiedCount === 0) {
            showWarning(message, "请先修改数据再更新");
            return;
          }
          await bucketUpdate(updatedID.value, diffInfo.data);
        } else {
          await bucketAdd({
            name: data.name as string,
            owners: data.owners as string[],
            description: data.description as string,
          });
        }
        toggle(Mode.List);
      } catch (err) {
        showError(message, err);
      } finally {
        processing.value = false;
      }
    };

    return {
      buckets,
      toggle,
      mode,
      processing,
      currentBucket,
      updatedID,
      onSubmit,
      onUpdate,
    };
  },
  render() {
    const {
      buckets,
      toggle,
      mode,
      processing,
      currentBucket,
      updatedID,
      onSubmit,
      onUpdate,
    } = this;
    if (mode == Mode.List) {
      const columns = getColumns();
      columns.push(
        newOPColumn((row) => {
          onUpdate(row.id as number);
        })
      );
      return (
        <NCard title={"图片Bucket"}>
          <ExTable columns={columns} data={buckets} fetch={bucketList} />
          <NButton
            size="large"
            class={addButtonClass}
            onClick={() => {
              toggle(Mode.Add);
            }}
          >
            增加图片Bucket
          </NButton>
        </NCard>
      );
    }

    const rules: FormRules = {
      name: newRequireRule("名称不能为空"),
      owners: newArrayRequireRule("拥有者不能为空"),
      description: newRequireRule("描述不能为空"),
    };

    return (
      <NSpin show={processing}>
        <NCard>
          <NPageHeader
            title={"更新/添加图片bucket"}
            onBack={() => {
              toggle(Mode.List);
            }}
          >
            <ExForm
              formItems={getFormItems(currentBucket)}
              submitText={updatedID !== 0 ? "更新" : "添加"}
              onSubmit={onSubmit}
              rules={rules}
            />
          </NPageHeader>
        </NCard>
      </NSpin>
    );
  },
});
