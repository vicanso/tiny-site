import { css } from "@linaria/core";
import { useMessage } from "naive-ui";
import { defineComponent, onMounted, PropType, ref } from "vue";
import { showError } from "../helpers/util";
import { fluxFindOne } from "../states/flux";
const infoListCalss = css`
  margin: 0;
  padding: 0 20px;
  max-width: 400px;
  word-wrap: break-word;
  word-break: break-all;
  white-space: normal;
  list-style-position: insied;
`;

export default defineComponent({
  name: "FluxDetailList",
  props: {
    measurement: {
      type: String,
      required: true,
    },
    time: {
      type: String,
      required: true,
    },
    tags: {
      type: Object as PropType<Record<string, string>>,
      required: true,
    },
  },
  setup(props) {
    const message = useMessage();
    const processing = ref(true);
    const values = ref([] as Record<string, unknown>[]);
    onMounted(async () => {
      try {
        const data = await fluxFindOne({
          measurement: props.measurement,
          time: props.time,
          tags: props.tags,
        });
        const ignoreKeys = [
          "_measurement",
          "_start",
          "_stop",
          "_time",
          "result",
          "table",
        ];
        Object.keys(data).forEach((key) => {
          if (ignoreKeys.includes(key)) {
            return;
          }
          values.value.push({
            name: key,
            value: data[key],
          });
        });
      } catch (err) {
        showError(message, err);
      } finally {
        processing.value = false;
      }
    });
    return {
      processing,
      values,
    };
  },
  render() {
    const { processing, values } = this;
    if (processing) {
      return <span>正在加载中...</span>;
    }
    if (values.length === 0) {
      return <span>很抱歉，无符合记录</span>;
    }
    const arr = values.map((item) => {
      return (
        <li>
          <span class="mright5">{item.name}:</span> {String(item.value)}
        </li>
      );
    });
    return <ul class={infoListCalss}>{arr}</ul>;
  },
});
