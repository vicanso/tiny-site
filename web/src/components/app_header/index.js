import React from "react";
import { Spin, message, Icon } from "antd";
import { Link } from "react-router-dom";
import PropTypes from "prop-types";

import { LOGIN_PATH, REGISTER_PATH } from "../../paths";
import "./app_header.sass";
import * as userService from "../../services/user";

class AppHeader extends React.Component {
  state = {
    loading: false
  };
  async logout(e) {
    const { setUserInfo } = this.props;
    e.preventDefault();
    this.setState({
      loading: true
    });
    try {
      await userService.logout();
      setUserInfo({
        account: ""
      });
    } catch (err) {
      message.error(err.message);
    } finally {
      this.setState({
        loading: false
      });
    }
  }
  renderFunctions() {
    const { account } = this.props;
    const { loading } = this.state;
    let content = null;
    if (loading) {
      content = <Spin spinning={loading}></Spin>;
    } else if (account) {
      content = (
        <div>
          <span>
            <Icon type="user-add" />
            {account}
          </span>
          <Link
            className="function"
            to={"/logout"}
            onClick={this.logout.bind(this)}
          >
            <Icon type="logout" />
            注销
          </Link>
        </div>
      );
    } else {
      content = (
        <div>
          <Link className="function" to={LOGIN_PATH}>
            <Icon type="login" />登 录
          </Link>
          <Link className="function" to={REGISTER_PATH}>
            <Icon type="user-add" />注 册
          </Link>
        </div>
      );
    }
    return <div className="functions">{content}</div>;
  }
  render() {
    return <div className="AppHeader">{this.renderFunctions()}</div>;
  }
}

AppHeader.propTypes = {
  account: PropTypes.string.isRequired,
  setUserInfo: PropTypes.func.isRequired
};

export default AppHeader;
