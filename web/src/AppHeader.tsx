import { User } from "@vicons/fa";
import {
  NButton,
  NH2,
  NIcon,
  NLayoutHeader,
  NSpace,
  NText,
  useMessage,
} from "naive-ui";
import { css } from "@linaria/core";
import { defineComponent, onBeforeMount } from "vue";
import { mainHeaderHeight, padding } from "./constants/style";
import { showError } from "./helpers/util";
import { goToHome, goToLogin, goToRegister, goToProfile } from "./routes";
import useUserState, { userFetchInfo, userLogout } from "./states/user";
import useCommonState, {
  commonGetSettings,
  commonUpdateSettingTheme,
} from "./states/common";

const userInfoClass = css`
  margin-right: 5px;
`;
const headerClass = css`
  height: ${mainHeaderHeight}px;
  line-height: ${mainHeaderHeight}px;
  padding: 0 ${3 * padding}px;
`;

const logoClass = css`
  float: left;
  cursor: pointer;
`;

export default defineComponent({
  name: "AppHeader",
  setup() {
    const { info } = useUserState();
    const { settings } = useCommonState();
    const message = useMessage();
    onBeforeMount(async () => {
      try {
        await commonGetSettings();
        await userFetchInfo();
      } catch (err) {
        showError(message, err);
      } finally {
        if (!info.account) {
          goToLogin();
        }
      }
    });

    // 退出登录
    const logout = async () => {
      try {
        await userLogout();
      } catch (err) {
        showError(message, err);
      }
    };

    // 主题选择
    const renderToggleTheme = () => {
      const isDark = settings.theme === "dark";
      const toggleTheme = isDark ? "light" : "dark";
      const text = isDark ? "浅色" : "深色";
      return (
        <NButton
          bordered={false}
          onClick={async () => {
            try {
              await commonUpdateSettingTheme(toggleTheme);
            } catch (err) {
              showError(message, err);
            }
          }}
        >
          {text}
        </NButton>
      );
    };

    // 用户信息
    const renderUserInfo = () => (
      <>
        <NButton bordered={false} onClick={() => goToProfile()}>
          <NIcon class={userInfoClass}>
            <User />
          </NIcon>
          {info.account}
        </NButton>
        <NButton bordered={false} onClick={logout}>
          退出登录
        </NButton>
      </>
    );
    // 登录注册等功能按钮
    const renderCtrls = () => (
      <>
        <NButton bordered={false} onClick={() => goToLogin()}>
          登录
        </NButton>
        <NButton bordered={false} onClick={() => goToRegister()}>
          注册
        </NButton>
      </>
    );
    return {
      renderToggleTheme,
      renderUserInfo,
      renderCtrls,
      userInfo: info,
    };
  },
  render() {
    const { processing, account } = this.userInfo;
    return (
      <NLayoutHeader bordered class={headerClass}>
        <NText tag="div" class={logoClass}>
          <NH2>
            <a
              href="#"
              style="color: var(--primary-color);text-decoration: none;"
              onClick={(e) => {
                e.preventDefault();
                goToHome();
              }}
            >
              tiny-site
            </a>
          </NH2>
        </NText>
        <NSpace justify="end">
          {this.renderToggleTheme()}
          {processing && <span>正在加载中，请稍候...</span>}
          {!processing && account != "" && this.renderUserInfo()}
          {!processing && account === "" && this.renderCtrls()}
        </NSpace>
      </NLayoutHeader>
    );
  },
});
