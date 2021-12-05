import { css } from "@linaria/core";
import { EditRegular } from "@vicons/fa";
import {
  NButton,
  NCard,
  NIcon,
  NPageHeader,
  NSpin,
  useMessage,
} from "naive-ui";
import { TableColumn } from "naive-ui/lib/data-table/src/interface";
import { defineComponent, onUnmounted, ref } from "vue";
import ExForm, { FormItemTypes } from "../components/ExForm";
import ExTable, { newLevelValueColumn } from "../components/ExTable";
import { diff, showError, showWarning } from "../helpers/util";
import useUserState, {
  userList,
  userListClear,
  userUpdateByID,
  UserAccount,
} from "../states/user";

const userRoleClass = css`
  margin: 0;
  padding: 0;
  list-style-position: insied;
`;

function getColumns(): TableColumn[] {
  return [
    {
      title: "账户",
      key: "account",
    },
    {
      title: "角色",
      key: "roles",
      render: (row: Record<string, unknown>) => {
        if (!row.roles) {
          return null;
        }
        const arr = (row.roles as string[]).map((role) => {
          return <li>{role}</li>;
        });
        return <ul class={userRoleClass}>{arr}</ul>;
      },
    },
    newLevelValueColumn({
      title: "状态",
      key: "status.desc",
    }),
    {
      title: "邮箱",
      key: "email",
    },
  ];
}

function getFilters() {
  return [
    {
      name: "角色：",
      key: "role",
      placeholder: "请选择要筛选的用户角色",
      type: FormItemTypes.Select,
      options: [
        {
          label: "所有",
          value: "",
        },
        {
          label: "超级用户",
          value: "su",
        },
        {
          label: "管理员",
          value: "admin",
        },
      ],
    },
    {
      name: "状态：",
      key: "status",
      placeholder: "请选择要筛选的账户状态",
      type: FormItemTypes.Select,
      options: [
        {
          label: "所有",
          value: "",
        },
        {
          label: "启用",
          value: "1",
        },
        {
          label: "禁用",
          value: "2",
        },
      ],
    },
    {
      name: "关键字：",
      key: "keyword",
      placeholder: "请输入搜索关键字",
    },
  ];
}

function getUpdateFormItems(updatedUser: UserAccount) {
  return [
    {
      name: "账号：",
      key: "account",
      disabled: true,
      defaultValue: updatedUser.account,
    },
    {
      name: "状态：",
      key: "status",
      type: FormItemTypes.Select,
      defaultValue: updatedUser.status.value,
      placeholder: "请选择账户状态",
      options: [
        {
          label: "启用",
          value: 1,
        },
        {
          label: "禁用",
          value: 2,
        },
      ],
    },
    {
      name: "角色：",
      key: "roles",
      type: FormItemTypes.MultiSelect,
      defaultValue: updatedUser.roles,
      placeholder: "请选择账户角色",
      options: [
        {
          label: "超级用户",
          value: "su",
        },
        {
          label: "管理员",
          value: "admin",
        },
      ],
    },
  ];
}

const listMode = "list";
const updateMode = "update";

export default defineComponent({
  name: "UsersView",
  setup() {
    const message = useMessage();
    const { users, info } = useUserState();
    const mode = ref(listMode);
    const updatedUser = ref({} as UserAccount);
    const updating = ref(false);

    const fetchUsers = async (params: {
      limit: number;
      offset: number;
      keyword?: string;
      status?: string;
      role?: string;
    }) =>
      userList({
        limit: params.limit,
        offset: params.offset,
        keyword: params.keyword || "",
        status: params.status || "",
        role: params.role || "",
      });

    const update = async (data: Record<string, unknown>) => {
      const diffInfo = diff(data, updatedUser.value);
      if (diffInfo.modifiedCount === 0) {
        showWarning(message, "请先更新信息");
        return;
      }
      if (updating.value) {
        return;
      }
      try {
        updating.value = true;
        await userUpdateByID({
          id: updatedUser.value.id as number,
          data: diffInfo.data,
        });
        mode.value = listMode;
      } catch (err) {
        showError(message, err);
      } finally {
        updating.value = false;
      }
    };

    onUnmounted(() => {
      userListClear();
    });

    return {
      mode,
      updatedUser,
      userInfo: info,
      users,
      updating,
      fetchUsers,
      update,
    };
  },
  render() {
    const { users, fetchUsers, userInfo, mode, updatedUser, update, updating } =
      this;
    if (mode === updateMode) {
      const formItems = getUpdateFormItems(updatedUser);
      return (
        <NSpin show={updating}>
          <NCard>
            <NPageHeader
              title="用户信息更新"
              onBack={() => {
                this.mode = listMode;
              }}
            >
              <ExForm
                formItems={formItems}
                onSubmit={update}
                submitText="更新"
              />
            </NPageHeader>
          </NCard>
        </NSpin>
      );
    }
    const columns = getColumns();
    const { roles } = userInfo;
    if (roles.includes("su") || roles.includes("admin")) {
      const render = (row: Record<string, unknown>) => {
        return (
          <NButton
            bordered={false}
            onClick={() => {
              this.updatedUser = row as UserAccount;
              this.mode = updateMode;
            }}
          >
            <NIcon class="mright5">
              <EditRegular />
            </NIcon>
            更新
          </NButton>
        );
      };
      columns.push({
        title: "操作",
        key: "op",
        render,
      });
    }
    return (
      <ExTable
        title={"用户查询"}
        columns={columns}
        filters={getFilters()}
        data={users}
        fetch={fetchUsers}
      />
    );
  },
});
