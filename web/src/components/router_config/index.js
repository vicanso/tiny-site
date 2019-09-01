import React from "react";
import PropTypes from "prop-types";
import {
  Button,
  Card,
  Typography,
  Select,
  Form,
  Col,
  notification,
  Input,
  message
} from "antd";

import ConfigEditor from "../config_editor";
import ConfigTable from "../config_table";
import "./router_config.sass";
import * as routerService from "../../services/router";
import { isAdminUser } from "../../helpers/util";

const { Paragraph } = Typography;
const Option = Select.Option;
const { TextArea } = Input;
const editMode = "edit";

const routerConfigCategory = "routerConfig";

class RouterConfig extends React.Component {
  state = {
    routers: null,
    mode: "",

    method: "",
    route: "",
    status: null,
    contentType: "",
    response: "",

    currentData: null,
    currentKeys: null
  };
  async componentWillMount() {
    try {
      const data = await routerService.list();
      this.setState({
        routers: data.routers
      });
    } catch (err) {
      message.error(err);
    }
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
  reset() {
    this.setState({
      mode: "",
      currentData: null
    });
  }
  renderConfigEditor() {
    const {
      mode,
      currentData,
      routers,
      method,
      route,
      status,
      contentType,
      response
    } = this.state;
    if (mode !== editMode) {
      return;
    }
    const originalData = currentData || {
      category: routerConfigCategory
    };
    const opts = routers.map(item => {
      const { method, path } = item;
      const key = `${method} ${path}`;
      return (
        <Option key={key} value={key}>
          {key}
        </Option>
      );
    });
    const contentTypes = [
      "application/json; charset=UTF-8",
      "text/plain; charset=UTF-8",
      "text/html; charset=UTF-8"
    ].map(item => {
      return (
        <Option key={item} value={item}>
          {item}
        </Option>
      );
    });
    const colSpan = 8;
    const content = (
      <div>
        <Col span={colSpan}>
          <Form.Item label="路由选择">
            <Select
              showSearch
              placeholder="请选择要配置的路由"
              defaultValue={`${method} ${route}`}
              onSelect={value => {
                const arr = value.split(" ");
                this.setState({
                  method: arr[0],
                  route: arr[1]
                });
              }}
            >
              {opts}
            </Select>
          </Form.Item>
        </Col>
        <Col span={colSpan}>
          <Form.Item label="状态码">
            <Input
              defaultValue={status}
              type="number"
              placeholder="请输入响应状态码"
              onChange={e => {
                this.setState({
                  status: e.target.valueAsNumber
                });
              }}
            />
          </Form.Item>
        </Col>
        <Col span={colSpan}>
          <Form.Item label="响应类型">
            <Select
              placeholder="请选择响应数据类型"
              defaultValue={contentType}
              onChange={value => {
                this.setState({
                  contentType: value
                });
              }}
            >
              {contentTypes}
            </Select>
          </Form.Item>
        </Col>
        <Col span={24}>
          <Form.Item label="响应数据">
            <TextArea
              autosize={{
                minRows: 4
              }}
              defaultValue={response}
              onChange={e => {
                this.setState({
                  response: e.target.value.trim()
                });
              }}
            />
          </Form.Item>
        </Col>
      </div>
    );
    return (
      <Card size="small" title="添加/更新路由配置">
        <Paragraph>用于设置路由启用、禁用等配置</Paragraph>
        <ConfigEditor
          originalData={originalData}
          content={content}
          getConfigData={() => {
            const { method, route, status, response, contentType } = this.state;
            const data = {
              method,
              route,
              status,
              contentType,
              response
            };
            return JSON.stringify(data);
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
          category: routerConfigCategory
        }}
        formatData={data => {
          if (!data) {
            return "";
          }
          const result = JSON.parse(data);
          if (result.response && result.response[0] === "{") {
            result.response = JSON.parse(result.response);
          }
          return <pre>{JSON.stringify(result, null, 2)}</pre>;
        }}
        onUpdate={data => {
          const routerData = JSON.parse(data.data);
          this.setState({
            method: routerData.method,
            route: routerData.route,
            status: routerData.status,
            contentType: routerData.contentType,
            response: routerData.response,

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
                method: "",
                route: "",
                status: null,
                contentType: "",
                response: "",

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
    return <div className="RouterConfig">{this.renderContent()}</div>;
  }
}

RouterConfig.propTypes = {
  account: PropTypes.string.isRequired,
  roles: PropTypes.array.isRequired
};

export default RouterConfig;
