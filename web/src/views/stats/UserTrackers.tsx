import { NEllipsis, useMessage } from "naive-ui";
import { TableColumn } from "naive-ui/lib/data-table/src/interface";
import { defineComponent, onMounted, onUnmounted } from "vue";
import ExFluxDetail from "../../components/ExFluxDetail";
import { FormItemTypes } from "../../components/ExForm";
import ExLoading from "../../components/ExLoading";
import ExTable from "../../components/ExTable";
import { formatJSON, showError, today } from "../../helpers/util";
import useFluxState, {
  fluxListUserTrackAction,
  fluxListUserTracker,
  fluxListUserTrackerClear,
  measurementUserTracker,
} from "../../states/flux";

function getColumns(): TableColumn[] {
  return [
    {
      title: "账户",
      key: "account",
      width: 100,
    },
    {
      type: "expand",
      expandable: (data: Record<string, unknown>) =>
        !!data.form || !!data.query || !!data.params,
      renderExpand: (data: Record<string, unknown>) => {
        const arr = [];
        if (data.form) {
          arr.push(<pre>form: {formatJSON(data.form as string)}</pre>);
        }
        if (data.query) {
          arr.push(<pre>query: {formatJSON(data.query as string)}</pre>);
        }
        if (data.params) {
          arr.push(<pre>params: {formatJSON(data.params as string)}</pre>);
        }
        return <div>{arr}</div>;
      },
    },
    {
      title: "类别",
      key: "action",
      width: 160,
    },
    {
      title: "结果",
      key: "resultDesc",
      width: 80,
    },
    {
      title: "TrackID",
      key: "tid",
      width: 120,
      ellipsis: true,
    },
    {
      title: "SessionID",
      key: "sid",
      width: 120,
      ellipsis: true,
    },
    {
      title: "IP",
      key: "ip",
      width: 100,
      ellipsis: {
        tooltip: true,
      },
    },
    {
      title: "出错信息",
      key: "error",
      render: (row: Record<string, unknown>) => {
        const text = row.error as string;
        if (!text) {
          return;
        }
        const tooltip = {
          width: 250,
        };
        return <NEllipsis tooltip={tooltip}>{text}</NEllipsis>;
      },
    },
    {
      title: "完整记录",
      key: "userTrackerDetail",
      width: 90,
      align: "center",
      render(row: Record<string, unknown>) {
        return (
          <ExFluxDetail
            measurement={measurementUserTracker}
            data={row}
            tagKeys={["action", "result", "step"]}
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

const actionOptions = [
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
      key: "action",
      name: "类别：",
      placeholder: "请选择要筛选的类别",
      type: FormItemTypes.Select,
      options: actionOptions,
    },
    {
      key: "result",
      name: "结果：",
      type: FormItemTypes.Select,
      placeholder: "请选择要筛选的结果",
      options: [
        {
          label: "所有",
          value: "",
        },
        {
          label: "成功",
          value: "0",
        },
        {
          label: "失败",
          value: "1",
        },
      ],
    },
    {
      key: "limit",
      name: "查询数量",
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
  name: "UserTrackerStats",
  setup() {
    const message = useMessage();
    const { userTrackers, userTrackerActions } = useFluxState();
    // 加载用户行为类别
    onMounted(async () => {
      try {
        await fluxListUserTrackAction();
      } catch (err) {
        showError(message, err);
      }
    });
    // 清除数据
    onUnmounted(() => {
      fluxListUserTrackerClear();
    });

    return {
      userTrackers,
      userTrackerActions,
    };
  },
  render() {
    const { userTrackers, userTrackerActions } = this;
    if (userTrackerActions.processing) {
      return <ExLoading />;
    }
    // 添加类别选项
    if (actionOptions.length === 1 && userTrackerActions.items.length !== 0) {
      userTrackerActions.items.forEach((item) =>
        actionOptions.push({
          label: item,
          value: item,
        })
      );
    }
    return (
      <ExTable
        disableAutoFetch={true}
        hidePagination={true}
        filters={getFilters()}
        title={"用户行为查询"}
        columns={getColumns()}
        data={userTrackers}
        fetch={fluxListUserTracker}
      />
    );
  },
});
