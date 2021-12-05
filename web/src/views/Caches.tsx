import { defineComponent, ref } from "vue";
import { NCard, NGrid, NGridItem, NInput, NButton, useMessage } from "naive-ui";
import { adminCleanCacheByKey, adminFindCacheByKey } from "../states/admin";
import { showError, showWarning, toast } from "../helpers/util";

export default defineComponent({
  name: "CachesView",
  setup() {
    const message = useMessage();
    const key = ref("");
    const processing = ref(false);
    const cacheData = ref("");

    const fetch = async () => {
      if (!key.value) {
        showWarning(message, "请输入要查询的key");
        return;
      }
      if (processing.value) {
        return;
      }
      processing.value = true;
      try {
        cacheData.value = "";
        const result = await adminFindCacheByKey(key.value);
        try {
          const json = JSON.parse(result.data);
          cacheData.value = JSON.stringify(json, null, 2);
        } catch (err) {
          cacheData.value = result.data;
        }
      } catch (err) {
        showError(message, err);
      } finally {
        processing.value = false;
      }
    };

    const del = async () => {
      if (!key.value) {
        showWarning(message, "请输入要删除的key");
        return;
      }
      if (processing.value) {
        return;
      }
      try {
        cacheData.value = "";
        await adminCleanCacheByKey(key.value);
        toast(message, "已成功清除数据");
      } catch (err) {
        showError(message, err);
      } finally {
        processing.value = false;
      }
    };

    return {
      processing,
      key,
      fetch,
      del,
      cacheData,
    };
  },
  render() {
    const size = "large";
    const { fetch, cacheData, del } = this;
    return (
      <NCard title="缓存查询与清除">
        <p>session的缓存格式 ss:sessionID</p>
        <NGrid xGap={24}>
          <NGridItem span={12}>
            <NInput
              placeholder="请输入缓存的key"
              size={size}
              clearable
              onUpdateValue={(value) => {
                this.key = value;
              }}
            />
          </NGridItem>
          <NGridItem span={6}>
            <NButton class="widthFull" size={size} onClick={() => fetch()}>
              查询
            </NButton>
          </NGridItem>
          <NGridItem span={6}>
            <NButton class="widthFull" size={size} onClick={() => del()}>
              清除
            </NButton>
          </NGridItem>
          {cacheData && (
            <NGridItem span={24}>
              <pre>{cacheData}</pre>
            </NGridItem>
          )}
        </NGrid>
      </NCard>
    );
  },
});
