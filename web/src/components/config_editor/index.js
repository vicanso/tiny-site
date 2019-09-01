import React from "react";
import PropTypes from "prop-types";
import moment from "moment";
import {
  Form,
  Input,
  Row,
  Col,
  Select,
  DatePicker,
  Button,
  Spin,
  message
} from "antd";

import "./config_editor.sass";
import { TIME_FORMAT } from "../../vars";
import * as configService from "../../services/configuration";

const Option = Select.Option;
const editMode = "edit";

class ConfigEditor extends React.Component {
  state = {
    id: 0,
    colSpan: 8,
    submitting: false,
    mode: "",

    name: "",
    category: "",
    status: null,
    beginDate: "",
    endDate: ""
  };
  constructor(props) {
    super(props);
    if (props.originalData) {
      Object.assign(this.state, props.originalData);
    }
    if (this.state.id !== 0) {
      this.state.mode = editMode;
    }
  }
  async handleSubmit(e) {
    const { onSuccess, originalData } = this.props;
    e.preventDefault();
    const {
      submitting,
      id,
      mode,

      name,
      status,
      category,
      enabled,
      beginDate,
      endDate
    } = this.state;
    if (submitting) {
      return;
    }
    try {
      if (!name) {
        throw new Error("配置名称不能为空");
      }
      if (!beginDate || !endDate) {
        throw new Error("开始与结束日期不能为空");
      }
      if (!status) {
        throw new Error("配置状态不能为空");
      }
      this.setState({
        submitting: true
      });
      const configData = {
        name,
        status,
        category,
        enabled,
        beginDate,
        endDate,
        data: this.props.getConfigData()
      };
      if (mode === editMode) {
        Object.keys(configData).forEach(key => {
          if (!originalData) {
            return;
          }
          if (configData[key] === originalData[key]) {
            delete configData[key];
          }
        });
        if (Object.keys(configData).length !== 0) {
          await configService.updateByID(id, configData);
        }
        message.info("更新配置已成功");
      } else {
        await configService.add(configData);
        message.info("添加配置已成功");
      }
      if (onSuccess) {
        setTimeout(onSuccess, 100);
      }
    } catch (err) {
      message.error(err.message);
    } finally {
      this.setState({
        submitting: false
      });
    }
  }
  render() {
    const { content, originalData } = this.props;
    const {
      colSpan,
      submitting,
      mode,

      status,
      name,
      category,
      beginDate,
      endDate
    } = this.state;
    let beginTime = null;
    if (beginDate) {
      beginTime = moment(beginDate);
    }
    let endTime = null;
    if (endDate) {
      endTime = moment(endDate);
    }
    return (
      <Form onSubmit={this.handleSubmit.bind(this)} className="ConfigEditor">
        <Spin spinning={submitting}>
          <Row gutter={24}>
            <Col span={colSpan}>
              <Form.Item label="名称">
                <Input
                  defaultValue={name}
                  disabled={!!(originalData && originalData.name)}
                  type="text"
                  placeholder="请输入配置名称(唯一)"
                  onChange={e => {
                    this.setState({
                      name: e.target.value.trim()
                    });
                  }}
                />
              </Form.Item>
            </Col>
            <Col span={colSpan}>
              <Form.Item label="分类">
                <Input
                  defaultValue={category}
                  type="text"
                  placeholder="请输入配置分类(可选)"
                  onChange={e => {
                    this.setState({
                      category: e.target.value.trim()
                    });
                  }}
                />
              </Form.Item>
            </Col>
            <Col span={colSpan}>
              <Form.Item label="是否启用">
                <Select
                  defaultValue={status}
                  onChange={value => {
                    this.setState({
                      status: value
                    });
                  }}
                >
                  <Option value={1}>启用</Option>
                  <Option value={2}>禁用</Option>
                </Select>
                {/* <Switch
                  defaultChecked={enabled}
                  onChange={checked => {
                    this.setState({
                      enabled: checked
                    });
                  }}
                /> */}
              </Form.Item>
            </Col>
            <Col span={colSpan}>
              <Form.Item label="开始时间">
                <DatePicker
                  defaultValue={beginTime}
                  className="datePicker"
                  showTime
                  format={TIME_FORMAT}
                  placeholder="请选择开始时间"
                  onChange={date => {
                    let value = "";
                    if (date) {
                      value = date.toISOString();
                    }
                    this.setState({
                      beginDate: value
                    });
                  }}
                />
              </Form.Item>
            </Col>
            <Col span={colSpan}>
              <Form.Item label="结束时间">
                <DatePicker
                  defaultValue={endTime}
                  className="datePicker"
                  showTime
                  format={TIME_FORMAT}
                  placeholder="请选择结束时间"
                  onChange={date => {
                    let value = "";
                    if (date) {
                      value = date.toISOString();
                    }
                    this.setState({
                      endDate: value
                    });
                  }}
                />
              </Form.Item>
            </Col>
            <Col span={colSpan}>
              <Form.Item label="操作">
                <Button className="submit" type="primary" htmlType="submit">
                  {mode === editMode ? "更新" : "保存"}
                </Button>
              </Form.Item>
            </Col>
            {content}
          </Row>
        </Spin>
      </Form>
    );
  }
}

ConfigEditor.propTypes = {
  originalData: PropTypes.object,
  content: PropTypes.element.isRequired,
  getConfigData: PropTypes.func.isRequired,
  onSuccess: PropTypes.func
};

export default ConfigEditor;
