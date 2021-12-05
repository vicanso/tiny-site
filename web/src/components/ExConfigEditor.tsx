import {
  NCard,
  NPageHeader,
  NSpin,
  useMessage,
  FormRules,
  FormItemRule,
} from "naive-ui";
import { defineComponent, PropType, ref, Ref } from "vue";
import { get, isObject } from "lodash-es";
import { showError, showWarning } from "../helpers/util";
import {
  Config,
  configAdd,
  configFindByID,
  ConfigStatus,
  configUpdateByID,
} from "../states/configs";
import ExForm, { FormItem, FormItemTypes } from "./ExForm";
import ExLoading from "./ExLoading";

const statusKey = "status.value";

export enum FormItemKey {
  name = "name",
  category = "category",
  status = "status",
  startedAt = "startedAt",
  endedAt = "endedAt",
}

export function getDefaultFormItems(params: {
  category: string;
  name?: string;
}): FormItem[] {
  return [
    {
      name: "名称：",
      key: FormItemKey.name,
      disabled: params.name != null,
      defaultValue: params.name,
      placeholder: "请输入配置名称",
    },
    {
      name: "分类：",
      key: FormItemKey.category,
      disabled: true,
      defaultValue: params.category,
    },
    {
      name: "状态：",
      key: statusKey,
      type: FormItemTypes.Select,
      placeholder: "请选择配置状态",
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
      name: "开始时间：",
      key: FormItemKey.startedAt,
      type: FormItemTypes.DateTime,
      placeholder: "请选择配置生效开始时间",
    },
    {
      name: "结束时间：",
      key: FormItemKey.endedAt,
      type: FormItemTypes.DateTime,
      placeholder: "请选择配置生效结束时间",
    },
  ];
}

export function newRequireRule(message: string): FormItemRule {
  return {
    required: true,
    message: message,
    trigger: "blur",
  };
}

export function getDefaultFormRules(extra?: FormRules): FormRules {
  const defaultRules: FormRules = {
    [FormItemKey.name]: newRequireRule("配置名称不能为空"),
    [FormItemKey.category]: newRequireRule("配置分类不能为空"),
    [FormItemKey.status]: {
      value: {
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
    },
    [FormItemKey.startedAt]: newRequireRule("配置生效开始时间不能为空"),
    [FormItemKey.endedAt]: newRequireRule("配置生效结束时间不能为空"),
  };
  if (!extra) {
    return defaultRules;
  }
  return Object.assign(defaultRules, extra);
}

function convertDataToConfig(data: Record<string, unknown>): Config {
  // 转换配置数据
  const dataKeyPrefix = "data.";
  const configData: Record<string, unknown> = {};
  if (isObject(data.data)) {
    Object.assign(configData, data.data);
  }
  Object.keys(data).forEach((key) => {
    if (!key.startsWith(dataKeyPrefix)) {
      return;
    }
    configData[key.substring(dataKeyPrefix.length)] = data[key];
  });
  let configDataStr = data.data as string;
  if (Object.keys(configData).length !== 0) {
    configDataStr = JSON.stringify(configData);
  }
  return {
    name: data.name,
    status: {
      value: get(data, statusKey),
    },
    category: data.category,
    startedAt: data.startedAt,
    endedAt: data.endedAt,
    data: configDataStr,
    description: data.description,
  } as Config;
}

function diffConfig(
  newConfig: Config,
  current: Config
): Record<string, unknown> {
  const data: Record<string, unknown> = {};
  if (newConfig.name != current.name) {
    data.name = newConfig.name;
  }
  if (newConfig.status.value != current.status.value) {
    data.status = newConfig.status;
  }
  if (newConfig.category !== current.category) {
    data.category = newConfig.category;
  }
  if (newConfig.startedAt !== current.startedAt) {
    data.startedAt = newConfig.startedAt;
  }
  if (newConfig.endedAt !== current.endedAt) {
    data.endedAt = newConfig.endedAt;
  }
  if (newConfig.data !== current.data) {
    data.data = newConfig.data;
  }
  if (newConfig.description !== current.description) {
    data.description = newConfig.description;
  }
  return data;
}

function noop(): void {
  // 无操作
}

export default defineComponent({
  name: "ExConfigEditor",
  props: {
    id: {
      type: Number,
      default: 0,
    },
    title: {
      type: String,
      required: true,
    },
    description: {
      type: String,
      required: true,
    },
    formItems: {
      type: Array as PropType<FormItem[]>,
      required: true,
    },
    onSubmitDone: {
      type: Function as PropType<() => void>,
      default: noop,
    },
    onBack: {
      type: Function as PropType<() => void>,
      default: noop,
    },
    rules: {
      type: Object as PropType<FormRules>,
      default: null,
    },
  },
  setup(props) {
    const message = useMessage();
    const isUpdatedMode = props.id !== 0;
    const processing = ref(false);
    const currentConfig: Ref<Config> = ref({} as Config);
    // 提交数据
    const onSubmit = async (data: Record<string, unknown>) => {
      if (processing.value) {
        return;
      }
      const { name, category, startedAt, endedAt } = data;
      const status = get(data, statusKey);
      if (!name || !category || !status || !startedAt || !endedAt) {
        showWarning(
          message,
          "配置名称、分类、状态、开始时间、结束时间均不能为空"
        );
        return;
      }

      try {
        processing.value = true;
        const configData = convertDataToConfig(data);
        if (isUpdatedMode) {
          const updateData = diffConfig(configData, currentConfig.value);
          if (Object.keys(updateData).length === 0) {
            showWarning(message, "数据未修改无需要更新");
            return;
          }
          await configUpdateByID({
            id: props.id,
            data: updateData,
          });
          currentConfig.value = configData;
        } else {
          await configAdd(configData);
        }
        props.onSubmitDone();
      } catch (err) {
        showError(message, err);
      } finally {
        processing.value = false;
      }
    };
    const items: FormItem[] = [];
    // 由于会对form item的元素填写默认值，因此重新clone
    props.formItems.forEach((item) => {
      items.push(Object.assign({}, item));
    });
    items.push({
      name: "配置描述：",
      key: "description",
      type: FormItemTypes.TextArea,
      placeholder: "请输入配置描述",
      span: 24,
    });
    // 如果指定了id，则拉取配置
    const fetch = async () => {
      processing.value = true;
      try {
        const configData = await configFindByID(props.id);
        currentConfig.value = configData;
        const result = configData as Record<string, unknown>;
        items.forEach((item) => {
          const { key } = item;
          if (!key) {
            return;
          }
          const arr = key.split(".");
          if (arr.length === 2 && arr[0] === "data") {
            try {
              const data = result["data"] as string;
              item.defaultValue = JSON.parse(data)[arr[1]];
            } catch (err) {
              console.error(err);
            }
            return;
          }

          item.defaultValue = get(result, key);
        });
      } finally {
        processing.value = false;
      }
    };
    if (isUpdatedMode) {
      fetch();
    }

    return {
      currentConfig,
      processing,
      onSubmit,
      items,
    };
  },
  render() {
    const { title, description, id, onBack, rules } = this.$props;
    const { onSubmit, processing, items, currentConfig } = this;
    // 如果指定了id，则展示加载中
    if (processing && id && !currentConfig.id) {
      return <ExLoading />;
    }
    return (
      <NSpin show={processing}>
        <NCard>
          <NPageHeader
            title={title}
            onBack={onBack == noop ? undefined : onBack}
            subtitle={description}
          >
            <ExForm
              formItems={items}
              rules={rules}
              onSubmit={onSubmit}
              submitText={id !== 0 ? "更新" : "添加"}
            />
          </NPageHeader>
        </NCard>
      </NSpin>
    );
  },
});
