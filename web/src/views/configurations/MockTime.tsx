import { useMessage } from "naive-ui";
import { defineComponent, ref } from "vue";
import ExConfigEditor, {
  getDefaultFormItems,
  getDefaultFormRules,
  newRequireRule,
} from "../../components/ExConfigEditor";
import { FormItemTypes } from "../../components/ExForm";
import ExLoading from "../../components/ExLoading";
import { showError } from "../../helpers/util";
import { ConfigCategory, configGetMockTime } from "../../states/configs";

export default defineComponent({
  name: "MockTimeConfigs",
  setup() {
    const message = useMessage();
    const id = ref(0);
    const processing = ref(true);

    const fetchMockTimeID = async () => {
      processing.value = true;
      try {
        const data = await configGetMockTime();
        if (data.id) {
          id.value = data.id;
        }
      } catch (err) {
        showError(message, err);
      } finally {
        processing.value = false;
      }
    };
    fetchMockTimeID();

    return {
      id,
      processing,
    };
  },
  render() {
    const { id, processing } = this;
    if (processing) {
      return <ExLoading />;
    }
    const formItems = getDefaultFormItems({
      category: ConfigCategory.MockTime,
      name: ConfigCategory.MockTime,
    });
    formItems.push({
      name: "时间配置：",
      key: "data",
      type: FormItemTypes.DateTime,
      placeholder: "请选择要Mock的时间",
    });
    const rules = getDefaultFormRules({
      data: newRequireRule("Mock时间不能为空"),
    });
    return (
      <ExConfigEditor
        id={id}
        title="添加/更新MockTime配置"
        description="针对应用时间Mock，用于测试环境中调整应用时间"
        formItems={formItems}
        rules={rules}
      />
    );
  },
});
