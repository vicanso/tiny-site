import { NSelect, SelectOption, useMessage } from "naive-ui";
import { defineComponent, onBeforeMount, PropType, ref } from "vue";
import { Value } from "naive-ui/lib/select/src/interface";

import { bucketSearch } from "../states/image";
import { showError } from "../helpers/util";

export default defineComponent({
  name: "ExBucketSelect",
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
      type: [String, Array] as PropType<Value>,
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
        const buckets = await bucketSearch(query);
        options.value = buckets.map((item) => {
          return {
            label: item.name,
            value: item.name,
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
        size="large"
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
      />
    );
  },
});
