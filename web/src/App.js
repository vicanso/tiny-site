import React from "react";
import { Route, HashRouter } from "react-router-dom";
import axios from "axios";
import { message, Spin } from "antd";

import "./app.sass";
import {
  ALL_CONFIG_PATH,
  BASIC_CONFIG_PATH,
  SIGNED_KEYS_CONFIG_PATH,
  ROUTER_CONFIG_PATH,
  IP_BLOCK_CONFIG_PATH,
  REGISTER_PATH,
  LOGIN_PATH,
  USER_PATH,
  USER_LOGIN_RECORDS_PATH,
  HOME_PATH,
  FILE_ZONE_PATH
} from "./paths";
import { USERS_ME } from "./urls";
import AppMenu from "./components/app_menu";
import AppHeader from "./components/app_header";
import BasicConfig from "./components/basic_config";
import SignedKeysConfig from "./components/signed_keys_config";
import Login from "./components/login";
import Register from "./components/register";
import RouterConfig from "./components/router_config";
import ConfigList from "./components/config_list";
import UserList from "./components/user_list";
import UserLoginRecordList from "./components/user_login_record_list";
import IPBlockList from "./components/ip_block_list";
import Home from "./components/home";
import FileZone from "./components/file_zone";

function NeedLoginRoute({ component: Component, account, roles, ...rest }) {
  return (
    <Route
      {...rest}
      render={props => {
        const { history } = props;
        if (!account) {
          history.push(LOGIN_PATH);
          return;
        }
        return <Component {...props} account={account} roles={roles} />;
      }}
    />
  );
}

class App extends React.Component {
  state = {
    loading: false,
    account: "",
    roles: null
  };
  async componentWillMount() {
    this.setState({
      loading: true
    });
    try {
      const { data } = await axios.get(USERS_ME);
      this.setUserInfo(data);
    } catch (err) {
      message.error(err.message);
    } finally {
      this.setState({
        loading: false
      });
    }
    // 更新session与cookie有效期
    setTimeout(() => {
      axios.patch(USERS_ME);
    }, 5 * 1000);
  }
  setUserInfo(data) {
    this.setState({
      account: data.account || "",
      roles: data.roles || []
    });
  }
  render() {
    const { account, roles, loading } = this.state;
    return (
      <div className="App">
        <HashRouter>
          <AppMenu account={account} />
          {loading && (
            <div className="loadingWrapper">
              <Spin tip="加载中..." />
            </div>
          )}
          {!loading && (
            <div className="contentWrapper">
              <AppHeader
                account={account}
                setUserInfo={this.setUserInfo.bind(this)}
              />

              <Route
                path={LOGIN_PATH}
                render={props => (
                  <Login {...props} setUserInfo={this.setUserInfo.bind(this)} />
                )}
              />
              <Route path={REGISTER_PATH} component={Register} />
              <NeedLoginRoute
                path={ALL_CONFIG_PATH}
                component={ConfigList}
                account={account}
                roles={roles}
              />
              <NeedLoginRoute
                path={BASIC_CONFIG_PATH}
                component={BasicConfig}
                account={account}
                roles={roles}
              />
              <NeedLoginRoute
                path={SIGNED_KEYS_CONFIG_PATH}
                component={SignedKeysConfig}
                account={account}
                roles={roles}
              />
              <NeedLoginRoute
                path={ROUTER_CONFIG_PATH}
                component={RouterConfig}
                account={account}
                roles={roles}
              />
              <NeedLoginRoute
                exact
                path={USER_PATH}
                component={UserList}
                account={account}
                roles={roles}
              />
              <NeedLoginRoute
                path={USER_LOGIN_RECORDS_PATH}
                component={UserLoginRecordList}
                account={account}
                roles={roles}
              />
              <NeedLoginRoute
                path={IP_BLOCK_CONFIG_PATH}
                component={IPBlockList}
                account={account}
                roles={roles}
              />
              <NeedLoginRoute
                path={FILE_ZONE_PATH}
                component={FileZone}
                account={account}
                roles={roles}
              />
              <Route path={HOME_PATH} component={Home} exact />
            </div>
          )}
        </HashRouter>
      </div>
    );
  }
}

export default App;
