import { NLayout, NLayoutSider, useLoadingBar } from "naive-ui";
import { css } from "@linaria/core";
import { defineComponent, onMounted } from "vue";
import AppHeader from "./AppHeader";
import AppNavigation from "./AppNavigation";
import {
  mainHeaderHeight,
  mainNavigationWidth,
  padding,
} from "./constants/style";
import "./main.css";
import { setLoadingEvent } from "./routes/router";
import useCommonState, { commonUpdateSettingCollapsed } from "./states/common";

const layoutClass = css`
  top: ${mainHeaderHeight}px !important;
`;

const contentLayoutClass = css`
  padding: ${2 * padding}px;
`;

export default defineComponent({
  name: "App",
  setup() {
    const { settings } = useCommonState();
    const loadingBar = useLoadingBar();
    if (loadingBar != undefined) {
      setLoadingEvent(loadingBar.start, loadingBar.finish);
      onMounted(() => {
        loadingBar.finish();
      });
    }
    return {
      settings,
    };
  },
  render() {
    const { settings } = this;
    return (
      <div>
        <AppHeader />
        <NLayout hasSider position="absolute" class={layoutClass}>
          <NLayoutSider
            bordered
            collapseMode="width"
            collapsed={settings.collapsed}
            collapsedWidth={64}
            width={mainNavigationWidth}
            showTrigger
            onCollapse={() => {
              commonUpdateSettingCollapsed(true);
            }}
            onExpand={() => {
              commonUpdateSettingCollapsed(false);
            }}
          >
            <AppNavigation />
          </NLayoutSider>
          <NLayout class={contentLayoutClass}>
            <router-view />
          </NLayout>
        </NLayout>
      </div>
    );
  },
});
