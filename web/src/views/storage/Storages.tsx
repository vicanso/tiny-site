import { defineComponent, ref } from "vue";
import { FormRules, NButton, NCard, NPageHeader, useMessage } from "naive-ui";
import { css } from "@linaria/core";
import { TableColumn } from "naive-ui/lib/data-table/src/interface";

import { ConfigStatus } from "../../states/configs";
import useStorageState, {
  storageList,
  storageFindByID,
  storageAdd,
  storageUpdateByID,
} from "../../states/storage";
import ExForm, { FormItem, FormItemTypes } from "../../components/ExForm";
import ExTable, {
  newOPColumn,
  newLevelValueColumn,
} from "../../components/ExTable";
import { padding } from "../../constants/style";
import { Mode } from "../../states/common";
import ExLoading from "../../components/ExLoading";
import { IStatus } from "../../states/interface";
import { diff, showError, showWarning } from "../../helpers/util";
import { newRequireRule } from "../../components/ExConfigEditor";

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
      title: "类型",
      key: "category",
    },
    newLevelValueColumn({
      title: "状态",
      key: "status.desc",
    }),
    {
      title: "连接串",
      key: "uri",
    },
    {
      title: "描述",
      key: "description",
    },
  ];
}

export function getFormItems(params: {
  name?: string;
  category?: string;
  uri?: string;
  description?: string;
  status?: IStatus;
}): FormItem[] {
  return [
    {
      name: "名称：",
      key: "name",
      defaultValue: params.name,
      placeholder: "请输入存储服务名称",
    },
    {
      name: "类型：",
      key: "category",
      placeholder: "请选择存储类型",
      defaultValue: params.category,
      options: ["http", "minio", "oss", "gridfs"].map((item) => {
        return {
          label: item,
          value: item,
        };
      }),
      type: FormItemTypes.Select,
    },
    {
      name: "状态：",
      key: "status.value",
      type: FormItemTypes.Select,
      placeholder: "请选择配置状态",
      defaultValue: params.status?.value,
      options: [
        {
          label: "启用",
          value: ConfigStatus.Enabled,
        },
        {
          label: "禁用",
          value: ConfigStatus.Disabled,
        },
      ],
    },
    {
      name: "连接串：",
      key: "uri",
      span: 24,
      placeholder: "请输入该存储服务对应的连接串",
      defaultValue: params.uri,
    },
    {
      name: "描述：",
      key: "description",
      defaultValue: params.description,
      span: 24,
      type: FormItemTypes.TextArea,
      placeholder: "请输入该存储服务的描述",
    },
  ];
}

function diffStorage(
  newStorage: Record<string, unknown>,
  original: Record<string, unknown>
): Record<string, unknown> {
  const { data, modifiedCount } = diff(newStorage, original);
  if (modifiedCount === 0) {
    return data;
  }
  delete data.status;
  const newStatus = newStorage.status as Record<string, unknown>;
  const status = original.status as Record<string, unknown>;
  if (newStatus?.value !== status?.value) {
    data.status = newStatus.value;
  }

  return data;
}

export default defineComponent({
  name: "StorageList",
  setup() {
    const message = useMessage();
    const { storages } = useStorageState();
    const updatedID = ref(0);
    const currentStorage = ref({});
    const mode = ref(Mode.List);
    const toggle = (value: Mode) => {
      if (value === Mode.List) {
        currentStorage.value = {};
        updatedID.value = 0;
      }
      mode.value = value;
    };
    const processing = ref(false);
    const onUpdate = async (id: number) => {
      updatedID.value = id;
      mode.value = Mode.Update;
      try {
        processing.value = true;
        const storage = await storageFindByID(id);
        currentStorage.value = storage;
      } finally {
        processing.value = false;
      }
    };
    const update = async (data: Record<string, unknown>) => {
      const updatedData = diffStorage(data, currentStorage.value);
      if (Object.keys(updatedData).length === 0) {
        showWarning(message, "数据未修改无需要更新");
      }
      try {
        await storageUpdateByID(updatedID.value, updatedData);
        toggle(Mode.List);
      } catch (err) {
        showError(message, err);
      }
    };
    const onSubmit = async (data: Record<string, unknown>) => {
      if (updatedID.value === 0) {
        try {
          await storageAdd(data);
          toggle(Mode.List);
        } catch (err) {
          showError(message, err);
        }
        return;
      }
      update(data);
    };
    return {
      updatedID,
      storages,
      mode,
      toggle,
      onUpdate,
      processing,
      currentStorage,
      onSubmit,
    };
  },
  render() {
    const {
      storages,
      mode,
      onUpdate,
      toggle,
      updatedID,
      processing,
      currentStorage,
      onSubmit,
    } = this;
    if (mode == Mode.List) {
      const columns = getColumns();
      columns.push(
        newOPColumn((row) => {
          onUpdate(row.id as number);
        })
      );
      return (
        <NCard title={"图片存储服务"}>
          <ExTable
            columns={columns}
            data={storages}
            fetch={storageList}
            hidePagination={true}
          />
          <NButton
            size="large"
            class={addButtonClass}
            onClick={() => {
              toggle(Mode.Add);
            }}
          >
            增加服务
          </NButton>
        </NCard>
      );
    }
    if (processing) {
      return <ExLoading />;
    }
    const rules: FormRules = {
      name: newRequireRule("名称不能为空"),
      category: newRequireRule("类型不能为空"),
      status: newRequireRule("状态不能为空"),
      uri: newRequireRule("连接串不能为空"),
      description: newRequireRule("描述不能为空"),
    };
    return (
      <NCard>
        <NPageHeader
          title={"更新/添加图片存储服务"}
          onBack={() => {
            toggle(Mode.List);
          }}
        >
          <ExForm
            formItems={getFormItems(currentStorage)}
            submitText={updatedID !== 0 ? "更新" : "添加"}
            onSubmit={onSubmit}
            rules={rules}
          />
        </NPageHeader>
      </NCard>
    );
  },
});
