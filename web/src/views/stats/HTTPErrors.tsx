import { useMessage } from "naive-ui";
import { defineComponent, onMounted, onUnmounted } from "vue";
import { TableColumn } from "naive-ui/lib/data-table/src/interface";
import ExTable from "../../components/ExTable";
import { showError } from "../../helpers/util";
import useFluxState, {
  fluxListHTTPCategory,
  fluxListHTTPError,
  fluxListHTTPErrorClear,
  measurementHttpError,
} from "../../states/flux";
import { today } from "../../helpers/util";
import ExLoading from "../../components/ExLoading";
import { FormItemTypes } from "../../components/ExForm";
import ExFluxDetail from "../../components/ExFluxDetail";

function getColumns(): TableColumn[] {
  return [
    {
      title: "账户",
      key: "account",
      width: 100,
      fixed: "left",
    },
    {
      title: "方法",
      key: "method",
      width: 80,
    },
    {
      title: "分类",
      key: "category",
      width: 120,
    },
    {
      title: "状态码",
      key: "status",
      width: 80,
    },
    {
      title: "异常",
      key: "exception",
      width: 80,
      render(row: Record<string, unknown>) {
        if (row.exception) {
          return "是";
        }
        return "否";
      },
    },
    {
      title: "TrackID",
      key: "tid",
      width: 220,
    },
    {
      title: "SessionID",
      key: "sid",
      width: 220,
    },
    {
      title: "IP",
      key: "ip",
      width: 140,
    },
    {
      title: "URI",
      key: "uri",
      width: 200,
    },
    {
      title: "出错信息",
      key: "error",
      width: 80,
      ellipsis: {
        tooltip: true,
      },
    },
    {
      title: "完整记录",
      key: "httpErrorDetail",
      width: 90,
      align: "center",
      render(row: Record<string, unknown>) {
        return (
          <ExFluxDetail
            measurement={measurementHttpError}
            data={row}
            tagKeys={["method", "route", "category"]}
          />
        );
      },
    },
    {
      title: "时间",
      key: "createdAt",
      width: 180,
      fixed: "right",
    },
  ];
}

// 共用的options
const categoryOptions = [
  {
    label: "所有",
    value: "",
  },
];
function getFilters() {
  return [
    {
      key: "account",
      name: "账户：",
      placeholder: "请输入要筛选的账号",
    },
    {
      key: "category",
      name: "分类：",
      placeholder: "请选择要筛选的分类",
      type: FormItemTypes.Select,
      options: categoryOptions,
    },
    {
      key: "exception",
      name: "是否异常",
      placeholder: "请选择是否筛选异常出错",
      type: FormItemTypes.Select,
      options: [
        {
          label: "所有",
          value: "",
        },
        {
          label: "是",
          value: "true",
        },
        {
          label: "否",
          value: "false",
        },
      ],
    },
    {
      key: "limit",
      name: "查询数量：",
      type: FormItemTypes.InputNumber,
      placeholder: "请输入要查询的记录数量",
    },
    {
      key: "begin:end",
      name: "开始结束时间：",
      type: FormItemTypes.DateRange,
      span: 12,
      defaultValue: [today().toISOString(), new Date().toISOString()],
    },
  ];
}

export default defineComponent({
  name: "HTTPErrorStats",
  setup() {
    const message = useMessage();
    const { httpErrorCategories, httpErrors } = useFluxState();

    // 加载http响应出错依赖服务分类（服务名称)
    onMounted(async () => {
      try {
        await fluxListHTTPCategory();
      } catch (err) {
        showError(message, err);
      }
    });

    // 清除数据
    onUnmounted(() => {
      fluxListHTTPErrorClear();
    });

    return {
      httpErrorCategories,
      httpErrors,
    };
  },
  render() {
    const { httpErrors, httpErrorCategories } = this;
    if (httpErrorCategories.processing) {
      return <ExLoading />;
    }
    if (
      categoryOptions.length === 1 &&
      httpErrorCategories.items.length !== 0
    ) {
      httpErrorCategories.items.forEach((item) => {
        categoryOptions.push({
          label: item,
          value: item,
        });
      });
    }
    return (
      <ExTable
        disableAutoFetch={true}
        hidePagination={true}
        title={"HTTP响应出错"}
        filters={getFilters()}
        columns={getColumns()}
        data={httpErrors}
        fetch={fluxListHTTPError}
      />
    );
  },
});
