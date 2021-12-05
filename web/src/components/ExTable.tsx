import {
  NButton,
  NCard,
  NDataTable,
  NDatePicker,
  NForm,
  NFormItem,
  NGrid,
  NGridItem,
  NInput,
  NInputNumber,
  NPagination,
  NSelect,
  NSpin,
  NIcon,
  useMessage,
} from "naive-ui";
import { EditRegular } from "@vicons/fa";
import { TableColumn } from "naive-ui/lib/data-table/src/interface";
import { Value } from "naive-ui/lib/select/src/interface";
import { Component, defineComponent, onMounted, ref, PropType } from "vue";
import { css } from "@linaria/core";
import { padding } from "../constants/style";
import { getDaysAgo, showError, today, yesterday } from "../helpers/util";
import { FormItemTypes } from "./ExForm";

interface Filter {
  type?: string;
  key: string;
  name: string;
  span?: number;
  placeholder?: string;
  options?: [];
  defaultValue?: unknown;
}

interface TableData {
  processing: boolean;
  items: [];
  count: number;
}

const paginationClass = css`
  margin-top: ${padding}px;
  float: right;
`;
const ulClass = css`
  padding: 0;
  list-style: inside;
`;

// 列表形式的column
export function newListColumn(params: {
  key: string;
  title: string;
}): TableColumn {
  return {
    title: params.title,
    key: params.key,
    render: (row: Record<string, unknown>) => {
      const values = row[params.key] as unknown[];
      if (!values || values.length === 0) {
        return "--";
      }
      const arr = values.map((value) => <li>{value}</li>);
      return <ul class={ulClass}>{arr}</ul>;
    },
  };
}

export function newLevelValueColumn(params: {
  key: string;
  title: string;
  width?: number;
}): TableColumn {
  return {
    title: params.title,
    width: params.width,
    key: params.key,
    render: (row: Record<string, unknown>) => {
      const keys = params.key.split(".");
      const lastKey = keys.pop() || "";
      let data = row;
      keys.forEach((key) => (data = data[key] as Record<string, unknown>));
      return <span>{data[lastKey]}</span>;
    },
  };
}

// 操作相关
export function newOPColumn(
  fn: (params: Record<string, unknown>) => void
): TableColumn {
  return {
    title: "操作",
    key: "op",
    render(row: Record<string, unknown>) {
      return (
        <NButton
          bordered={false}
          onClick={() => {
            fn(row);
          }}
        >
          <NIcon>
            <EditRegular />
          </NIcon>
          更新
        </NButton>
      );
    },
  };
}

function isJSON(str: string): boolean {
  if (str.length < 2) {
    return false;
  }
  const firstLetter = str[0];
  const lastLetter = str[str.length - 1];
  // 示判断{]的场景
  if (firstLetter !== "{" && firstLetter !== "[") {
    return false;
  }
  if (lastLetter !== "}" && lastLetter !== "]") {
    return false;
  }
  return true;
}

// json形式的展开栏
export function newJSONRenderExpand(key: string): TableColumn {
  return {
    type: "expand",
    expandable: () => true,
    renderExpand: (row: Record<string, unknown>) => {
      const str = row[key] as string;
      if (isJSON(str.trim())) {
        return <pre>{JSON.stringify(JSON.parse(str), null, 2)}</pre>;
      }
      return str;
    },
  };
}

export default defineComponent({
  name: "ExTable",
  props: {
    title: {
      type: String,
      default: "",
    },
    columns: {
      type: Array as PropType<TableColumn[]>,
      required: true,
    },
    data: {
      type: Object,
      required: true,
    },
    fetch: {
      type: Function,
      required: true,
    },
    limit: {
      type: Number,
      default: 10,
    },
    filters: {
      type: Array,
      default: () => [],
    },
    hidePagination: {
      type: Boolean,
      default: () => false,
    },
    disableAutoFetch: {
      type: Boolean,
      default: () => false,
    },
    scrollX: {
      type: Number,
      default: () => null,
    },
  },
  setup(props) {
    const message = useMessage();
    const processing = ref(false);
    const offset = ref(0);
    const filterParams = ref({} as Record<string, unknown>);
    props.filters?.forEach((item) => {
      const filter = item as Filter;
      // 无默认值的忽略
      if (!filter.defaultValue) {
        return;
      }
      switch (filter.type) {
        case FormItemTypes.DateRange: {
          const arr = filter.key.split(":");
          const values = filter.defaultValue as [];
          arr.forEach((key, index) => {
            filterParams.value[key] = values[index];
          });
          break;
        }
        default:
          filterParams.value[filter.key] = filter.defaultValue;
      }
    });

    // 拉取数据
    const fetchData = async () => {
      if (processing.value) {
        return;
      }
      try {
        processing.value = true;
        await props.fetch(
          Object.assign(
            {
              limit: props.limit,
              offset: offset.value,
            },
            filterParams.value
          )
        );
      } catch (err) {
        showError(message, err);
      } finally {
        processing.value = false;
      }
    };
    if (!props.disableAutoFetch) {
      onMounted(() => {
        fetchData();
      });
    }
    return {
      processing,
      fetchData,
      filterParams,
      offset,
    };
  },
  render() {
    const { offset, fetchData, filterParams, $slots, processing } = this;
    const { columns, data, limit, filters, title, hidePagination, scrollX } =
      this.$props;

    const tableData = data as TableData;

    // 分页
    const pageCount = Math.ceil(tableData.count / limit);
    const page = Math.floor(offset / limit) + 1;

    const ranges: Record<string, [number, number]> = {
      最近1小时: [Date.now() - 3600 * 1000, Date.now()],
      最近3小时: [Date.now() - 3 * 3600 * 1000, Date.now()],
      最近6小时: [Date.now() - 3 * 3600 * 1000, Date.now()],
      今天: [today().getTime(), Date.now()],
      昨天: [yesterday().getTime(), today().getTime() - 1],
      最近3天: [getDaysAgo(2).getTime(), Date.now()],
      最近7天: [getDaysAgo(6).getTime(), Date.now()],
    };

    // 筛选功能
    const size = "large";
    let currentSpan = 0;
    const filterGrids = filters.map((item) => {
      const filterItem = item as Filter;
      const span = filterItem.span || 6;
      currentSpan += span;
      let component: Component;
      switch (filterItem.type) {
        case FormItemTypes.Select:
          component = (
            <NSelect
              filterable
              size={size}
              defaultValue={filterItem.defaultValue as Value}
              options={filterItem.options || []}
              placeholder={filterItem.placeholder}
              onUpdateValue={(value) => {
                filterParams[filterItem.key] = value;
              }}
            ></NSelect>
          );
          break;
        case FormItemTypes.InputNumber:
          component = (
            <NInputNumber
              class="widthFull"
              size={size}
              placeholder={filterItem.placeholder}
              defaultValue={filterItem.defaultValue as number}
              onUpdate:value={(value) => {
                filterParams[filterItem.key] = value;
              }}
            />
          );
          break;
        case FormItemTypes.DateRange: {
          let defaultValue: [number, number] | null = null;
          if (filterItem.defaultValue) {
            const arr = filterItem.defaultValue as [];
            defaultValue = [Date.now(), Date.now()];
            arr.forEach((value, index) => {
              if (defaultValue) {
                defaultValue[index] = new Date(value).getTime();
              }
            });
          }
          let startPlaceholder = filterItem.placeholder;
          let endPlaceholder = filterItem.placeholder;
          if (filterItem.placeholder) {
            const arr = filterItem.placeholder.split(":");
            if (arr.length === 2) {
              startPlaceholder = arr[0];
              endPlaceholder = arr[1];
            }
          }
          component = (
            <NDatePicker
              class="widthFull"
              format="yyyy-MM-dd HH:mm:ss"
              size={size}
              defaultValue={defaultValue}
              ranges={ranges}
              type="daterange"
              startPlaceholder={startPlaceholder}
              endPlaceholder={endPlaceholder}
              clearable
              onUpdateValue={(value) => {
                const arr = filterItem.key.split(":");
                arr.forEach((key, index) => {
                  if (!value || value.length <= index) {
                    filterParams[key] = "";
                    return;
                  }
                  filterParams[key] = new Date(value[index]).toISOString();
                });
              }}
            />
          );
          break;
        }
        default:
          component = (
            <NInput
              clearable
              size={size}
              defaultValue={filterItem.defaultValue as string}
              placeholder={filterItem.placeholder}
              onUpdateValue={(value) => {
                filterParams[filterItem.key] = value.trim();
              }}
            />
          );
          break;
      }
      return (
        <NGridItem span={span}>
          <NFormItem label={filterItem.name}>{component}</NFormItem>
        </NGridItem>
      );
    });
    if (filterGrids.length !== 0) {
      const span = 24 - (currentSpan % 24);
      filterGrids.push(
        <NGridItem span={span}>
          <NButton
            class="widthFull"
            size={size}
            onClick={() => {
              this.offset = 0;
              fetchData();
            }}
          >
            筛选
          </NButton>
        </NGridItem>
      );
    }

    const table = (
      <div class="clearfix">
        <NSpin show={processing}>
          {filterGrids.length !== 0 && (
            <NForm labelPlacement="left">
              <NGrid xGap={24}>{filterGrids}</NGrid>
            </NForm>
          )}
          <NDataTable
            columns={columns}
            data={tableData.items}
            scrollX={scrollX}
          ></NDataTable>
          {!hidePagination && pageCount > 1 && (
            <NPagination
              page={page}
              pageCount={pageCount}
              class={paginationClass}
              onUpdatePage={(value) => {
                this.offset = (value - 1) * limit;
                fetchData();
              }}
            />
          )}
          {$slots.default && $slots.default()}
        </NSpin>
      </div>
    );
    if (!title) {
      return table;
    }

    return <NCard title={title}>{table}</NCard>;
  },
});
