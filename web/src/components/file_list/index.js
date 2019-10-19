import React from "react";
import PropTypes from "prop-types";
import {
  Spin,
  message,
  Table,
  Icon,
  Card,
  Form,
  Button,
  Input,
  Row,
  Col
} from "antd";
import moment from "moment";
import { Link } from "react-router-dom";
import bytes from "bytes";

import "./file_list.sass";
import * as fileService from "../../services/file";
import { getQueryParams } from "../../helpers/util";
import { TIME_FORMAT } from "../../vars";
import { FILE_HANDLER_PATH, PREVIEW_PATH } from "../../paths";

class FileList extends React.Component {
  state = {
    loading: false,
    zoneName: "",
    zone: 0,
    sort: "-updatedAt",
    fields:
      "id,updatedAt,name,maxAge,zone,type,size,width,height,description,creator,thumbnail",
    keyword: "",
    files: null,
    pagination: {
      pageSizeOptions: ["10", "20"],
      showSizeChanger: true,
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
    const { loading, zone, fields, pagination, sort, keyword } = this.state;
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
        keyword,
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
  async handleSearch(e) {
    e.preventDefault();
    const { pagination } = this.state;
    const updateData = {};
    updateData.pagination = Object.assign(
      { ...pagination },
      {
        current: 1,
        total: 0
      }
    );

    this.setState(updateData, () => {
      this.fetchFiles();
    });
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
        title: "缓存有效期",
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
        title: "缩略图",
        dataIndex: "thumbnail",
        key: "thumbnail",
        width: "100px",
        render: (text, record) => {
          const { thumbnail } = record;
          if (!thumbnail) {
            return;
          }
          const data = `data:image/${record.type};base64,${record.thumbnail}`;
          return <img alt={"thumbnail"} src={data} />;
        }
      },
      {
        title: "宽/高",
        key: "widthHeight",
        width: "110px",
        render: (text, record) => {
          const { width, height } = record;
          if (!width && !height) {
            return "";
          }
          return `${width}/${height}`;
        }
      },
      {
        title: "大小",
        key: "size",
        dataIndex: "size",
        width: "100px",
        render: (text, record) => {
          return bytes(record.size);
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
        width: "200px",
        render: (text, record) => {
          let updateLink = null;
          const url =
            FILE_HANDLER_PATH.replace(":fileZoneID", zone) +
            `?name=${zoneName}&fileID=${record.id}`;
          if (record.creator === account) {
            updateLink = (
              <Link to={url}>
                <Icon type="edit" />
                更新
              </Link>
            );
          }
          return (
            <div className="op">
              {updateLink}
              <Link to={PREVIEW_PATH.replace(":id", record.id)}>
                <Icon type="file-jpg" />
                预览
              </Link>
            </div>
          );
        }
      }
    ];
    return (
      <div>
        <Card title="文件筛选" size="small" className="filter">
          <Form onSubmit={this.handleSearch.bind(this)}>
            <Row gutter={12}>
              <Col span={20}>
                <Form.Item>
                  <Input
                    allowClear
                    type="text"
                    placeholder="请输入搜索关键字，支持名称与描述的模糊搜索"
                    onChange={e => {
                      this.setState({
                        keyword: e.target.value.trim()
                      });
                    }}
                  />
                </Form.Item>
              </Col>
              <Col span={4}>
                <Form.Item>
                  <Button
                    htmlType="submit"
                    style={{
                      width: "100%",
                      letterSpacing: "5px"
                    }}
                    type="primary"
                    icon="search"
                  >
                    搜索
                  </Button>
                </Form.Item>
              </Col>
            </Row>
          </Form>
        </Card>

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
      </div>
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
