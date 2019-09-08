import React from "react";
import PropTypes from "prop-types";
import { Menu, Icon } from "antd";
import { Link, withRouter } from "react-router-dom";

import {
  HOME_PATH,
  USER_PATH,
  USER_LOGIN_RECORDS_PATH,
  ALL_CONFIG_PATH,
  BASIC_CONFIG_PATH,
  ROUTER_CONFIG_PATH,
  IP_BLOCK_CONFIG_PATH,
  SIGNED_KEYS_CONFIG_PATH,
  FILE_ZONE_PATH,
  MY_FILE_ZONE_PATH
} from "../../paths";
import "./app_menu.sass";

const { SubMenu } = Menu;

const configMenu = {
  key: "configuration",
  title: (
    <span>
      <Icon type="setting" />
      <span>配置</span>
    </span>
  ),
  children: [
    {
      key: "all-config",
      url: ALL_CONFIG_PATH,
      title: "所有配置"
    },
    {
      key: "basic-config",
      url: BASIC_CONFIG_PATH,
      title: "基本配置"
    },
    {
      key: "router-config",
      url: ROUTER_CONFIG_PATH,
      title: "路由配置"
    },
    {
      key: "signed-keys-config",
      url: SIGNED_KEYS_CONFIG_PATH,
      title: "签名配置"
    },
    {
      key: "ip-block-config",
      url: IP_BLOCK_CONFIG_PATH,
      title: "黑名单IP"
    }
  ]
};

const userMenu = {
  key: "user",
  title: (
    <span>
      <Icon type="user" />
      <span>用户</span>
    </span>
  ),
  children: [
    {
      key: "users",
      url: USER_PATH,
      title: "用户列表"
    },
    {
      key: "users-login-records",
      url: USER_LOGIN_RECORDS_PATH,
      title: "用户登录查询"
    }
  ]
};

const fileMenu = {
  key: "file",
  title: (
    <span>
      <Icon type="file" />
      <span>文件</span>
    </span>
  ),
  children: [
    {
      key: "filezones",
      url: FILE_ZONE_PATH,
      title: "所有文件空间"
    },
    {
      key: "my-filezones",
      url: MY_FILE_ZONE_PATH,
      title: "我的文件空间"
    }
  ]
};

const menuList = [configMenu, userMenu, fileMenu];

class AppMenu extends React.Component {
  state = {
    defaultOpenKeys: null,
    defaultSelectedKeys: null
  };
  constructor(props) {
    super(props);
    const { pathname } = props.location;
    const defaultSelectedKeys = [];
    const defaultOpenKeys = [];
    menuList.forEach(menu => {
      menu.children.forEach(item => {
        if (item.url === pathname) {
          defaultSelectedKeys.push(item.key);
          defaultOpenKeys.push(menu.key);
        }
      });
    });

    this.state.defaultSelectedKeys = defaultSelectedKeys;
    this.state.defaultOpenKeys = defaultOpenKeys;
  }
  renderMenus(data) {
    const arr = data.children.map(item => {
      return (
        <Menu.Item key={item.key}>
          <Link to={item.url}>{item.title}</Link>
        </Menu.Item>
      );
    });
    return (
      <SubMenu key={data.key} title={data.title}>
        {arr}
      </SubMenu>
    );
  }
  renderAllMenu() {
    const { account } = this.props;
    if (!account) {
      return;
    }
    const arr = menuList.map(item => {
      return this.renderMenus(item);
    });
    return arr;
  }
  render() {
    const { defaultSelectedKeys, defaultOpenKeys } = this.state;
    return (
      <div className="AppMenu">
        <Link className="logo" to={HOME_PATH}>
          <Icon type="home" />
          tiny-site
        </Link>
        <Menu
          mode="inline"
          theme="dark"
          defaultOpenKeys={defaultOpenKeys}
          defaultSelectedKeys={defaultSelectedKeys}
        >
          {this.renderAllMenu()}
        </Menu>
      </div>
    );
  }
}

AppMenu.propTypes = {
  account: PropTypes.string
};

export default withRouter(AppMenu);
