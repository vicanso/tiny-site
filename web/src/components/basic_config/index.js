import React from "react";
import {
  Card,
  Typography,
  Col,
  Form,
  DatePicker,
  Spin,
  notification,
  message
} from "antd";
import moment from "moment";
import PropTypes from "prop-types";

import "./basic_config.sass";
import ConfigEditor from "../config_editor";
import * as configService from "../../services/configuration";
import { isAdminUser } from "../../helpers/util";

const { Paragraph } = Typography;

const mockTime = "mockTime";

class BasicConfig extends React.Component {
  state = {
    mockTime: "",
    mockTimeConfig: {
      name: "mockTime"
    },

    signedKeys: null,

    loading: false
  };
  componentWillMount() {
    this.fetchUserInfo(this.props);
  }
  componentWillReceiveProps(props) {
    this.fetchUserInfo(props);
  }
  async fetchUserInfo(props) {
    const { roles } = props;
    if (!isAdminUser(roles)) {
      notification.open({
        message: "请使用管理员登录",
        description: "此功能需要先登录并有管理权限才可使用"
      });
      return;
    }

    this.setState({
      loading: true
    });
    try {
      const configs = await configService.list({
        name: [mockTime].join(",")
      });
      configs.forEach(item => {
        if (item.name === mockTime) {
          this.setState({
            mockTimeConfig: item,
            mockTime: item.data
          });
        }
      });
    } catch (err) {
      message.error(err.message);
    } finally {
      this.setState({
        loading: false
      });
    }
  }
  renderMockTimeConfig() {
    const { mockTimeConfig } = this.state;
    let time = null;
    if (mockTimeConfig && mockTimeConfig.data) {
      time = moment(mockTimeConfig.data);
    }
    const content = (
      <Col span={8}>
        <Form.Item label="时间配置">
          <DatePicker
            defaultValue={time}
            className="datePicker"
            showTime
            placeholder="请选择要要设置的时间"
            onChange={date => {
              let value = "";
              if (date) {
                value = date.toISOString();
              }
              this.setState({
                mockTime: value
              });
            }}
          ></DatePicker>
        </Form.Item>
      </Col>
    );
    return (
      <Card size="small" title="时间Mock">
        <Paragraph>针对应用时间Mock，用于测试环境中调整应用时间</Paragraph>
        <ConfigEditor
          originalData={mockTimeConfig}
          content={content}
          getConfigData={() => {
            return this.state.mockTime;
          }}
        />
      </Card>
    );
  }
  renderConfig() {
    const { loading } = this.state;
    if (loading) {
      return (
        <div className="loadingWrapper">
          <Spin spinning={loading} tip={"加载中..."}></Spin>
        </div>
      );
    }
    return <div>{this.renderMockTimeConfig()}</div>;
  }
  render() {
    return <div className="BasicConfig">{this.renderConfig()}</div>;
  }
}

BasicConfig.propTypes = {
  account: PropTypes.string.isRequired,
  roles: PropTypes.array.isRequired
};

export default BasicConfig;
