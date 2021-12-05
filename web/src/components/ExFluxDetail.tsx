import { InfoCircle } from "@vicons/fa";
import { NIcon, NPopover } from "naive-ui";
import { defineComponent, PropType } from "vue";
import ExFluxDetailList from "./ExFluxDetailList";

export default defineComponent({
  name: "ExFluxDetail",
  props: {
    measurement: {
      type: String,
      required: true,
    },
    data: {
      type: Object as PropType<Record<string, unknown>>,
      required: true,
    },
    tagKeys: {
      type: Array as PropType<string[]>,
      required: true,
    },
  },
  render() {
    const { data, measurement, tagKeys } = this.$props;
    const slots = {
      trigger: () => (
        <NIcon>
          <InfoCircle />
        </NIcon>
      ),
    };
    const tags: Record<string, string> = {};
    Object.keys(data).forEach((key) => {
      if (!tagKeys.includes(key)) {
        return;
      }
      const v = data[key];
      if (!v) {
        return;
      }
      tags[key] = v as string;
    });
    return (
      <NPopover
        trigger="hover"
        placement="top-end"
        delay={500}
        duration={1000}
        v-slots={slots}
      >
        <ExFluxDetailList
          measurement={measurement}
          time={data._time as string}
          tags={tags}
        />
      </NPopover>
    );
  },
});
