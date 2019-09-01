import React from "react";
import PropTypes from "prop-types";
import { Button, Card, Typography, Form, Col, notification, Input } from "antd";

import ConfigEditor from "../config_editor";
import ConfigTable from "../config_table";
import "./ip_block_list.sass";
import { isAdminUser } from "../../helpers/util";

const { Paragraph } = Typography;
const ipBlockCategory = "ipBlock";
const editMode = "edit";

class IPBlockList extends React.Component {
  state = {
    currentData: null,
    currentIP: ""
  };
  reset() {
    this.setState({
      mode: "",
      currentIP: "",
      currentData: null
    });
  }
  renderConfigEditor() {
    const { mode, currentData, currentIP } = this.state;
    if (mode !== editMode) {
      return;
    }
    const originalData = currentData || {
      category: ipBlockCategory
    };
    const content = (
      <div>
        <Col span={8}>
          <Form.Item label="IP">
            <Input
              defaultValue={currentIP}
              placeholder="请输入要拦截的IP或IP段"
              onChange={e => {
                this.setState({
                  currentIP: e.target.value
                });
              }}
            />
          </Form.Item>
        </Col>
      </div>
    );
    return (
      <Card size="small" title="添加/更新IP黑名单配置">
        <Paragraph>用于拦截访问IP</Paragraph>
        <ConfigEditor
          originalData={originalData}
          content={content}
          getConfigData={() => {
            return this.state.currentIP;
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
          category: ipBlockCategory
        }}
        formatData={data => {
          return data;
        }}
        onUpdate={data => {
          this.setState({
            mode: editMode,
            currentData: data,
            currentIP: data.data
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
  componentDidMount() {
    const { roles } = this.props;
    if (!isAdminUser(roles)) {
      notification.open({
        message: "请使用管理员登录",
        description: "此功能需要先登录并有管理权限才可使用"
      });
    }
  }
  render() {
    return <div className="IPBlockList">{this.renderContent()}</div>;
  }
}

IPBlockList.propTypes = {
  account: PropTypes.string.isRequired,
  roles: PropTypes.array.isRequired
};

export default IPBlockList;
