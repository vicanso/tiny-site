import React from "react";
import PropTypes from "prop-types";
import {
  Button,
  Form,
  Input,
  Row,
  Col,
  message,
  Spin,
  Table,
  Card
} from "antd";
import moment from "moment";

import "./my_file_zone.sass";
import * as fileService from "../../services/file";
import { TIME_FORMAT } from "../../vars";

const { TextArea } = Input;
const createMode = "create";
const listMode = "list";
const updateMode = "update";

class MyFileZone extends React.Component {
  state = {
    loading: true,
    mode: listMode,
    zones: null,
    // 保存需要更新的filezone数据，用于对比过滤更新字段
    updateFileZone: null,
    fileID: 0,
    fileZoneName: "",
    fileZoneOwner: "",
    fileZoneDescription: ""
  };
  componentDidMount() {
    this.fetchZoneList();
  }
  goBack() {
    this.setState({
      mode: listMode
    });
    this.fetchZoneList();
  }
  async handleFileZoneConfirm() {
    const {
      updateFileZone,
      fileZoneDescription,
      fileZoneName,
      fileZoneOwner,
      fileID,
      mode
    } = this.state;
    if (!fileZoneName || !fileZoneOwner || !fileZoneDescription) {
      message.error("文件空间名称、拥有者以及描述均不能为空");
      return;
    }
    const tips = mode === updateMode ? "更新" : "创建";
    try {
      const params = {
        name: fileZoneName,
        owner: fileZoneOwner,
        description: fileZoneDescription
      };
      if (mode === updateMode) {
        Object.keys(params).forEach(key => {
          if (params[key] === updateFileZone[key]) {
            delete params[key];
          }
        });
        if (Object.keys(params).length === 0) {
          message.error("无信息修改，请确认是否需要修改信息");
          return;
        }

        params.id = fileID;
        await fileService.updateZone(params);
      } else {
        await fileService.addZone(params);
      }
      message.info(`${tips}文件空间成功`);
      this.goBack();
    } catch (err) {
      message.error(`${tips}文件空间失败:${err.message}`);
    }
  }
  async fetchZoneList() {
    this.setState({
      loading: true
    });
    try {
      const data = await fileService.listMyZone();
      this.setState({
        zones: data
      });
    } catch (err) {
      message.error(err.message);
    } finally {
      this.setState({
        loading: false
      });
    }
  }
  renderEditor() {
    const {
      mode,
      fileZoneName,
      fileZoneOwner,
      fileZoneDescription
    } = this.state;
    if (mode === listMode) {
      return;
    }
    return (
      <Card size="small" title={"新建文件空间"}>
        <Form className="createForm">
          <Row gutter={12}>
            <Col span={12}>
              <Form.Item label="名称">
                <Input
                  type="text"
                  placeholder="请输入文件空间名称"
                  defaultValue={fileZoneName}
                  onChange={e => {
                    this.setState({
                      fileZoneName: e.target.value.trim()
                    });
                  }}
                />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item label="拥有者">
                <Input
                  type="text"
                  placeholder="请输入文件空间拥有者"
                  defaultValue={fileZoneOwner}
                  onChange={e => {
                    this.setState({
                      fileZoneOwner: e.target.value.trim()
                    });
                  }}
                />
              </Form.Item>
            </Col>
            <Col span={24}>
              <Form.Item label="描述">
                <TextArea
                  defaultValue={fileZoneDescription}
                  autosize={{ minRows: 3, maxRows: 6 }}
                  type="textarea"
                  placeholder="请输入文件空间的描述"
                  onChange={e => {
                    this.setState({
                      fileZoneDescription: e.target.value.trim()
                    });
                  }}
                />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Button
                type="primary"
                onClick={() => this.handleFileZoneConfirm()}
              >
                {mode === updateMode ? "更新" : "确认"}
              </Button>
            </Col>
            <Col span={12}>
              <Button
                onClick={() => {
                  this.goBack();
                }}
              >
                返回
              </Button>
            </Col>
          </Row>
        </Form>
      </Card>
    );
  }
  renderList() {
    const { mode, zones } = this.state;
    if (mode !== listMode) {
      return;
    }
    const columns = [
      {
        title: "名称",
        dataIndex: "name",
        key: "name",
        width: "120px"
      },
      {
        title: "拥有者",
        dataIndex: "owner",
        key: "owner",
        width: "120px"
      },
      {
        title: "描述",
        dataIndex: "description",
        key: "description"
      },
      {
        title: "更新时间",
        dataIndex: "updatedAt",
        key: "updatedAt",
        width: "180px",
        render: text => {
          if (!text) {
            return;
          }
          return moment(text).format(TIME_FORMAT);
        }
      },
      {
        title: "操作",
        key: "op",
        width: "120px",
        render: (text, record) => {
          return (
            <a
              href="/update"
              onClick={e => {
                e.preventDefault();
                this.setState({
                  updateFileZone: record,
                  mode: updateMode,
                  fileID: record.id,
                  fileZoneName: record.name,
                  fileZoneOwner: record.owner,
                  fileZoneDescription: record.description
                });
              }}
            >
              更新
            </a>
          );
        }
      }
    ];
    return (
      <div className="fileList">
        <Table
          className="zoneTable"
          rowKey={"id"}
          dataSource={zones}
          columns={columns}
        />
        <Button
          type="primary"
          className="createFileZone"
          onClick={() => {
            this.setState({
              mode: createMode
            });
          }}
        >
          新建文件空间
        </Button>
      </div>
    );
  }
  render() {
    const { loading } = this.state;
    return (
      <div className="MyFileZone">
        <Spin spinning={loading}>
          {this.renderEditor()}
          {this.renderList()}
        </Spin>
      </div>
    );
  }
}

MyFileZone.propTypes = {
  account: PropTypes.string.isRequired,
  roles: PropTypes.array.isRequired
};

export default MyFileZone;
