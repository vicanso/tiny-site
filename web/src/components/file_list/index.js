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
  Col,
  Select
} from "antd";
import moment from "moment";
import { Link } from "react-router-dom";
import debounce from "debounce";
import bytes from "bytes";

import "./file_list.sass";
import * as fileService from "../../services/file";
import * as imageService from "../../services/image";
import { getQueryParams, copy } from "../../helpers/util";
import { TIME_FORMAT } from "../../vars";
import { FILE_HANDLER_PATH } from "../../paths";

const { Option } = Select;

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
    clientX: 0,
    optiming: false,
    imageConfig: null,
    previewImage: null,
    previewImageData: "",
    optimImageInfo: null,
    optimWidth: 0,
    optimHeight: 0,
    optimQuality: 0,
    optimType: "",
    pagination: {
      pageSizeOptions: ["10", "20"],
      showSizeChanger: true,
      current: 1,
      pageSize: 10,
      total: 0
    }
  };
  previewImageRef = React.createRef();
  constructor(props) {
    super(props);
    this.state.zone = Number.parseInt(props.match.params.fileZoneID);
    this.state.zoneName = getQueryParams(props.location.search, "name");
    this.debounceUpdateOptimParams = debounce((...args) => {
      this.updateOptimParams(...args);
    }, 1000);
  }
  componentDidMount() {
    this.fetchFiles();
    this.fetchConfig();
  }
  async fetchConfig() {
    try {
      const data = await imageService.getConfig();
      this.setState({
        imageConfig: data
      });
    } catch (err) {
      message.error(err.message);
    }
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
  async fetchOriginalImageData(fileID) {
    try {
      const data = await fileService.getByID(fileID, {
        fields: "*"
      });
      this.setState({
        previewImageData: data.data
      });
    } catch (err) {
      message.error(err.message);
    }
  }
  preview(item) {
    const { imageConfig } = this.state;
    if (!imageConfig) {
      message.error("获取图片相关配置失败，请刷新重试");
      return;
    }
    this.setState(
      {
        optimImageInfo: null,
        previewImageData: null,
        previewImage: item
      },
      () => {
        this.optimImage();
      }
    );
    this.fetchOriginalImageData(item.id);
  }
  async optimImage() {
    const file = this.getFileName();
    this.setState({
      optiming: true
    });
    try {
      const data = await imageService.optim(file);
      this.setState({
        clientX: 0,
        optimImageInfo: data
      });
    } catch (err) {
      message.error(err.message);
    } finally {
      this.setState({
        optiming: false
      });
    }
  }
  getFileName() {
    const {
      optimQuality,
      optimWidth,
      optimHeight,
      previewImage,
      optimType
    } = this.state;
    const file = `${
      previewImage.name
    }-${optimQuality}-${optimWidth}-${optimHeight}.${optimType ||
      previewImage.type}`;
    return file;
  }
  updateOptimParams(key, value) {
    const { previewImage, optimQuality } = this.state;
    if (Number.isInteger(value) && value < 0) {
      message.error("参数不能小于0");
      return;
    }

    const update = {
      optimImageInfo: null
    };
    // 如果是转换为webp图片
    if (key === "optimType" && value === "webp") {
      // 原图片为jpeg，而且未选择转换质量，则设置为80来转换
      if (previewImage.type === "jpeg" && optimQuality === 0) {
        update.optimQuality = 80;
      }
    }

    update[key] = value;
    this.setState(update, () => {
      this.optimImage();
    });
  }
  renderPreview() {
    const {
      optiming,
      previewImage,
      previewImageData,
      optimImageInfo,
      optimWidth,
      optimType,
      optimHeight,
      optimQuality,
      imageConfig,
      clientX
    } = this.state;
    if (!previewImage) {
      return;
    }
    // 高度-(60 + 50) 60为顶部高度，50为底部工具栏高度
    // 宽度-25
    const { innerHeight, innerWidth } = window;
    const maxWidth = innerWidth - 25;
    const maxHeight = innerHeight - 110;
    let originalWidth = previewImage.width;
    let originalHeight = previewImage.height;
    let currentWidth = 0;
    let currentHeight = 0;
    // 宽高都有设置
    if (optimWidth && optimHeight) {
      currentWidth = optimWidth;
      currentHeight = optimHeight;
    } else if (optimHeight) {
      // 只设置高度，宽度自适应
      currentHeight = optimHeight;
      currentWidth = (originalWidth / originalHeight) * currentHeight;
    } else if (optimWidth) {
      // 只设置宽度，高度自适应
      currentWidth = optimWidth;
      currentHeight = (originalHeight / originalWidth) * currentWidth;
    } else {
      currentWidth = originalWidth;
      currentHeight = originalHeight;
    }
    const wTimes = currentWidth / maxWidth;
    const hTimes = currentHeight / maxHeight;
    // 如果宽度或高度超过显示区域
    if (wTimes > 1 || hTimes > 1) {
      if (wTimes > hTimes) {
        currentHeight *= maxWidth / currentWidth;
        currentWidth = maxWidth;
      } else {
        currentWidth *= maxHeight / currentHeight;
        currentHeight = maxHeight;
      }
    }

    const close = (
      <a
        href="/close"
        className="close"
        onClick={e => {
          e.preventDefault();
          this.setState({
            previewImage: null
          });
        }}
      >
        <Icon type="close-square" />
      </a>
    );
    const tips = [];
    if (!previewImageData) {
      tips.push("正在加载原图片数据");
    }
    if (!optimImageInfo) {
      tips.push("正在加载优化图片数据");
    }
    if (tips.length !== 0) {
      tips.push("请稍候...");
    }
    if (previewImage && optimImageInfo) {
      const originalSize = previewImage.size;
      const optimSize = optimImageInfo.size;
      tips.push(
        `优化后，文件大小${bytes(originalSize)} => ${bytes(
          optimSize
        )}，减少${100 - Math.floor((100 * optimSize) / originalSize)}%数据量`
      );
    }
    const title = `图片预览：${previewImage.name}`;
    const marginLeft = (maxWidth - currentWidth) / 2;
    const imgStyle = {
      position: "relative",
      marginLeft: `${marginLeft}px`,
      width: `${currentWidth}px`,
      height: `${currentHeight}px`,
      backgroundSize: "100% 100%",
      backgroundImage: `url("data:image/${previewImage.type};base64,${previewImageData}")`
    };
    const halfOffset = currentWidth / 2;
    let leftValue = halfOffset;
    if (clientX && this.previewImageRef.current) {
      leftValue = clientX - this.previewImageRef.current.offsetLeft;
    }
    if (leftValue <= 0) {
      leftValue = halfOffset;
    } else if (leftValue >= currentWidth) {
      leftValue = halfOffset;
    }
    const optimStyle = {
      borderLeft: "1px solid #fff",
      position: "absolute",
      top: "0px",
      right: "0px",
      bottom: "0px",
      backgroundPosition: "right center",
      backgroundSize: "auto 100%",
      // 还要送去左边框的1px
      left: `${leftValue - 1}px`
    };
    if (optimImageInfo) {
      optimStyle.backgroundImage = `url("data:image/${optimImageInfo.type};base64,${optimImageInfo.data}")`;
    } else {
      optimStyle.backgroundColor = "rgba(255, 255, 255, 0.4)";
    }
    const file = this.getFileName();
    const url = imageConfig.url.replace(":file", file);
    const colList = [];
    if (!optiming) {
      colList.push(
        <Col key="qualityCol" span={3}>
          <Input
            defaultValue={optimQuality}
            addonBefore="图片质量："
            type="number"
            onChange={e => {
              this.debounceUpdateOptimParams(
                "optimQuality",
                e.target.valueAsNumber || 0
              );
            }}
          />
        </Col>
      );
      colList.push(
        <Col key="widthCol" span={3}>
          <Input
            defaultValue={optimWidth}
            addonBefore="图片宽度："
            type="number"
            onChange={e => {
              this.debounceUpdateOptimParams(
                "optimWidth",
                e.target.valueAsNumber || 0
              );
            }}
          />
        </Col>
      );
      colList.push(
        <Col key="heightCol" span={3}>
          <Input
            defaultValue={optimHeight}
            addonBefore="图片高度："
            type="number"
            onChange={e => {
              this.debounceUpdateOptimParams(
                "optimHeight",
                e.target.valueAsNumber || 0
              );
            }}
          />
        </Col>
      );
      colList.push(
        <Col key="typeCol" span={2}>
          <Select
            placeholder="图片类型"
            style={{
              width: "100%"
            }}
            defaultValue={optimType || previewImage.type}
            onChange={value => {
              this.updateOptimParams("optimType", value);
            }}
          >
            <Option value="png">PNG</Option>
            <Option value="jpeg">JPEG</Option>
            <Option value="webp">WEBP</Option>
          </Select>
        </Col>
      );
      colList.push(
        <Col key="urlCol" span={5}>
          <Input
            readOnly
            addonBefore="图片地址："
            addonAfter={
              <Icon
                title="点击复制图片地址"
                onClick={e => {
                  copy(url, e.target);
                  message.info("已成功复制图片地址");
                }}
                type="copy"
              />
            }
            defaultValue={url}
          />
        </Col>
      );
    }

    return (
      <Card title={title} extra={close} size="small" className="previewWrapper">
        <div className="content">
          <div
            className="imgWrapper"
            style={imgStyle}
            ref={this.previewImageRef}
          >
            <div style={optimStyle}></div>
            <div className="imgOriginal">原图</div>
            <div className="imgOptim">预览图</div>
          </div>
          <Row className="functions" gutter={12}>
            {colList}
            <Col span={8} className="tips">
              <Icon
                type="info-circle"
                style={{
                  marginRight: "3px"
                }}
              />
              {tips.join("，")}
            </Col>
          </Row>
        </div>
      </Card>
    );
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
              <a
                style={{
                  marginLeft: "5px"
                }}
                href="/copy"
                onClick={e => {
                  e.preventDefault();
                  this.preview(record);
                }}
              >
                <Icon type="file-jpg" />
                预览
              </a>
            </div>
          );
        }
      }
    ];
    // TODO 增加搜索功能
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
  handleMouseMove(e) {
    const { previewImage } = this.state;
    // 如果非预览，则直接返回
    if (!previewImage) {
      return;
    }
    const clientX = e.clientX;
    if (clientX % 3 === 0) {
      this.setState({
        clientX
      });
    }
  }
  render() {
    const { loading } = this.state;
    return (
      <div className="FileList" onMouseMove={this.handleMouseMove.bind(this)}>
        <Spin spinning={loading}>
          {this.renderList()}
          {this.renderPreview()}
        </Spin>
      </div>
    );
  }
}

FileList.propTypes = {
  account: PropTypes.string.isRequired,
  roles: PropTypes.array.isRequired
};

export default FileList;
