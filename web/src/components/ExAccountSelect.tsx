import { NSelect, SelectOption, useMessage } from "naive-ui";
import { defineComponent, onBeforeMount, PropType, ref } from "vue";
import { Value } from "naive-ui/lib/select/src/interface";

import { userListAccount } from "../states/user";
import { showError } from "../helpers/util";

export default defineComponent({
  name: "ExAccountSelect",
  props: {
    placeholder: {
      type: String,
      default: () => "",
    },
    onUpdateValue: {
      type: Function as PropType<(value: string[]) => void>,
      required: true,
    },
    defaultValue: {
      type: Array as PropType<Value>,
      default: () => null,
    },
  },
  setup() {
    const loading = ref(false);
    const options = ref<SelectOption[]>([]);
    const message = useMessage();
    const onSearch = async (query: string) => {
      loading.value = true;
      try {
        const accounts = await userListAccount(query);
        options.value = accounts.map((item) => {
          return {
            label: item,
            value: item,
          };
        });
      } catch (err) {
        showError(message, err);
      } finally {
        loading.value = false;
      }
    };
    onBeforeMount(() => onSearch(""));
    return {
      loading,
      options,
      onSearch,
    };
  },
  render() {
    const { loading, options, onSearch } = this;
    return (
      <NSelect
        multiple
        filterable
        defaultValue={this.$props.defaultValue}
        placeholder={this.$props.placeholder}
        options={options}
        loading={loading}
        clearable
        remote
        onSearch={onSearch}
        onUpdateValue={(value) => {
          this.$props.onUpdateValue(value);
        }}
      ></NSelect>
    );
  },
});
