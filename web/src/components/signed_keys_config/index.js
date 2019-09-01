import React from "react";
import {
  Button,
  Card,
  Typography,
  Select,
  Form,
  Col,
  message,
  notification
} from "antd";
import PropTypes from "prop-types";

import ConfigEditor from "../config_editor";
import ConfigTable from "../config_table";
import "./signed_keys_config.sass";
import * as commonSerivce from "../../services/common";
import { isAdminUser } from "../../helpers/util";

const { Paragraph } = Typography;
const signedKeyCategory = "signedKey";
const editMode = "edit";

class SignedKeysConfig extends React.Component {
  state = {
    mode: "",
    randomKey: "",
    currentData: null,
    currentKeys: null
  };
  componentDidMount() {
    const { roles } = this.props;
    if (!isAdminUser(roles)) {
      notification.open({
        message: "请使用管理员登录",
        description: "此功能需要先登录并有管理权限才可使用"
      });
    }
  }
  reset() {
    this.setState({
      mode: "",
      currentKeys: null,
      currentData: null
    });
  }
  async updateRandomString() {
    try {
      const data = await commonSerivce.getRandomKeys({
        n: 10
      });
      this.setState({
        randomKey: data.keys[0]
      });
    } catch (err) {
      message.error(err.message);
    }
  }
  renderConfigEditor() {
    const { mode, currentData, currentKeys, randomKey } = this.state;
    if (mode !== editMode) {
      return;
    }
    const originalData = currentData || {
      category: signedKeyCategory
    };
    const colSpan = 8;
    const content = (
      <div>
        <Col span={colSpan}>
          <Form.Item label="Key列表">
            <Select
              defaultValue={currentKeys || []}
              mode="tags"
              placeholder="请输入需要配置的key"
              onChange={value => {
                this.setState({
                  currentKeys: value
                });
              }}
            ></Select>
          </Form.Item>
        </Col>
        <Col span={colSpan}>
          <Form.Item label="随机串">
            <Button type="primary" onClick={this.updateRandomString.bind(this)}>
              立即生成
            </Button>
            <span
              style={{
                marginLeft: "15px"
              }}
            >
              {randomKey}
            </span>
          </Form.Item>
        </Col>
      </div>
    );
    return (
      <Card size="small" title="添加/更新签名配置">
        <Paragraph>用于生成session的cookie认证</Paragraph>
        <ConfigEditor
          originalData={originalData}
          content={content}
          getConfigData={() => {
            const { currentKeys } = this.state;
            if (!currentKeys) {
              return "";
            }
            return currentKeys.join(",");
          }}
          onSuccess={this.reset.bind(this)}
        />
        <Button className="back" onClick={this.reset.bind(this)}>
          返回
        </Button>
      </Card>
    );
  }
  renderTable() {
    const { mode } = this.state;
    if (mode) {
      return;
    }
    return (
      <ConfigTable
        params={{
          category: signedKeyCategory
        }}
        formatData={data => {
          if (!data) {
            return "[]";
          }
          return JSON.stringify(data.split(","));
        }}
        onUpdate={data => {
          this.setState({
            currentKeys: data.data.split(","),
            mode: editMode,
            currentData: data
          });
        }}
      />
    );
  }
  renderContent() {
    const { mode } = this.state;
    const { roles } = this.props;
    if (!isAdminUser(roles)) {
      return;
    }
    return (
      <div>
        {this.renderTable()}
        {this.renderConfigEditor()}
        {!mode && (
          <Button
            onClick={() => {
              this.setState({
                mode: editMode
              });
            }}
            type="primary"
            className="add"
          >
            添加
          </Button>
        )}
      </div>
    );
  }
  render() {
    return <div className="SignedKeysConfig">{this.renderContent()}</div>;
  }
}

SignedKeysConfig.propTypes = {
  account: PropTypes.string.isRequired,
  roles: PropTypes.array.isRequired
};

export default SignedKeysConfig;
