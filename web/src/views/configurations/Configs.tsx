import { defineComponent, ref } from "vue";
import { NTabPane, NTabs, useMessage } from "naive-ui";
import ExConfigTable from "../../components/ExConfigTable";
import { configGetCurrentValid } from "../../states/configs";
import { showError } from "../../helpers/util";
import ExLoading from "../../components/ExLoading";

const validTab = "当前生效配置";

export default defineComponent({
  name: "ConfigList",
  setup() {
    const message = useMessage();
    const tab = ref("所有配置");
    const currentValid = ref("");
    const fetchingCurrentValid = ref(false);
    const fetchValid = async () => {
      if (fetchingCurrentValid.value) {
        return;
      }
      fetchingCurrentValid.value = true;
      try {
        const data = await configGetCurrentValid();
        currentValid.value = JSON.stringify(data, null, 2);
      } catch (err) {
        showError(message, err);
      } finally {
        fetchingCurrentValid.value = false;
      }
    };
    return {
      fetchValid,
      tab,
      currentValid,
      fetchingCurrentValid,
    };
  },
  render() {
    const { tab, fetchValid, currentValid, fetchingCurrentValid } = this;
    return (
      <NTabs
        defaultValue={tab}
        onUpdateValue={(value) => {
          if (value === validTab) {
            fetchValid();
          }
        }}
      >
        <NTabPane name="所有配置">
          <ExConfigTable />
        </NTabPane>
        <NTabPane name={validTab}>
          {fetchingCurrentValid && <ExLoading />}
          {!fetchingCurrentValid && <pre>{currentValid}</pre>}
        </NTabPane>
      </NTabs>
    );
  },
});
