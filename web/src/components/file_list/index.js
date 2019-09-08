import React from "react";
import PropTypes from "prop-types";
import { Spin, message, Table, Icon } from "antd";
import moment from "moment";
import { Link } from "react-router-dom";

import "./file_list.sass";
import * as fileService from "../../services/file";
import { getQueryParams } from "../../helpers/util";
import { TIME_FORMAT } from "../../vars";
import { FILE_HANDLER_PATH } from "../../paths";

class FileList extends React.Component {
  state = {
    loading: false,
    zoneName: "",
    zone: 0,
    sort: "-updatedAt",
    fields:
      "id,updatedAt,name,maxAge,zone,type,size,width,height,description,creator",
    files: null,
    pagination: {
      current: 1,
      pageSize: 10,
      total: 0
    }
  };
  constructor(props) {
    super(props);
    this.state.zone = Number.parseInt(props.match.params.fileZoneID);
    this.state.zoneName = getQueryParams(props.location.search, "name");
  }
  componentDidMount() {
    this.fetchFiles();
  }
  async fetchFiles() {
    const { loading, zone, fields, pagination, sort } = this.state;
    if (loading) {
      return;
    }
    this.setState({
      loading: true
    });

    try {
      const limit = pagination.pageSize;
      const offset = (pagination.current - 1) * limit;
      const data = await fileService.list({
        offset,
        limit,
        zone,
        sort,
        fields
      });
      const updateData = {
        files: data.files
      };
      if (data.count >= 0) {
        updateData.pagination = Object.assign(
          { ...pagination },
          {
            total: data.count
          }
        );
      }
      this.setState(updateData);
    } catch (err) {
      message.error(err.message);
    } finally {
      this.setState({
        loading: false
      });
    }
  }
  renderList() {
    const { account } = this.props;
    const { files, pagination, zone, zoneName } = this.state;
    const columns = [
      {
        title: "名称",
        dataIndex: "name",
        key: "name",
        width: "300px"
      },
      {
        title: "缓存有效欺",
        dataIndex: "maxAge",
        key: "maxAge",
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
        title: "上传者",
        dataIndex: "creator",
        key: "creator",
        width: "120px"
      },
      {
        title: "操作",
        key: "op",
        width: "80px",
        render: (text, record) => {
          if (record.creator !== account) {
            return;
          }
          const url =
            FILE_HANDLER_PATH.replace(":fileZoneID", zone) +
            `?name=${zoneName}&fileID=${record.id}`;
          return (
            <div className="op">
              <Link to={url}>
                <Icon type="edit" />
                更新
              </Link>
            </div>
          );
        }
      }
    ];
    return (
      <Table
        className="files"
        rowKey={"id"}
        columns={columns}
        pagination={pagination}
        dataSource={files}
        onChange={pagination => {
          this.setState(
            {
              pagination: { ...pagination }
            },
            () => {
              this.fetchFiles();
            }
          );
        }}
      />
    );
  }
  render() {
    const { loading } = this.state;
    return (
      <div className="FileList">
        <Spin spinning={loading}>{this.renderList()}</Spin>
      </div>
    );
  }
}

FileList.propTypes = {
  account: PropTypes.string.isRequired,
  roles: PropTypes.array.isRequired
};

export default FileList;
