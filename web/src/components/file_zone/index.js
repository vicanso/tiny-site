import React from "react";
import { message, Spin, Table, Icon } from "antd";
import { Link } from "react-router-dom";

import "./file_zone.sass";
import * as fileService from "../../services/file";
import { FILE_LIST_PATH, FILE_HANDLER_PATH } from "../../paths";

class FileZone extends React.Component {
  state = {
    loading: true,
    zones: null
  };
  async fetchZones() {
    this.setState({
      loading: true
    });
    try {
      const data = await fileService.listZone();
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
  componentDidMount() {
    this.fetchZones();
  }
  renderList() {
    const { zones } = this.state;
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
        title: "操作",
        key: "op",
        width: "180px",
        render: (text, record) => {
          const listURL = FILE_LIST_PATH.replace(":fileZoneID", record.id);
          const addURL = FILE_HANDLER_PATH.replace(":fileZoneID", record.id);
          return (
            <div className="op">
              <Link to={`${listURL}?name=${record.name}`}>
                <Icon type="unordered-list" />
                浏览
              </Link>
              <Link to={`${addURL}?name=${record.name}`}>
                <Icon type="plus-square" />
                添加
              </Link>
            </div>
          );
        }
      }
    ];
    return (
      <Table
        className="zoneTable"
        rowKey={"id"}
        dataSource={zones}
        columns={columns}
      />
    );
  }
  render() {
    const { loading } = this.state;
    return (
      <div className="FileZone">
        <Spin spinning={loading}>{this.renderList()}</Spin>
      </div>
    );
  }
}

export default FileZone;
