import {
  FormRules,
  FormInst,
  NButton,
  NDatePicker,
  NForm,
  NFormItem,
  NGrid,
  NGridItem,
  NInput,
  NInputNumber,
  NSelect,
  NInputGroup,
  NInputGroupLabel,
  useMessage,
} from "naive-ui";
import { set } from "lodash-es";
import { Value } from "naive-ui/lib/select/src/interface";
import { Component, defineComponent, PropType, ref } from "vue";
import { showError, durationToSeconds } from "../helpers/util";

export interface FormItem {
  name: string;
  key: string;
  type?: string;
  placeholder?: string;
  span?: number;
  defaultValue?: unknown;
  disabled?: boolean;
  // TODO 确认是否有其它方式表示
  // eslint-disable-next-line
  options?: any[];
}

export enum FormItemTypes {
  Select = "select",
  MultiSelect = "multiselect",
  DateTime = "datetime",
  DateRange = "datrange",
  InputNumber = "inputnumber",
  InputDuration = "inputDuration",
  TextArea = "textarea",
  Blank = "blank",
}

export default defineComponent({
  name: "ExForm",
  props: {
    formItems: {
      type: Array as PropType<FormItem[]>,
      required: true,
    },
    onSubmit: {
      type: Function as PropType<
        (data: Record<string, unknown>) => Promise<void>
      >,
      required: true,
    },
    submitText: {
      type: String,
      default: "提交",
    },
    rules: {
      type: Object as PropType<FormRules>,
      default: null,
    },
  },
  setup(props) {
    const formRef = ref({} as FormInst);
    const params = ref({} as Record<string, unknown>);
    props.formItems.forEach((item) => {
      if (item.defaultValue) {
        set(params.value, item.key, item.defaultValue);
      }
    });
    const formValidate = (): Promise<void> => {
      return new Promise((resolve, reject) => {
        formRef.value.validate((errors) => {
          if (errors) {
            const msgList = errors.map((arr) => {
              return arr.map((item) => item.message).join(",");
            });
            reject(new Error(msgList.join(";")));
            return;
          }
          resolve();
        });
      });
    };
    const message = useMessage();
    return {
      handleSubmit: async (data: Record<string, unknown>) => {
        try {
          await formValidate();
        } catch (err) {
          // 如果出错，则不再提交
          showError(message, err);
          return;
        }
        return props.onSubmit(data);
      },
      params,
      formRef,
    };
  },
  render() {
    const { submitText, rules } = this.$props;
    const { params, handleSubmit } = this;
    const size = "large";
    const createSelect = (item: FormItem, multiple: boolean) => {
      return (
        <NSelect
          filterable
          multiple={multiple}
          defaultValue={item.defaultValue as Value}
          options={item.options || []}
          placeholder={item.placeholder}
          onUpdateValue={(value) => {
            set(params, item.key, value);
          }}
        />
      );
    };
    const formItems = this.$props.formItems as FormItem[];
    const arr = formItems.map((item) => {
      let component: Component;
      switch (item.type) {
        case FormItemTypes.Blank:
          component = <div />;
          break;
        case FormItemTypes.MultiSelect:
          component = createSelect(item, true);
          break;
        case FormItemTypes.Select:
          component = createSelect(item, false);
          break;
        case FormItemTypes.DateTime:
          {
            let defaultValue = null;
            if (item.defaultValue) {
              defaultValue = new Date(item.defaultValue as string).getTime();
            }
            component = (
              <NDatePicker
                type="datetime"
                class="widthFull"
                placeholder={item.placeholder}
                defaultValue={defaultValue}
                clearable
                onUpdateValue={(value) => {
                  if (!value) {
                    params[item.key] = "";
                  } else {
                    params[item.key] = new Date(value).toISOString();
                  }
                }}
              />
            );
          }
          break;
        case FormItemTypes.InputNumber:
          {
            component = (
              <NInputNumber
                class="widthFull"
                disabled={item.disabled || false}
                placeholder={item.placeholder}
                defaultValue={(item.defaultValue || null) as number}
                onUpdate:value={(value) => {
                  params[item.key] = value;
                }}
              />
            );
          }
          break;
        case FormItemTypes.InputDuration:
          {
            component = (
              <NInputGroup>
                <NInputNumber
                  class="widthFull"
                  placeholder={item.placeholder}
                  defaultValue={
                    (durationToSeconds(item.defaultValue as string) ||
                      null) as number
                  }
                  onUpdate:value={(value) => {
                    params[item.key] = `${value || 0}s`;
                  }}
                />
                <NInputGroupLabel size={size}>秒</NInputGroupLabel>
              </NInputGroup>
            );
          }
          break;
        case FormItemTypes.TextArea:
          {
            component = (
              <NInput
                type="textarea"
                autosize={{
                  minRows: 3,
                  maxRows: 5,
                }}
                disabled={item.disabled || false}
                placeholder={item.placeholder}
                defaultValue={(item.defaultValue || "") as string}
                onUpdateValue={(value) => {
                  params[item.key] = value;
                }}
                clearable
              />
            );
          }
          break;
        default:
          component = (
            <NInput
              disabled={item.disabled || false}
              placeholder={item.placeholder}
              defaultValue={(item.defaultValue || "") as string}
              onUpdateValue={(value) => {
                set(params, item.key, value);
              }}
              clearable
            />
          );
          break;
      }
      return (
        <NGridItem span={item.span || 8}>
          <NFormItem label={item.name} path={item.key}>
            {component}
          </NFormItem>
        </NGridItem>
      );
    });
    arr.push(
      <NGridItem span={24}>
        <NFormItem>
          <NButton class="widthFull" onClick={() => handleSubmit(params)}>
            {submitText}
          </NButton>
        </NFormItem>
      </NGridItem>
    );
    return (
      <NForm
        labelPlacement="left"
        rules={rules}
        model={params}
        ref="formRef"
        size={size}
      >
        <NGrid xGap={24}>{arr}</NGrid>
      </NForm>
    );
  },
});
