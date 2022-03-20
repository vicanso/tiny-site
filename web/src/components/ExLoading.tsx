import { NSpin } from "naive-ui";
import { defineComponent } from "vue";
import { css } from "@linaria/core";

const spinClass = css`
  float: left;
  margin-right: 10px;
`;
const loadingClass = css`
  margin: auto;
  width: 200px;
`;

export default defineComponent({
  name: "ExLoading",
  props: {
    style: {
      type: Object,
      default: () => {
        return {
          marginTop: "60px",
        };
      },
    },
  },
  render() {
    const { style } = this.$props;
    return (
      <div style={style} class={loadingClass}>
        <NSpin size="small" class={spinClass} />
        正在加载中，请稍候...
      </div>
    );
  },
});
